---
id: T-2304
type: task
nature: feature
title: wait_for_events answers the orchestrator end true with why reasons done when nothing queued remains for it
status: done
parent: S-0536
owner: alex
created: 2026-10-28T21:46:09Z
updated: 2026-10-29T08:01:33Z
transitions:
  - to: ready
    at: 2026-10-29T07:46:25Z
    by: agent-S-0536
  - to: in-progress
    at: 2026-10-29T07:46:26Z
    by: agent-S-0536
  - to: done
    at: 2026-10-29T08:01:33Z
    by: agent-S-0536
stream: S-0536
tags: [flai, mcp, orchestrator]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/orchestrator_end_test.go]
after: [T-2302]
usage:
  source: log
  seconds: 907
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 18
      output: 6805
      cache_read: 1317925
      cache_write: 41887
      cost: 0.6365
---
# T-2304 wait_for_events answers the orchestrator end true with why reasons done when nothing queued remains for it

## Work

ADR-0209 §4, criterion 3 (the end). `wait_for_events` in `flai/internal/mcpserver/server.go` answers a story's agent `end: true` with a `why` from `endWhy` (ADR-0154). Give the orchestrator's session (`s.orchestrator()`, `FLAI_ROLE=orchestrate`) the same answer, `end: true`, `why: reasons done`, after any events already behind its cursor, when the reason queue T-2302 keeps in the serve folder holds nothing for its project. A session that cannot find or read the queue, such as one run by hand outside flai serve, holds as today.

Update the tool's description and the `End` field's jsonschema text, and the server's instructions where they say the orchestrator holds `wait_for_events` between decisions.

Waits for T-2302, which writes the queue this reads.

## Done when

- A test in `orchestrator_end_test.go` gives an orchestrator session an empty queue and gets `end: true` and `why: reasons done`; with a reason queued, the call holds; for a story's agent nothing changes.
- The description and the instructions say when the orchestrator is told to end.
- `flai test` passes on the changed paths.

## Notes
