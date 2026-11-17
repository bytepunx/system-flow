---
id: T-3598
type: task
nature: remediation
title: A state the project adds with limit true takes a WIP limit through flai board limit, whose help names no baseline states, and the baseline-only LimitedColumns and DefaultWIPLimits are retired
status: backlog
parent: S-0844
owner: alex
created: 2026-11-17T12:37:00Z
updated: 2026-11-17T12:37:00Z
transitions: []
stream: S-0844
tags: [flai, board]
touches: [flai/internal/workitem/board.go, flai/internal/workitem/board_update_test.go, flai/internal/workitem/boardview_test.go, flai/cmd/board.go, flai/cmd/board_limit_test.go, docs/users/flai-reference.md]
---
# T-3598 A state the project adds with limit true takes a WIP limit through flai board limit, whose help names no baseline states, and the baseline-only LimitedColumns and DefaultWIPLimits are retired

## Work

This task waits for no other task. It shares no file with the other tasks of the first layer.

I-0380's instance 4 says that `workitem.LimitedColumns` is fixed to ready, in-progress, and review, that `Board.SetLimit` refuses any other state, and that `Board.Save` drops a limit on one. That was true when the issue was filed. Since then `LoadBoard` gives the board the repository's workflow, so `SetLimit` and `Save` follow it, and `TestBoardLimitsFollowTheWorkflow` in `flai/internal/workitem/boardview_test.go` covers that. Three things are left:

- `LimitedColumns` and `DefaultWIPLimits` in `flai/internal/workitem/board.go` still read only the baseline, and no production code calls them. Remove them. Where a test needs them, use `limitedStates` of the workflow it builds, so no caller can read the baseline's three columns again.
- The help of `flai board limit` in `flai/cmd/board.go` says "Set the WIP limit of ready, in-progress, or review". Make it say "a state whose limit is true in the workflow", with the baseline's three as the example. This is also the second half of I-0403. Run `make flai-reference` so that `docs/users/flai-reference.md` follows the help.
- Add a case to `flai/cmd/board_limit_test.go` that reproduces instance 4. In a project whose `design/workflows/board.md` adds a state `testing` with `limit: true`, `flai board limit testing 2` writes the limit. A later write of the board, such as `flai order`, keeps it. `flai board limit` refuses a state without `limit: true`.

## Done when

- `LimitedColumns` and `DefaultWIPLimits` are gone, and nothing in `flai/` reads the baseline's limited states, apart from the baseline workflow itself.
- `flai board limit --help` names no baseline state as the only choice, and `docs/users/flai-reference.md` matches it.
- The new test passes, and a limit on a state the project adds survives a later write of the board.
- `flai test` on the changed paths passes.

## Notes
