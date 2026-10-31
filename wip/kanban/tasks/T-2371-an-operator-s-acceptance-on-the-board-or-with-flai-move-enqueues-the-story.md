---
id: T-2371
type: task
nature: improvement
title: An operator's acceptance on the board or with flai move enqueues the story while the queue batches, and drops no batch
status: done
parent: S-0552
owner: alex
created: 2026-10-31T07:33:07Z
updated: 2026-10-31T08:04:38Z
transitions:
  - to: ready
    at: 2026-10-31T07:34:56Z
    by: system-flow
  - to: in-progress
    at: 2026-10-31T07:34:57Z
    by: system-flow
  - to: done
    at: 2026-10-31T08:04:38Z
    by: agent-S-0552
stream: S-0552
tags: [flai]
touches: [flai/cmd/accept.go, flai/cmd/accept_queue_test.go, flai/cmd/move.go, flai/cmd/move_test.go, flai/internal/preview/accept.go, docs/users/flai-reference.md, flai/cmd/accept_queue_tiers_test.go, flai/cmd/move_refused_test.go]
usage:
  source: log
  seconds: 1781
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 90
      output: 441
      cache_read: 3973174
      cache_write: 253957
      cost: 1.47
---
# T-2371 An operator's acceptance on the board or with flai move enqueues the story while the queue batches, and drops no batch

## Work

While the merge queue can batch (`queueBatches`) and is not held, `flai accept <story>` by the operator, and `flai move <story> done`, which runs it, enqueue the story with `mergequeue.Add` as the operator instead of landing it, answer that it is enqueued (text and `--json`), and change nothing else: no merge, no move, no batch dropped. A story already enqueued or in the batch under way is answered as such, unchanged. A refusal of the enqueue (not eligible, not in review) is the acceptance's refusal. `--batch` (the queue's own landing), `--dry-run`, and the orchestrator's path stay as they are. While the queue is held, or cannot batch, the hand acceptance lands as today under ADR-0221's gate and `--without-queue-tiers`. Waits for nothing: it shares no file with the other layer-1 tasks.

## Done when

- [ ] Tests show the operator's accept and move to done enqueue a story in review while the queue batches, leave a batch under way as it was, and land as before while the queue is held.

## Notes

Layer 1.
