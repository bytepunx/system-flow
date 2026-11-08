---
id: T-2667
type: task
nature: improvement
title: threads.List and Get read threads through the kept store
status: done
parent: S-0602
owner: alex
created: 2026-11-07T23:49:55Z
updated: 2026-11-08T03:22:31Z
transitions:
  - to: ready
    at: 2026-11-08T03:08:06Z
    by: agent-S-0602
  - to: in-progress
    at: 2026-11-08T03:08:07Z
    by: agent-S-0602
  - to: done
    at: 2026-11-08T03:22:31Z
    by: agent-S-0602
stream: S-0602
tags: [cli]
touches: [flai/internal/threads/threads.go, flai/internal/threads/threads_test.go]
after: [T-2660]
usage:
  source: log
  seconds: 864
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 16
      output: 114
      cache_read: 2401134
      cache_write: 17444
      cost: 0.8752
---
# T-2667 threads.List and Get read threads through the kept store

## Work

Move `List`, `Get`, and `For` in `flai/internal/threads/threads.go` onto a package-level `workitem.KeptFiles[*Thread]`, with `Read` as its parse and a deep copy of a thread and its entries as its clone. Keep `List`'s order by ID and its `WarnUnknown` call, once per parse. The writes (`Save`, `Reply`, `New`, `Resolve`, `Confirm`) save through `atomicfile`, which replaces the file, so the next read sees them by stat.

It waits for T-2660, whose store it uses.

## Done when

- `threads_test.go` shows that a reply written, a thread added, and a thread removed are each seen on the next `List`, and that a second `List` with nothing changed parses no file.
- It shows that two `List` calls at once return equal threads, and that changing a returned thread changes nothing the next `List` returns.
- `flai test flai/internal/threads/threads.go flai/internal/threads/threads_test.go` passes.

## Notes
