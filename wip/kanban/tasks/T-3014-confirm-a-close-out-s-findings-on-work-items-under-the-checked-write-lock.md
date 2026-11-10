---
id: T-3014
type: task
nature: remediation
title: Confirm a close-out's findings on work items under the checked-write lock before recording them
status: backlog
parent: S-0679
owner: alex
created: 2026-11-10T10:31:43Z
updated: 2026-11-10T10:31:43Z
transitions: []
stream: S-0679
tags: [flai, check]
touches: [flai/cmd/check.go, flai/cmd/check_test.go, flai/internal/check/itemwrite.go, docs/users/flai-reference.md]
---
# T-3014 Confirm a close-out's findings on work items under the checked-write lock before recording them

## Work

Fix the cause of I-0316 in `recordOutside` (`flai/cmd/check.go`), which both `flai check --story --record-issues` and `flai verify --record-issues` (through `recordNotes`) call. The one instance: S-0599's close-out at 2026-11-09T07:36:46Z found S-0636's body missing `## Tasks`, but S-0636 was created at 07:36:47Z and its first commit (`1a0ee3c7b`) has every section. The close-out's check read S-0636 while a checked creation was between its write and its keep or removal, which [ADR-0272](../../../design/adrs/0272-a-checked-creation-holds-a-checked-write-lock-from-its-write-until-it-keeps-or.md) says only flai serve's look reads under the checked-write lock. Confirm this from the instance before you build: if S-0636's history or the serve journal says otherwise, say so on the plan's thread.

Before it records a finding outside the story whose path is a work item under `wip/kanban`, have `recordOutside` take `(*workitem.Repo).LockCheckedWrites`, run `check.RunItemWrite` on those items, and record only the findings that run still gives. A finding gone under the lock was a checked creation's or edit's passing state ([ADR-0275](../../../design/adrs/0275-a-checked-edit-of-a-work-item-holds-the-checked-write-lock-from-its-first-save.md)): drop it, and say in the output how many were passed over. Findings on other files are recorded as today. When the lock cannot be taken, record as today and warn, so a close-out never stops on it. If the confirmation needs a helper, add it to `flai/internal/check/itemwrite.go`. If you change the help text of `--record-issues`, regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Add a test in `flai/cmd/check_test.go` that reproduces the cause: a finding outside the story on a work item that a checked write fixes before the record step (for example, hold the lock, write an item missing a section, run the record step in a goroutine, then write the item whole and let go) is not recorded, and a finding that still holds under the lock is. The test fails before the fix.

This task waits for nothing: it is the first layer, beside the container mount, with no path in common.

## Done when

- The new test fails on the code before the fix and passes after it.
- A finding outside the story on a work item is recorded only when it still holds under the checked-write lock, and a lock that cannot be taken warns and records as before.
- `flai test` on the changed paths passes.

## Notes
