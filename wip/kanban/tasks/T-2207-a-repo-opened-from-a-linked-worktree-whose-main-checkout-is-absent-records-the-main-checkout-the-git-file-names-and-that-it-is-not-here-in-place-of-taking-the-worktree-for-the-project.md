---
id: T-2207
type: task
nature: remediation
title: A Repo opened from a linked worktree whose main checkout is absent records the main checkout the .git file names and that it is not here, in place of taking the worktree for the project
status: backlog
parent: S-0521
owner: alex
created: 2026-10-28T13:02:35Z
updated: 2026-10-28T13:02:35Z
transitions: []
stream: S-0521
tags: [cli, docker]
touches: [flai/internal/workitem/repo.go, flai/internal/workitem/worktree_test.go]
---
# T-2207 A Repo opened from a linked worktree whose main checkout is absent records the main checkout the .git file names and that it is not here, in place of taking the worktree for the project

## Work

In `flai/internal/workitem/repo.go`, `mainCheckout` returns the worktree itself when the main checkout its `.git` file names has no `system-flow.yaml`, which is the container's case (I-0219): the worktree is mounted at its host path and the main checkout's files are not. Change `Open` so that a `Repo` opened from such a worktree records what it knows: the main checkout's path the `.git` file names and that it is absent, for example a field `MainAbsent string` beside `MainRoot`, with `MainRoot` left as today so that nothing that reads it changes, and a method such as `RequireMain(what string) error` that returns nil for a normal checkout and a linked worktree whose main checkout is present, and otherwise an error naming the command (`what`), the story worktree, and the main checkout path, saying that the main checkout is not here and to run the command on the host or through the MCP tools `verify`, `task_done`, and `stream_sync`, which run there. `WorktreePath(id)` returns `Root` when `Root` is that story's worktree and the main checkout is absent, so that no caller builds `<worktree>/.flai-cache/worktrees/S-nnnn`. Nothing is refused here: the commands decide in a later task. Add the case to `worktree_test.go` beside `TestOpenFromLinkedWorktreeUsesMainWip`: a worktree with a `.git` file naming a main checkout that has no `system-flow.yaml`, opened, records the absent main, `RequireMain` errors with that path, and `WorktreePath` of the worktree's own story is `Root`. This task waits for nothing.

## Done when

- `Open` on a linked worktree whose main checkout's `system-flow.yaml` is absent records the main checkout's path and that it is absent, and `IsWorktree`, `WipDir`, and `CacheDir` keep today's values
- `RequireMain` is nil for a normal checkout and for a linked worktree with its main checkout present, and an error naming the command, the main checkout path, and the MCP tools otherwise
- `WorktreePath` of the worktree's own story is `Root` when the main checkout is absent
- `worktree_test.go` covers the three, and `go test ./internal/workitem/...` passes

## Notes

The `.git` file's `gitdir:` line and the `/.git/worktrees/` marker already give the main checkout's path (`mainCheckout`, `repo.go`); only the fallback when its manifest is absent changes.
