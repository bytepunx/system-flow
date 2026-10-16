---
id: ADR-0163
title: "flai runs a sub-project's own Vite dev server on 127.0.0.1 behind the components host action, stops it when idle, and the dashboard frames its harness cross-origin and talks to it by postMessage under a fixed protocol"
status: accepted
date: 2026-10-16
supersedes: []
superseded_by: []
refines: [ADR-0029, ADR-0031]
---

# ADR-0163 flai runs a sub-project's own Vite dev server on 127.0.0.1 behind the components host action, stops it when idle, and the dashboard frames its harness cross-origin and talks to it by postMessage under a fixed protocol

## Context

E-0026 renders a project's own components live in the dashboard. Rendering them needs the project's own Svelte compiler, Vite config, plugins, and dependencies, which only the host has. The dashboard's container holds nothing of the project and reaches it only through flai on the host ([ADR-0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md), [ADR-0031](0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md)). [ADR-0162](0162-the-component-library-reaches-sub-projects-of-kind-sveltekit-or-svelte-finds.md) settled which sub-projects are reached and how their props are read.

Three questions remain: how the operator's browser reaches a rendered component, what gates running the project's code on the host, and how the dashboard's page and the rendered component talk. The planner asked them on TH-0483, questions 1 and 5, recommending the project's own Vite dev server on a loopback port, framed cross-origin, behind a new host action. The operator answered "1. perfect" and "5. lets gate it behind the feature for now".

## Decision

flai runs a sub-project's own Vite dev server on the host, bound to 127.0.0.1 on a port flai picks, behind the `components` host action, and the dashboard frames its harness page cross-origin and talks to it by `postMessage` under a fixed protocol.

- **The dev server.** flai writes a Vite config under `.flai-cache/components/<project>/<sub-project>/` that loads the sub-project's own config and merges a plugin flai embeds into it. It runs the sub-project's own `vite`, from its `node_modules`, in the sub-project's folder, with `--host 127.0.0.1`, a free port, and `--strictPort`. Aliases, Tailwind, the design stylesheet, and hot module replacement behave as in the app. A sub-project without `vite` installed is refused with the fix.
- **The harness.** The plugin serves a page at `/__flai/components/`. It imports the manifest's `setup` module, then mounts the one component named by `?c=<path>` or by a `select` message, inside an error boundary. A component that throws posts `error` rather than leaving a blank frame.
- **How the browser reaches it.** The page in the dashboard frames the harness URL directly, cross-origin. The browser must run on the host. The container still holds nothing of the project and opens no connection to the host: flai opens the port, and the browser reaches it, never the container. This refines ADR-0029 and ADR-0031: the channel stays the dashboard's only way to the project, and the frame is the browser's.
- **The host action.** Starting a dev server runs the project's code and dependencies on the host, so it is gated by a new host action, `components`, off until the operator runs `flai serve enable components`, as `checks` is. `components.start` and `components.stop` `{project?}` are gated by it and journalled. `components.status` and `components.list` are reads. `project` names the sub-project, as `flai components list --project` does; it may be left out when one sub-project is supported.
- **The idle stop.** `flai serve` supervises one dev server per sub-project. It stops a dev server after `components.idle`, a host setting of 30 minutes by default, without a `components.status` call, and stops every dev server when `flai serve` stops.
- **The protocol.** Every message is an object with `source: "flai-components"`, `version: 1`, and `type`. The harness URL carries the page's origin as `?origin=`. The harness accepts a message only from that origin and posts only to it; the page accepts a message only from the harness's origin, from its own frame's window. With no `origin`, the harness renders `?c=` and posts nothing.

| From | Type | Fields | Meaning |
|------|------|--------|---------|
| page | `select` | `component`, `props?` | Mount this component, by its path in the sub-project, with these props. |
| page | `props` | `props`, `functions?` | Set the mounted component's props without remounting it. Each name in `functions` gets a logger that posts `event`. |
| page | `theme` | `theme` | Set `data-theme` on the document to `light` or `dark`. |
| page | `tokens` | `light`, `dark` | Lay these CSS custom properties, by name, over the stylesheet's, under `:root` and `[data-theme='dark']`; empty objects clear them. |
| harness | `ready` | `component?` | The harness has loaded and its `setup` module has run. |
| harness | `rendered` | `component` | The component mounted, or took new props. |
| harness | `error` | `component`, `message`, `file?`, `line?`, `column?` | The component failed to import or render; `file` and `line` say where. |
| harness | `event` | `component`, `name`, `args` | A function prop was called, with its arguments as far as JSON carries them. |
| harness | `size` | `width`, `height` | The rendered content's size, in CSS pixels, each time it changes. |

## Consequences

- A component renders as the project builds it, and a saved change to its source or to the stylesheet reaches the frame as a hot update, without a rebuild or a redeploy.
- The Components page works only with the browser on the host, as the operator's is at 127.0.0.1:4242. A dashboard reached from another machine shows the catalogue but cannot frame the harness.
- The project's code runs on the host only once the operator enables `components`, and an idle dev server does not outlive its use by more than `components.idle`.
- The protocol is the contract between flai's harness and the dashboard's page. A change to it is a new `version` and an ADR.
- [component-library.md](../system/component-library.md) describes the dev server and the protocol, and [dashboard-host-channel.md](../system/dashboard-host-channel.md) the action and its methods.

## Alternatives considered

- **Proxy the dev server through the channel.** The dashboard would serve the harness and its hot-update WebSocket same-origin and authenticated, and it would work with the dashboard off the host. It refines ADR-0029 with a byte-stream method and is larger. It is the later way for a dashboard off the host.
- **A static `vite build` per change, served through a host method.** No hot updates: each change takes a rebuild, which the epic rules out.
- **No gate.** The operator's own code runs either way, but a gate keeps running it on the dashboard's word the operator's choice, as for every other host action.
