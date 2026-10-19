---
id: T-1659
type: task
nature: feature
title: A green batch lands each story in queue order through the acceptance flow, as whoever enqueued it
status: done
parent: S-0400
owner: alex
created: 2026-10-16T05:02:07Z
updated: 2026-10-19T22:40:27Z
transitions:
  - to: ready
    at: 2026-10-19T22:32:54Z
    by: agent-S-0400
  - to: in-progress
    at: 2026-10-19T22:32:54Z
    by: agent-S-0400
  - to: done
    at: 2026-10-19T22:40:27Z
    by: agent-S-0400
stream: S-0400
tags: [cli]
touches: [flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go]
after: [T-1644]
usage:
  source: log
  seconds: 453
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 35
      output: 17614
      cache_read: 2550035
      cache_write: 72583
      cost: 1.282
---
# T-1659 A green batch lands each story in queue order through the acceptance flow, as whoever enqueued it

## Work

- `land.go`: `Land` accepts each story of `Batch.Landing` in queue order, as whoever enqueued it, through an `Accepter` that flai serve gives as `flai accept <story> --by <by> --batch <id>`. Before each one it checks that the main branch is where the landing left it: the batch's base, then the head the last acceptance left. After each it records that head.
- `Landed` is what `flai accept --batch` calls in place of `AcceptedByHand`. It takes the story out of the queue, adds it to `Landing.Landed`, records `OutcomeLanded`, and leaves the batch in place.
- `Landable` says whether the batch lets a story land now, and at which verified commit, for the acceptance preview. The batch must be the one named, not dropped, and passed or bisected. The story must be the next to land, its branch still at the head the batch rebased, and main where the landing left it.
- `DropIfMainMoved` drops a batch whose main branch moved other than by its own landing, so that the next run builds it again from the new main.
- `Advance` takes a verified batch one step: a dropped batch is left for `RunBatch` to replace; a passed batch gets its landing; a failed one is bisected and sent back, then landed; a batch fully landed is removed, its worktree and branch with it.
- Waits for T-1644, whose `Bisect` and `SendBack` `Advance` calls.

## Done when

`Land` and `Advance` land a green batch in order, stop and drop the batch when main moved, and finish a bisected one, with a fake accepter and fake verify. The tests pass.

## Notes
