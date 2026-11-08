---
id: T-2660
type: task
nature: improvement
title: workitem keeps parsed files of any kind by path, stat, and the racy window, for the five readers to share
status: done
parent: S-0602
owner: alex
created: 2026-11-07T23:49:37Z
updated: 2026-11-08T03:07:36Z
transitions:
  - to: ready
    at: 2026-11-08T03:01:52Z
    by: agent-S-0602
  - to: in-progress
    at: 2026-11-08T03:01:53Z
    by: agent-S-0602
  - to: done
    at: 2026-11-08T03:07:36Z
    by: agent-S-0602
stream: S-0602
tags: [cli]
touches: [flai/internal/workitem/keptfiles.go, flai/internal/workitem/keptfiles_test.go]
usage:
  source: log
  seconds: 343
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 26
      output: 207
      cache_read: 3453586
      cache_write: 24818
      cost: 1.2588
---
# T-2660 workitem keeps parsed files of any kind by path, stat, and the racy window, for the five readers to share

## Work

Add `flai/internal/workitem/keptfiles.go` (S-0600 owns `kept.go`, its kept items file): a generic kept-files store, one package-level value per kind, on the model of `store.go`'s `store`, `folder`, and `kept`, which this task leaves as they are.

- A `KeptFiles[T]` built with a parse function (path to `T`) and a clone function, so that each caller gets its own copy.
- `List(dir, match)` reads the folder once, stats each file that matches, parses again only a file that is new, is not the same file (`os.SameFile`), or whose size or modification time changed, and drops the files gone from the folder. `Get(path)` does the same for one file.
- A file read within `racyWindow` of its modification time is parsed again on every read until it is older, as `store.go` does.
- One folder lock is held while a folder is looked at, so that readers at once share one parse of each changed file.
- A parse error is not kept, so the next read tries the file again.
- A parse counter, for the tests.

It waits for nothing: it is the first layer, and every reader task builds on it.

## Done when

- `keptfiles_test.go` shows that a second read parses nothing, and that a changed, an added, a removed, and a renamed-over file are each seen on the next read.
- It shows that a file written within the racy window is parsed again, and that two readers at once get equal answers from one parse.
- `flai test flai/internal/workitem/keptfiles.go flai/internal/workitem/keptfiles_test.go` passes.

## Notes
