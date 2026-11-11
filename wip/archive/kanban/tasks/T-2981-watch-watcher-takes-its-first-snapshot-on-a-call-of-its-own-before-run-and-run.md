---
id: T-2981
type: task
nature: remediation
title: watch.Watcher takes its first snapshot on a call of its own before Run, and Run compares its first tick with it
status: done
parent: S-0690
owner: alex
created: 2026-11-10T10:29:24Z
updated: 2026-11-11T02:16:58Z
transitions:
  - to: ready
    at: 2026-11-11T02:12:25Z
    by: agent-S-0690
  - to: in-progress
    at: 2026-11-11T02:12:25Z
    by: agent-S-0690
  - to: done
    at: 2026-11-11T02:16:58Z
    by: agent-S-0690
stream: S-0690
tags: [flai, watch]
touches: [flai/internal/watch/watch.go, flai/internal/watch/watch_test.go]
usage:
  source: log
  seconds: 273
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 640
      cache_read: 191636
      cache_write: 5475
      cost: 0.0949
---
# T-2981 watch.Watcher takes its first snapshot on a call of its own before Run, and Run compares its first tick with it

## Work

Today `Watcher.Run` in `flai/internal/watch/watch.go` takes its first snapshot inside itself (`newDebounce(w.snapshot())`), so a caller that starts `Run` in a goroutine cannot know when the baseline was taken. A change written before that snapshot is already in it and is never reported (I-0273).

Add a method that takes the first snapshot at once, in the caller's goroutine, and keeps it, such as `Watcher.Mark()`. `Run` compares its first tick with the kept snapshot when there is one, and takes its own as today when there is none, so the merge queue's `watchRefs` and every other caller keep working unchanged. A second `Mark`, or a `Mark` after `Run` started, needs no support: say so in the method's comment.

Add a test in `flai/internal/watch/watch_test.go`: a file added, one changed, and one removed after `Mark` and before `Run` starts are each reported by `Run`'s first ticks. A second test shows that `Run` with no `Mark` still reports a change made after it started.

This is the first layer: it waits for no task.

## Done when

- `Watcher` has a method that takes the first snapshot before `Run`, and `Run` compares with it.
- A change made between that call and `Run` is reported, in a test that fails without the method.
- The watch package's existing tests pass unchanged.
- `flai test flai/internal/watch` passes.

## Notes

The method is `Watcher.Mark`. `TestAChangeMadeBetweenMarkAndRunIsReported` covers a fast and a slow path. The "Run with no Mark" test is the existing `TestReportsChangedAddedAndRemovedFilesUnderTheWatchedPathsOnly`, unchanged. `flai test` on `flai/internal/watch` passed.
