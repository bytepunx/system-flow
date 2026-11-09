---
id: T-2668
type: task
nature: improvement
title: The item store counts a generation per root, which moves when an item file is added, changed, or removed, in snapshot mode or not
status: in-progress
parent: S-0608
owner: alex
created: 2026-11-07T23:49:55Z
updated: 2026-11-09T09:12:44Z
transitions:
  - to: ready
    at: 2026-11-09T09:12:43Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T09:12:44Z
    by: system-flow
stream: S-0608
tags: [cli]
touches: [flai/internal/workitem/store.go, flai/internal/workitem/repo.go, flai/internal/workitem/store_test.go]
---
# T-2668 The item store counts a generation per root, which moves when an item file is added, changed, or removed, in snapshot mode or not

## Work

Give the parsed item store (`flai/internal/workitem/store.go`) an in-memory generation count per repository root. `Repo` answers it (`flai/internal/workitem/repo.go`), so that a caller can tell that the items it computed from have not changed. Nothing is written to disk.

The count moves on every path by which the store's items change:

- With snapshot mode off, as in a one-shot `flai` and `flai mcp`: when a list or get parses a file again because its modification time or size changed, finds a new file, or drops a removed one.
- With S-0606's snapshot mode on, as in `flai serve`: when a change the store is told of replaces the snapshot or an entry of it. A list or get then stats nothing, so that is the only place it can move.
- In either mode: when `Repo.Save` in this process saves an item.

Waits for no task of this story. The story waits for S-0606, whose snapshot mode (T-2615) this builds on.

## Done when

- `store_test.go` shows, with snapshot mode off, that the generation stays the same across two lists with no change, and moves on an edit, an added file, a removed file, and a save.
- `store_test.go` shows, with snapshot mode on, that the generation stays the same across two lists with no change, and moves on a change the store is told of and on a save.
- `go test ./internal/workitem` passes.

## Notes
