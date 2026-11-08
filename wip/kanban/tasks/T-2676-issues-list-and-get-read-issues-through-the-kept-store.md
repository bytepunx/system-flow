---
id: T-2676
type: task
nature: improvement
title: issues.List and Get read issues through the kept store
status: backlog
parent: S-0602
owner: alex
created: 2026-11-07T23:50:21Z
updated: 2026-11-07T23:50:21Z
transitions: []
stream: S-0602
tags: [cli]
touches: [flai/internal/issues/issues.go, flai/internal/issues/issues_test.go]
after: [T-2660]
---
# T-2676 issues.List and Get read issues through the kept store

## Work

Move `List` and `Get` in `flai/internal/issues/issues.go` onto a package-level `workitem.Kept[*Issue]`, with `Read` as its parse and a deep copy as its clone. Keep `List` leaving out `summary.md` and `README.md`, and keep its order. `Save` writes in place with `os.WriteFile`, so the store sees a bump by size or modification time and the racy window. The test covers a bump that keeps the file's size.

It waits for T-2660, whose store it uses.

## Done when

- `issues_test.go` shows that an issue bumped, added, and removed is each seen on the next `List`, the bump written in place at the same size included, and that a second `List` with nothing changed parses no file.
- It shows that two `List` calls at once return equal issues.
- `flai test flai/internal/issues/issues.go flai/internal/issues/issues_test.go` passes.

## Notes
