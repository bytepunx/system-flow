---
id: T-1142
type: task
nature: improvement
title: The threads package resolves the threads left open on items being archived, and the threads.archived finding names flai thread resolve
status: backlog
parent: S-0277
owner: alex
created: 2026-10-07T01:20:10Z
updated: 2026-10-07T01:20:10Z
transitions: []
stream: S-0277
tags: [flai]
touches: [flai/internal/threads/archived.go, flai/internal/threads/archived_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go]
---
# T-1142 The threads package resolves the threads left open on items being archived, and the threads.archived finding names flai thread resolve

## Work

I-0073 recurs because nothing resolves a thread when its item is archived. `flai accept` and `flai archive` move the item to `wip/archive`, and its threads stay `open` or `answered` in `wip/threads/`. `flai check` then warns `threads.archived` at every close-out until someone runs `flai thread resolve` by hand. TH-0112, TH-0127, TH-0194, TH-0203, and TH-0232 all followed that path.

- In a new `flai/internal/threads/archived.go`, add `OnItems(r, ids)`. It lists the threads whose anchor `item` is one of the IDs and that are still `Open()`, in ID order.
- Add `ResolveOnItems(r, ids, author, reason, now)`. It resolves each of those threads with `Resolve`, so the thread gets a dated entry saying why, and returns their IDs. Resolving no thread is not an error.
- In `flai/internal/check/check.go`, change the `threads.archived` message from "resolve it or move it" to name the command: `resolve it with flai thread resolve TH-nnnn`. Threads have no archive folder, so "move it" points at nothing.
- Test both functions in `archived_test.go` on a temporary repository: an open thread and an answered thread on the IDs, a resolved one, and one on another item. Update the message assertion in `check_test.go` if it checks the text.

This task waits for no other, and no other task touches its files.

## Done when

- `OnItems` and `ResolveOnItems` exist and are tested. Open and answered threads on the given items are resolved, each with an entry carrying the reason. Resolved threads and threads on other items are left unchanged.
- The `threads.archived` finding names `flai thread resolve`.
- `flai test flai/internal/threads flai/internal/check` passes.

## Notes
