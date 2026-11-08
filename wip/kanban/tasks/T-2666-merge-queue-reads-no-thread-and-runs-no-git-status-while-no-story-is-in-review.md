---
id: T-2666
type: task
nature: improvement
title: merge.queue reads no thread and runs no git status while no story is in review
status: backlog
parent: S-0604
owner: alex
created: 2026-11-07T23:49:54Z
updated: 2026-11-07T23:49:54Z
transitions: []
stream: S-0604
tags: [cli]
touches: [flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go]
---
# T-2666 merge.queue reads no thread and runs no git status while no story is in review

## Work

`ComputeAt` in `flai/internal/mergequeue/eligible.go` reads `threads.List(repo)` before its item loop finds a story in review. That read is 94% of `merge.queue`'s 185 to 630 ms. It also runs `mainDirt`, a `git status --porcelain`, every time.

Find the stories in review first. When there are none, read no threads and run no `mainDirt`, and answer as today. When there is one, read them once, as today. If anything on the board needs `mainDirt` with no story in review, keep that read, and say why in this task's Notes.

Waits for no task: it shares no file with the release or hostapi tasks.

## Done when

- [ ] A test in `flai/internal/mergequeue/eligible_test.go` shows that `ComputeAt` with no story in review reads no thread, and answers what it answered before.
- [ ] `TestEachOtherStoryInReviewHasEveryReasonItIsNotEligible` and the rest of the package's tests pass.

## Notes
