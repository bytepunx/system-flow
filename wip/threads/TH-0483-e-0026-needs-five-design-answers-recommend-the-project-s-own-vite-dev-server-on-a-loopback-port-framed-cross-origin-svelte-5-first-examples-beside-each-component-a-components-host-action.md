---
id: TH-0483
title: "E-0026 needs five design answers: recommend the project's own Vite dev server on a loopback port, framed cross-origin; Svelte 5 first; examples beside each component; a components host action"
anchor:
  path: wip/kanban/epics/E-0026-component-library.md
  item: E-0026
status: resolved
participants: [planner-E-0026, alex]
created: 2026-10-16T21:23:52Z
updated: 2026-10-16T22:08:07Z
---

# TH-0483 E-0026 needs five design answers: recommend the project's own Vite dev server on a loopback port, framed cross-origin; Svelte 5 first; examples beside each component; a components host action

On wip/kanban/epics/E-0026-component-library.md.

## Entries

### 2026-10-16T21:23:52Z planner-E-0026
E-0026 asks for the project's own components, rendered live. `design/system/design-system-tool.md` § The project's own components left exactly this out of E-0020: it needs the project's own Svelte and Tailwind build, which only the host has, and ADR-0029 and ADR-0031 keep the dashboard's container away from the project. I am drafting the stories on the recommended answers below. The cost of delay input is already set (600 USD a week on the epic), so nothing else waits. Answer "all as recommended", or give the numbers you want otherwise.

## Questions

1. **How the browser reaches the components.** Recommended: flai runs the sub-project's own Vite dev server on the host, bound to 127.0.0.1 on a port flai picks, and the Components page frames it cross-origin. Vite's hot module replacement then works unchanged, for source and stylesheet alike. It needs the browser on the host, as yours is at 127.0.0.1:4242. Alternatives:
   - Proxy the dev server's HTTP and its HMR WebSocket through the channel, so the dashboard serves it same-origin and authenticated. Works when the dashboard runs off the host, but refines ADR-0029 with a byte-stream method. It is larger, and could come later.
   - A static `vite build` per change, served through a host method. No HMR, so a change takes a rebuild, which outcome 3 rules out.
2. **Frameworks.** Recommended: Svelte 5 with Vite (`kind: sveltekit` and a new `kind: svelte`) now, behind an adapter seam keyed by the sub-project's `kind`. React and Vue are later stories. Alternative: all three now, which triples the harness and the props extraction.
3. **Where props come from.** Recommended: flai runs a script it embeds with the project's own Node and its own `svelte/compiler`, which reads each component's `$props()` and its TypeScript type. Types it cannot resolve, such as an imported interface, get a JSON editor. Alternative: a Go parser in flai, which needs no Node but cannot read TypeScript types faithfully.
4. **Where saved examples live.** A component with required props, such as a card that needs an item, renders only from an example. Recommended: `<Component>.examples.json` beside the component, saved through `file.save` on main or on a story's branch, so that they are versioned and agents keep them up to date. Alternatives: under `design/components/<sub-project>/`, which keeps code folders clean; or in the browser only, which loses them on another machine.
5. **The host action.** Recommended: a new host action, `components`. It is off until you enable it with `flai serve enable components`, as `checks` is. It gates starting a dev server, which runs the project's own code and dependencies on the host, and flai stops an idle dev server after 30 minutes. Alternative: no gate, since the operator's own code runs either way.

## What the answers change

- 1 decides the ADR and the harness in the second story; a proxy instead adds a story.
- 2 sets which sub-projects the first story lists.
- 3 changes only the first story.
- 4 changes only the examples story's touches.
- 5 changes only the gating.

### 2026-10-16T21:29:20Z alex
1. perfect 2. yes, Svelte 5 and Vite first 3. let's start with that 4. lets put supporting files under design/components, 5. lets gate it behind the feature for now

### 2026-10-16T21:51:05Z planner-E-0026
Taken as answered. 1, 2, 3, and 5 stand as recommended. Answer 4 changed two stories:

- **S-0443** keeps each component's examples in one JSON file under `design/components/<sub-project>/`, at a path that mirrors the component's. The design folder is already in `file.save`'s reach, so ADR-0147 is unchanged; the save only learns to create a file. The documentation convention gains a second exception to "Everything is markdown", in both copies and the standard, through a new task, T-1894, with the folder's README here and in the template.
- **S-0444** now writes flaiover's examples under `design/components/flaiover/`. It changes no file of flaiover's own and no longer holds eight stories of E-0021, E-0022, and E-0024.

The whole plan is on the epic's plan thread.

### 2026-10-16T22:08:07Z alex
Resolved.
