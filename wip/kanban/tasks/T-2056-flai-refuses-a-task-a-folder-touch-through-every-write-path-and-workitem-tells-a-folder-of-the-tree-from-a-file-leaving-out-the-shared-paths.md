---
id: T-2056
type: task
nature: improvement
title: flai refuses a task a folder touch through every write path, and workitem tells a folder of the tree from a file, leaving out the shared paths
status: backlog
parent: S-0468
owner: alex
created: 2026-10-22T21:04:59Z
updated: 2026-10-22T21:04:59Z
transitions: []
stream: S-0468
tags: [cli]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/create.go, flai/internal/workitem/create_test.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/folder_test.go, flai/cmd/touches.go, flai/cmd/touches_test.go, flai/internal/mcpserver/items_write_test.go]
after: [T-2055]
---
# T-2056 flai refuses a task a folder touch through every write path, and workitem tells a folder of the tree from a file, leaving out the shared paths

## Work

Waits for the ADR task, so that the refusal's terms are the ADR's (points 1 and 2).

- `flai/internal/workitem/hold.go`, beside `path()`: a function the refusals and the check share, such as `FolderTouch(root string, projects, shared, entry) (files int, folder bool)`, that reads the entry as `path()` does (a component's name or tag is its sub-project's path), says whether it names a folder of the tree at `root` holding at least one file, and how many, and answers false for a path absent from the tree and for a folder lying wholly inside a `claims.shared` pattern (`manifest/shared.go` has the matcher). Nothing in storage tells a folder from a file today; the tree does.
- `flai/internal/workitem/create.go`: `CleanTouches` stays syntax only. Add the task rule beside it and call it from `Repo.Create` for a task, so that `flai task new --touches` (`cmd/items.go` reaches `Repo.Create`) refuses a folder touch, writing nothing. The refusal names the folder and its file count and says to name the files the task changes, or to leave the folder to the story's touches with a `### Planning` note.
- `flai/internal/itemedit/itemedit.go` `Apply`, for a task, the same rule: it is behind `flai edit --touches` and MCP `item_edit`. `itemedit/widen.go` is untouched: `flai task done` widens with files from `git diff --name-only` (ADR-0144).
- `flai/cmd/touches.go`: `flai touches` saves through `repo.Save` without `Apply`, so it applies the rule itself for a task, before the save.
- MCP `item_new` reaches `itemnew.Create` and `Repo.Create`, and `item_edit` reaches `Apply`, so `items_write.go` needs no change; a test in `items_write_test.go` shows the refusal comes through both tools with the message as data.

Tests: `hold_test.go` for `FolderTouch` (a folder with files, a file, an absent path, a folder inside a shared pattern, a component's tag); `create_test.go` for `task new` refused and nothing written, and accepted with a file, an absent path, and a shared folder; a new `itemedit/folder_test.go` for `Apply`; `cmd/touches_test.go` for `flai touches` on a task. Run the Go tiers of the packages the story touches, and fix any existing test that creates a task with a folder touch in a tree that holds files there; most test repositories hold no such folder, so most pass untouched.

## Done when

- A folder touch on a task is refused, writing nothing, by `flai task new`, `flai edit --touches`, `flai touches`, `item_new`, and `item_edit`, with the folder, its file count, and what to do instead in the message.
- A file, an absent path, and a folder inside `claims.shared` are accepted on a task as before.
- The tests above pass, `go vet`, `gofmt`, and `golangci-lint` are clean, and `flai test` on the changed paths passes.

## Notes

Layer 2, alone: the planner's rule and the check rule both take `FolderTouch` from here, so they wait for it.
