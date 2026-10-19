---
id: T-1669
type: task
nature: feature
title: flai serve lands a green batch, bisects a red one, sends the failing story back with a reason and a thread, and drops a batch when main moves
status: done
parent: S-0400
owner: alex
created: 2026-10-16T05:02:54Z
updated: 2026-10-19T23:20:00Z
transitions:
  - to: ready
    at: 2026-10-19T23:01:37Z
    by: agent-S-0400
  - to: in-progress
    at: 2026-10-19T23:01:37Z
    by: agent-S-0400
  - to: done
    at: 2026-10-19T23:20:00Z
    by: agent-S-0400
stream: S-0400
tags: [cli]
touches: [flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go]
after: [T-1643, T-1659]
usage:
  source: log
  seconds: 1103
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 174
      output: 87075
      cache_read: 12606100
      cache_write: 358813
      cost: 6.3376
---
# T-1669 flai serve lands a green batch, bisects a red one, sends the failing story back with a reason and a thread, and drops a batch when main moves

## Work

- After `RunBatch`, the queue pass calls `mergequeue.Advance`. It lands a green batch, bisects a red one, sends the failing story back with a reason and a thread, and drops a batch when main moves, all in the pass's own goroutine.
- The accepter runs this flai's own binary as `flai accept <story> --by <by> --batch <id>` in the main checkout.
- Each outcome is logged once as a lifecycle event.
- `mergequeue_test.go` drives a pass through a green batch, a red one, and main moving under a batch, with fake runners.
- Waits for T-1643 and T-1659: the accepter needs `--batch`, and the pass calls `Advance`.

## Done when

flai serve's queue pass lands, bisects, sends back, and drops in the tested cases, and the tests pass.

## Notes
