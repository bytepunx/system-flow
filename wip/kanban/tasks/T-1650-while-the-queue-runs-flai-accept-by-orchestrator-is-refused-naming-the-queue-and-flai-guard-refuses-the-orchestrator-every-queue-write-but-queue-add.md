---
id: T-1650
type: task
nature: feature
title: While the queue runs, flai accept --by orchestrator is refused, naming the queue, and flai guard refuses the orchestrator every queue write but queue_add
status: done
parent: S-0403
owner: alex
created: 2026-10-16T05:01:25Z
updated: 2026-10-20T00:19:45Z
transitions:
  - to: ready
    at: 2026-10-20T00:02:24Z
    by: agent-S-0403
  - to: in-progress
    at: 2026-10-20T00:02:25Z
    by: agent-S-0403
  - to: done
    at: 2026-10-20T00:19:45Z
    by: agent-S-0403
stream: S-0403
tags: [cli]
touches: [flai/cmd/accept.go, flai/cmd/accept_orchestrator_test.go, flai/internal/guard/writes.go, flai/internal/guard/orchestrate_test.go, flai/internal/guard/call.go, flai/internal/guard/guard.go, flai/cmd/guard.go]
after: [T-1649]
usage:
  source: log
  seconds: 1040
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 94
      output: 32758
      cache_read: 6872356
      cache_write: 249325
      cost: 3.5769
---
# T-1650 While the queue runs, flai accept --by orchestrator is refused, naming the queue, and flai guard refuses the orchestrator every queue write but queue_add

## Work

- In `flai accept`, refuse `--by orchestrator` while the `merge-queue` host action is on for the project, unless the queue itself is landing. Name `queue_add` in the refusal.
- In `flai guard`, under `FLAI_ROLE=orchestrate`, refuse `flai queue remove`, `hold`, and `resume`, and their host writes. Allow `flai queue add` and `queue_add`.
- Test both refusals and what stays allowed.
- When the queue lands a story the orchestrator enqueued (`flai accept --by orchestrator --batch`), write the evidence and verified commit recorded with the enqueue under `### Accepted by the orchestrator` in the story's Notes.
- Waits for T-1649, which records the evidence in the queue's state.

## Done when

The orchestrator can only enqueue while the queue runs, and the tests pass.

## Notes
