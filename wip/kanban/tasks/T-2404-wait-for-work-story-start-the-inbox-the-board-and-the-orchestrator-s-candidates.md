---
id: T-2404
type: task
nature: improvement
title: wait_for_work, story start, the inbox, the board, and the orchestrator's candidates judge an overlap by the trial merge
status: backlog
parent: S-0551
owner: alex
created: 2026-10-31T07:44:07Z
updated: 2026-10-31T07:44:07Z
transitions: []
stream: S-0551
tags: [flai, go]
touches: [flai/internal/mcpserver/work.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/orchestrate.go, flai/internal/workitem/promote.go, flai/internal/workitem/boardview.go, flai/internal/storystart/start.go, flai/internal/storystart/start_test.go, flai/internal/inbox/inbox.go]
after: [T-2389, T-2391]
---
# T-2404 wait_for_work, story start, the inbox, the board, and the orchestrator's candidates judge an overlap by the trial merge

## Work

Give the trial-merge judge to every other caller of the holds, so that what the board and the inbox show agrees with what `flai serve` starts:

- `wait_for_work` in `flai/internal/mcpserver/work.go` names an overlapping story to pull when its trial merge is clean.
- `refuse` in `flai/internal/storystart/start.go` lets `flai story start` and `story_start` start it.
- The inbox's `pull_hold` (`flai/internal/inbox/inbox.go`) and the board's `held` (`flai/internal/workitem/boardview.go`) mark it held only on a conflict, with the conflicting paths.
- The orchestrator's promotion candidates (`flai/internal/workitem/promote.go`, served by `flai/internal/mcpserver/orchestrate.go`) leave out a story only on a conflict.

`promote.go` and `boardview.go` are in `workitem`, which cannot import `storygit`: they take the judge from their caller. It waits for the storygit trial merge and the hold verdict's judge.

## Done when

- A test in `work_test.go` shows `wait_for_work` naming an overlapping story on a clean trial merge and holding it on a conflicting one, with the path in its reason.
- A test in `start_test.go` shows `story_start` starting an overlapping story on a clean trial merge.
- `flai test` passes on the changed packages.

## Notes

The folder variant of `wait_for_work`, `flai/internal/mcpserver/folder.go`, serves a project with no repository of its own; leave it to the verdict with no judge unless the story's agent finds it reachable here.
