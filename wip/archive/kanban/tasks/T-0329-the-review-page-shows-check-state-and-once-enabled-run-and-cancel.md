---
id: T-0329
type: task
nature: feature
title: The review page shows check state and, once enabled, Run and Cancel
status: done
parent: S-0082
owner: alex
created: 2026-09-22T22:37:07Z
updated: 2026-09-22T23:23:49Z
transitions:
  - to: ready
    at: 2026-09-22T23:01:11Z
    by: system-flow
  - to: in-progress
    at: 2026-09-22T23:01:11Z
    by: system-flow
  - to: done
    at: 2026-09-22T23:23:49Z
    by: system-flow
stream: S-0082
tags: []
---
# T-0329 The review page shows check state and, once enabled, Run and Cancel

## Work
`/api/items/[id]/checks` (`GET` = `checks.status`; `POST {action: "run"|"cancel"}` maps to `checks.run`/`checks.cancel`, the same `{action}` shape `/api/dashboard` uses) and `/api/items/[id]/checks/tail?from=N` (`GET` = `checks.tail`, forwarded as an NDJSON stream reusing the exact pattern `/api/items/[id]/accept` already established for progress, since `checks.tail`'s progress events are content lines, not an accept's steps).

On `Review.svelte`: a "Checks" section showing the current or last state (which named check, if any, is running; pass/fail per check so far; outcome, duration once ended); **Run checks** and, while running, **Cancel**, gated on `checks_enabled` from `project.info.host_actions`, off says what enables it (`flai serve enable checks`), the same shape `HostPanel.svelte` already uses for its own gated buttons. While a run is active the page polls `checks.status` for the summary and drives the tail endpoint in a loop (each call returns once it has new content or its own wait elapses, then the client calls again with the returned offset) for the live output, shown as scrolling text; the loop stops itself once `running` is false. A story's page opened after the run already finished still shows the last outcome and duration, not nothing.

## Done when
Component and route tests (mirroring `HostPanel.svelte.test.ts` and `dashboard.test.ts`'s patterns): the gated-off state, a live run's polling and tail-loop termination once it ends, cancel, and a story whose last run already finished. `pnpm run check`/`test:unit`/`lint` clean.

## Notes
Depends on T-0328. Reuses `/api/items/[id]/accept`'s NDJSON-over-one-request pattern for `tail` rather than inventing a second one; unlike accept's, a browser closing this connection loses nothing, since the run itself lives independently of any one request (T-0327) — the next `tail` call, from this tab or another, picks up wherever the log file is.
