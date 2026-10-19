---
id: T-1906
type: task
nature: feature
title: A framed component's requests to the app's own server are answered by the harness, never with a 500
status: done
parent: S-0449
owner: alex
created: 2026-10-19T05:04:46Z
updated: 2026-10-19T05:37:38Z
transitions:
  - to: ready
    at: 2026-10-19T05:21:18Z
    by: agent-S-0449
  - to: in-progress
    at: 2026-10-19T05:21:18Z
    by: agent-S-0449
  - to: done
    at: 2026-10-19T05:37:38Z
    by: agent-S-0449
stream: S-0449
tags: []
touches: [flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go, design/system/component-library.md, docs/users/flai.md, docs/users/flaiover.md, design/adrs, flai/internal/components/testdata/harness/src/lib/Fetcher.svelte]
after: [T-1904, T-1905]
usage:
  source: log
  seconds: 980
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 150
      output: 62273
      cache_read: 8568913
      cache_write: 257033
      cost: 4.5809
---
# T-1906 A framed component's requests to the app's own server are answered by the harness, never with a 500

## Work

A framed component's own requests are the 500s the Components page shows (criterion 5). The harness is served by the sub-project's own dev server, so a component's `fetch` of the app's API, such as `/api/...`, reaches that dev server's server code. With no dashboard token and no host channel there, it answers 500, and the component shows "Internal Server Error" ([component-library.md § What the pass found](../../../design/system/component-library.md)). The dashboard's own log shows no 5xx from the Components page's own routes.

- `harness.mjs`: a framed component's `fetch` to the dev server's own origin for the app's server code is not sent. The harness answers it with 501 Not Implemented and a JSON body `{"error": "..."}` saying the component harness does not answer the app's requests, naming the method and path. It posts the request as an `event` (name `fetch`, args the method and the URL) so the page's event log shows it. Vite's and flai's own requests pass through: paths under `/@`, `/__`, and `/node_modules/`, the Vite ping, and a path with a file extension, such as a static asset or a module. Requests to any other origin pass through.
- Record the decision with `flai adr new` (no `--commit`): the harness answers a framed component's requests to the app's own server with 501 and logs each as an event, refining ADR-0163.
- Tests: extend the browser checks in `flai/internal/components/integration_test.go`. A component that fetches `/api/x` gets 501 with the JSON body, and the page gets the `fetch` event. A fetch of an asset path still reaches the dev server.
- A headless pass, as the one `component-library.md` describes: render every listed flaiover component's examples, or its defaults, in the harness from `bin/flai components serve --story S-0449 --project flaiover`, and click each enabled button once. Record every response of 500 or more, and every request the harness answered. Report each 500 left, with the component, the request, and why. The story's agent files those as issues.
- Docs: `design/system/component-library.md` (the harness, and § What the pass found, as this pass finds it), `docs/users/flai.md` (§ Render a component live), and `docs/users/flaiover.md` (§ Components, the event log).

It waits for T-1904, which changes `harness.mjs`, `integration_test.go`, `component-library.md`, and `docs/users/flai.md` too. It waits for T-1905, which changes `docs/users/flaiover.md`, and its pass checks both tasks' work.

## Done when

- No framed component's request reaches the dev server's server code, and none answers 500.
- The pass's list of 500s is empty, or each one left is in the final message for the story's agent to file.
- The new integration checks pass with `flai test` on the changed paths, and fail without the change.

## Notes
