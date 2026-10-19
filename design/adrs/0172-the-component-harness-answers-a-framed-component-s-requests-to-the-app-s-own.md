---
id: ADR-0172
title: "The component harness answers a framed component's requests to the app's own server with 501 and logs each as an event"
status: accepted
date: 2026-10-19
supersedes: []
superseded_by: []
refines: [ADR-0163]
topics: [components, dashboard, cli]
---

# ADR-0172 The component harness answers a framed component's requests to the app's own server with 501 and logs each as an event

## Context

The harness is served by the sub-project's own Vite dev server ([ADR-0163](0163-flai-runs-a-sub-project-s-own-vite-dev-server-on-127-0-0-1-behind-the.md)). A component that fetches the app's API, such as `/api/...`, reaches that dev server's server code. There it has no dashboard token and no host channel, so it answers 500, and the component shows "Internal Server Error". An `EventSource` for `/api/events` fails the same way. S-0449's criterion 5 asks that no component interaction on the Components page result in a 500.

## Decision

Framed, with `?origin=`, the harness answers a component's requests to the app's own server itself, and posts each as an `event`.

- Before the `setup` module or any component runs, the harness replaces `window.fetch` and `window.EventSource`.
- A request is for the app's own server when it goes to the dev server's own origin, is not Vite's or flai's own, and is not for a file. Paths under `/@`, `/__`, and `/node_modules/`, Vite's ping (`Accept: text/x-vite-ping`), and a path whose last segment has a file extension, such as a module or a static asset, go through. So does a request to any other origin.
- A `fetch` of the app's own server is not sent. It is answered with status 501 Not Implemented, `content-type: application/json`, and `{"error": "flai's component harness does not answer the app's own requests: <METHOD> <path>"}`. The harness posts `{type: "event", component, name: "fetch", args: [method, path]}`, the path with its query.
- An `EventSource` for the app's own server opens no connection. It fails as a refused one does: an `error` event, then `readyState` 2, closed. The harness posts `{type: "event", component, name: "EventSource", args: [path]}`.
- The component is the one the page last selected.
- Opened directly, with no `?origin=`, the harness changes neither.

## Consequences

- No framed component's request reaches the app's server code, so none is answered 500. A component that shows its request's `error` shows the harness's message, which names the request.
- The page's event log lists each such request beside the calls of function props. The protocol does not change: these are `event` messages, with names the harness gives, `fetch` and `EventSource`.
- A component whose whole view comes from the app's server still shows nothing useful in the harness. Excluding it from the catalogue stays the way to leave it out.
- `XMLHttpRequest`, `navigator.sendBeacon`, a form's submission, and a `WebSocket` are not covered. No flaiover component uses them; a project whose components do still reaches its server with them.
- [component-library.md](../system/component-library.md) describes the rule under the harness.

## Alternatives considered

- **Answer with fixtures.** Examples could carry the answers a component's requests get. That is a larger format change, for the examples file and the page, and can follow.
- **Proxy the requests to the dashboard.** The dashboard would answer them with its token and host channel. A component's writes, such as an accept, would then act on the project from a preview.
- **Answer 404, or let the request fail as a network error.** 404 reads as a missing route, and a network error says nothing of why. 501 says the server does not do this, and the body says which request.
- **Answer only paths under `/api`.** That is flaiover's convention, not SvelteKit's: an app's routes may be anywhere, and a page route's `load` is server code too. The rule by origin, prefix, and file extension needs no setting per project.
