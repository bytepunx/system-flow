---
id: T-1643
type: task
nature: feature
title: The acceptance preview and flai accept take a passing batch's verification of the story's rebased head in place of the story's own verify
status: done
parent: S-0400
owner: alex
created: 2026-10-16T05:00:59Z
updated: 2026-10-19T22:55:47Z
transitions:
  - to: ready
    at: 2026-10-19T22:40:41Z
    by: agent-S-0400
  - to: in-progress
    at: 2026-10-19T22:40:41Z
    by: agent-S-0400
  - to: done
    at: 2026-10-19T22:55:47Z
    by: agent-S-0400
stream: S-0400
tags: [cli]
touches: [flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/preview/orchestrator.go, flai/cmd/accept.go, flai/cmd/accept_queue_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1659]
usage:
  source: log
  seconds: 906
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 70
      output: 34945
      cache_read: 5059056
      cache_write: 143998
      cost: 2.5434
---
# T-1643 The acceptance preview and flai accept take a passing batch's verification of the story's rebased head in place of the story's own verify

## Work

- `flai accept --batch <id>` is the queue's landing. Its preview takes `mergequeue.Landable` in place of the story head's own verify. For the orchestrator it also takes it in place of `--verified` equal to the branch head, and needs no `--evidence`. Each reason `Landable` gives is a blocker.
- With `--batch`, the acceptance calls `mergequeue.Landed`, not `AcceptedByHand`, and writes `### Accepted by the orchestrator` with the batch's verified commit when the orchestrator enqueued.
- `accept_queue_test.go` lands a story with `--batch` beside a story accepted by hand. It checks that the done transition, the acceptance commit, the archive, and the `flai stats` figures match. It also checks a refusal for each blocker.
- Waits for T-1659, whose `Landable` and `Landed` it calls.

## Done when

`flai accept --batch` lands a story as an acceptance by hand does, refuses one the batch does not let land, and the tests pass.

## Notes
