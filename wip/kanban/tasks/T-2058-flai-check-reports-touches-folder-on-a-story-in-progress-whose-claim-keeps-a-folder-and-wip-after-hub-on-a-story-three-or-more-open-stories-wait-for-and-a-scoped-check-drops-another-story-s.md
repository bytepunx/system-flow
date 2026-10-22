---
id: T-2058
type: task
nature: improvement
title: flai check reports touches.folder on a story in progress whose claim keeps a folder and wip.after-hub on a story three or more open stories wait for, and a scoped check drops another story's
status: backlog
parent: S-0468
owner: alex
created: 2026-10-22T21:05:53Z
updated: 2026-10-22T21:05:53Z
transitions: []
stream: S-0468
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go]
after: [T-2056]
---
# T-2058 flai check reports touches.folder on a story in progress whose claim keeps a folder and wip.after-hub on a story three or more open stories wait for, and a scoped check drops another story's

## Work

Waits for the task rule, whose `workitem.FolderTouch` the first finding reads. Builds the ADR's points 4 and 5. There is no rule registry: each rule is a method of `checker` in `check.go` that `Run` calls and that emits with `c.add` (`wip.overlap` in `overlap()`, `story.component-tag` in `componentTag()`, the `after` rules in `after()` are the models).

- `touches.folder`, a warning on the story's file: for each story in progress, take its claim from `Holds.Claim` (ADR-0096's narrowing applied, so a folder a task names a file inside is already gone) and report each entry `FolderTouch` says is a folder with files, with the folder, its file count, and the ready stories the hold rules hold through it (`Holds.Of` on each ready story, matching the reason's paths, or the overlap pairs `overlap()` already computes). The message says to name the files its tasks change with `flai touches`, or to say under `### Planning` why files may be added there.
- `wip.after-hub`, a warning on the hub's file: for each story not done or cancelled, count the stories in backlog, ready, or in progress that name it in `after`; at three or more, report the hub and name them, and say that `after` names only a story whose change the story needs on main, and that the pull order sets an order of preference.
- `scope.go`: a check scoped to a story (the close-out's) keeps `touches.folder` on that story and drops it on another story, recording no issue, in the list at the `item.archive`, `board.wip-limit`, `issues.no-story`, `threads.archived` drop; the same for `wip.after-hub` when the hub is not the story. The convention's close-out paragraph in `work-management.md` gains the two in the convention task.

Tests in `check_test.go`: a story in progress touching `flai/internal/serve` (a folder the test tree fills) with tasks naming files elsewhere gives one `touches.folder` naming the ready story it holds; a task naming a file inside gives none; a folder inside `claims.shared` gives none; a hub named by three open stories gives one `wip.after-hub`, by two gives none, and a done story that names it does not count. In `scope_test.go`: another story's `touches.folder` and `wip.after-hub` are dropped from a scoped check and recorded in no issue; the story's own is kept.

## Done when

- `flai check` reports the two findings as described, as warnings, with the folder, its count, and the stories in the message, and `flai check --strict` fails on them.
- A check scoped to another story drops them and records no issue; scoped to the story it keeps them.
- The tests above pass, `go vet`, `gofmt`, and `golangci-lint` are clean, and `flai test` on the changed paths passes, and the main checkout's `flai check` reports no `touches.folder` on a story this story did not change (none of the stories in progress keeps a folder the tree fills, as of planning; if one does, say so in the narrative rather than narrowing another story's touches).

## Notes

Layer 3, beside the planner's rule task; no path in common with it.
