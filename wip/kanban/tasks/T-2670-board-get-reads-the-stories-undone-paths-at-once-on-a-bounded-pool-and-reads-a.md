---
id: T-2670
type: task
nature: improvement
title: board.get reads the stories' undone paths at once on a bounded pool, and reads a story in review's once
status: backlog
parent: S-0604
owner: alex
created: 2026-11-07T23:50:03Z
updated: 2026-11-07T23:50:03Z
transitions: []
stream: S-0604
tags: [cli]
touches: [flai/internal/hostapi/hostapi.go, flai/internal/hostapi/reads_test.go, flai/internal/storygit/undone.go]
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
