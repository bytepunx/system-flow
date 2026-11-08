---
id: T-2670
type: task
nature: improvement
title: board.get reads the stories' undone paths at once on a bounded pool, and reads a story in review's once
status: done
parent: S-0604
owner: alex
created: 2026-11-07T23:50:03Z
updated: 2026-11-08T02:20:04Z
transitions:
  - to: ready
    at: 2026-11-08T02:07:54Z
    by: agent-S-0604
  - to: in-progress
    at: 2026-11-08T02:07:55Z
    by: agent-S-0604
  - to: done
    at: 2026-11-08T02:20:04Z
    by: agent-S-0604
stream: S-0604
tags: [cli]
touches: [flai/internal/hostapi/hostapi.go, flai/internal/hostapi/reads_test.go, flai/internal/storygit/undone.go]
usage:
  source: log
  seconds: 729
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 45
      output: 10848
      cache_read: 6237400
      cache_write: 217148
      cost: 2.6953
---
# T-2670 board.get reads the stories' undone paths at once on a bounded pool, and reads a story in review's once

## Work

`board.get`'s `undone.read` phase in `flai/internal/hostapi/hostapi.go` calls `undoesMain` for each ready, in-progress, and review card, one after another. Each call is `storygit.RecordedUndone` in `flai/internal/storygit/undone.go`, which runs one `git diff` per recorded path and may rewrite the story's record.

Read the cards at once on a small bounded pool, such as 8 workers, and keep the cards in the order they have today. Each story has its own record file, so stories do not share a write. Check that `RecordedUndone` and `MainBranch` are safe to call at once. A story in review is read here and again in `mergequeue.Compute`. Let one read serve both if the two calls can share it without changing either answer. Otherwise leave it, and say so in this task's Notes.

Waits for no task: the release and merge queue tasks change other files.

## Done when

- [ ] A test in `flai/internal/hostapi/reads_test.go` gives `board.get` several stories with undone records and checks that each card carries the same undone paths as before, in the same order.
- [ ] `go test -race ./internal/hostapi/... ./internal/storygit/...` passes, including `timing_test.go`'s phase names.

## Notes

- `RecordedUndone` and `MainBranch` are safe at once: each story reads and rewrites only its own record file, written atomically, and runs read-only git; `MainBranch` is one `rev-parse` with no shared state. `undone.go` is unchanged.
- A story in review is still read twice, here and in `mergequeue.Compute`. One read serving both would need `board.get` to pass what it read into the merge queue, a change to `mergequeue`'s API beyond this task; the T-2666 change already spares the queue every read while no story is in review.
