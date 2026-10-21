---
id: T-1938
type: task
nature: remediation
title: The dashboard's design, its host channel design, and the user and operator docs describe the live-events WebSocket
status: done
parent: S-0454
owner: alex
created: 2026-10-20T22:22:04Z
updated: 2026-10-21T01:16:30Z
transitions:
  - to: ready
    at: 2026-10-21T01:02:02Z
    by: agent-S-0454
  - to: in-progress
    at: 2026-10-21T01:02:03Z
    by: agent-S-0454
  - to: done
    at: 2026-10-21T01:16:30Z
    by: agent-S-0454
stream: S-0454
tags: [dashboard, docs]
touches: [design/system/flaiover-dashboard.md, design/system/dashboard-host-channel.md, design/tech/sveltekit.md, docs/users/flaiover.md, docs/operators/index.md, docs/users/flai.md]
after: [T-1934, T-1936, T-1937]
usage:
  source: log
  seconds: 867
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 33
      output: 12166
      cache_read: 1848371
      cache_write: 49098
      cost: 0.9095
---
# T-1938 The dashboard's design, its host channel design, and the user and operator docs describe the live-events WebSocket

## Work

Describe what T-1934, T-1936, and T-1937 built, linking the ADR of T-1931 rather than repeating it:

- `design/system/flaiover-dashboard.md`: the `GET /api/events` row of the table in § API (S-0011) becomes the WebSocket upgrade, its messages, its ping, and its 401 and 403. § Runtime ("push updates to open tabs with server-sent events"), § Where the data comes from (S-0073), and the line that says "reached an open event stream" say WebSocket. Bump `updated`.
- `design/system/dashboard-host-channel.md`: the Watching row that names `/api/events` and `sse.ts`, and the lines that say flaiover already uses server-sent events, say what is true now; the alternatives table keeps its history. Bump `updated`.
- `design/tech/sveltekit.md`: the `ws` row names the live-events socket beside `/agent`.
- `docs/users/flaiover.md`: § Components' Events line names a `WebSocket` beside the `EventSource`, and any word on live updates says many tabs no longer hang the dashboard.
- `docs/operators/index.md`: "even with event streams open" says what is true now.

It waits for the three tasks it describes, and T-1937's harness change is what its Components line names.

## Done when

- No document under `design/system`, `design/tech`, or `docs` says the dashboard's live events are server-sent events, but as history.
- The markdown lint and `flai check --strict` pass.

## Notes

The removal of the endpoint itself is the last task; write these as the state after it.
