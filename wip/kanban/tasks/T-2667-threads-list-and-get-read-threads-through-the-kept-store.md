---
id: T-2667
type: task
nature: improvement
title: threads.List and Get read threads through the kept store
status: backlog
parent: S-0602
owner: alex
created: 2026-11-07T23:49:55Z
updated: 2026-11-07T23:49:55Z
transitions: []
stream: S-0602
tags: [cli]
touches: [flai/internal/threads/threads.go, flai/internal/threads/threads_test.go]
after: [T-2660]
---
# T-2667 threads.List and Get read threads through the kept store

## Work

Move `List`, `Get`, and `For` in `flai/internal/threads/threads.go` onto a package-level `workitem.Kept[*Thread]`, with `Read` as its parse and a deep copy of a thread and its entries as its clone. Keep `List`'s order by ID and its `WarnUnknown` call, once per parse. The writes (`Save`, `Reply`, `New`, `Resolve`, `Confirm`) save through `atomicfile`, which replaces the file, so the next read sees them by stat.

It waits for T-2660, whose store it uses.

## Done when

- `threads_test.go` shows that a reply written, a thread added, and a thread removed are each seen on the next `List`, and that a second `List` with nothing changed parses no file.
- It shows that two `List` calls at once return equal threads, and that changing a returned thread changes nothing the next `List` returns.
- `flai test flai/internal/threads/threads.go flai/internal/threads/threads_test.go` passes.

## Notes
