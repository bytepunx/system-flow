---
id: T-1367
type: task
nature: remediation
title: TestRoundTripRepositoryItems round-trips each item from one read of its file and skips a file gone since the list
status: backlog
parent: S-0290
owner: alex
created: 2026-10-08T08:42:15Z
updated: 2026-10-08T08:42:15Z
transitions: []
stream: S-0290
tags: [flai, tests]
touches: [flai/internal/workitem/workitem_test.go]
---
# T-1367 TestRoundTripRepositoryItems round-trips each item from one read of its file and skips a file gone since the list

## Work

Every instance in I-0079 is one race. In a story worktree, `Open` puts `wip/` in the main checkout (`Repo.WipDir`, ADR-0019). `TestRoundTripRepositoryItems` parses every item through `r.List(true)`, then reads each file again with `os.ReadFile` and compares `it.Marshal()` with the second read. When an agent or flai rewrites an item between the two reads, the test compares two versions of the file and fails. `Repo.Save` writes through `atomicfile.WriteFile`, so a reader never sees a half-written file: one read is always a whole, consistent version.

- Move the per-item check into a helper in `flai/internal/workitem/workitem_test.go`. For each path `List` returned, it reads the file once and parses those bytes with `ParseItem`. It runs `Validate` on that item and compares its `Marshal()` with the same bytes.
- Skip a path whose read fails with `os.ErrNotExist`. `flai archive` or a move can take the file away after the list. Any other read error is still a failure; today it is ignored (`orig, _ :=`).
- Keep the rest of the test as it is: the live checkout, the skip under `-short`, the count of items, and the `NextID` check. The point is that every hand-written item round-trips byte for byte.
- Add a test that reproduces I-0079 in a temporary project (`newProject`). Create items and list them. Then rewrite one through `Repo.Save` with a changed field, such as its title or an `updated` timestamp, and remove another. Run the helper on the listed items and assert it reports nothing. Against the old two-read check, the same steps report the rewritten item as "marshal differs from file". Say so in the test's comment.

This task waits for nothing.

## Done when

- `TestRoundTripRepositoryItems` reads each item file once and compares the marshal with that read.
- A file gone since the list is skipped, and any other read error fails the test.
- The new reproduction test passes. It rewrites one listed item and removes another between the list and the check.
- `flai test flai/internal/workitem/workitem_test.go` passes, and `go test -run 'TestRoundTrip' ./internal/workitem` passes in `flai/` without `-short`.

## Notes

Drafted by planner-S-0290. Reading the story worktree's own committed `wip/` would also be stable, but it checks a stale copy and needs a hook into `Open`. One read per file keeps the check on the live items, so it is the smaller change.
