---
id: T-2674
type: task
nature: improvement
title: store.list parses an uncached folder's files on a bounded pool of workers
status: done
parent: S-0600
owner: alex
created: 2026-11-07T23:50:17Z
updated: 2026-11-08T00:15:16Z
transitions:
  - to: ready
    at: 2026-11-07T23:55:11Z
    by: agent-S-0600
  - to: in-progress
    at: 2026-11-07T23:55:16Z
    by: agent-S-0600
  - to: done
    at: 2026-11-08T00:15:16Z
    by: agent-S-0600
stream: S-0600
tags: [cli]
touches: [flai/internal/workitem/store.go, flai/internal/workitem/store_test.go]
usage:
  source: log
  seconds: 1200
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 19
      output: 104
      cache_read: 2188088
      cache_write: 20202
      cost: 0.7937
---
# T-2674 store.list parses an uncached folder's files on a bounded pool of workers

## Work

`store.list` (`flai/internal/workitem/store.go`, lines 61 to 96) calls `f.read` once per file in a plain loop while it holds the folder's mutex. Split it: stat every file first, reuse each `kept` entry whose key still matches as today, and hand the files that need a parse to a pool of at most `runtime.GOMAXPROCS(0)` workers, capped at 16. Gather the results under the folder's mutex, so `TestListsAtOnceShareOneRead` (eight lists at once parse each file once) still holds, and keep the order the caller sees. A parse error stays the error the loop returns now, for the first file in name order.

First layer: it waits for no other task. The kept parses under `.flai-cache` build on it.

## Done when

- [ ] A cold list of a folder parses its files on more than one goroutine, at most the bound, and a test shows each file is still parsed once.
- [ ] The existing tests in `flai/internal/workitem/store_test.go` pass unchanged, and `go test -race ./internal/workitem/` passes.
- [ ] A test lists a folder with a damaged item and gets the same error as before.

## Notes
