---
id: T-2267
type: task
nature: improvement
title: storygit.Committed keeps the stories' commits in a file under .flai-cache with its branch-tip key, so a fresh process whose key matches reads no git log
status: done
parent: S-0531
owner: alex
created: 2026-10-28T18:24:28Z
updated: 2026-10-28T19:00:49Z
transitions:
  - to: ready
    at: 2026-10-28T18:26:33Z
    by: agent-S-0531
  - to: in-progress
    at: 2026-10-28T18:26:37Z
    by: agent-S-0531
  - to: done
    at: 2026-10-28T18:31:03Z
    by: agent-S-0531
stream: S-0531
tags: [flai]
touches: [flai/internal/storygit/committed.go, flai/internal/storygit/committed_test.go, flai/internal/statsread/statsread_test.go, design/system/flai-cli.md, design/system/server-performance.md]
usage:
  source: log
  seconds: 266
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 27
      output: 9461
      cache_read: 1538224
      cache_write: 59651
      cost: 0.902
---
# T-2267 storygit.Committed keeps the stories' commits in a file under .flai-cache with its branch-tip key, so a fresh process whose key matches reads no git log

## Work

In `flai/internal/storygit/committed.go`, keep `Committed`'s answer in a file under the main checkout's `.flai-cache` (`repo.CacheDir()`) beside the per-process memory S-0500 added. The file holds the key, the main branch, the wip folder, and `git show-ref --heads`, and the answer. On a miss in memory, read the file: when its key matches, answer from it, keep it in memory, and run no `git log`. After a log read, write the file atomically, as the other kept files under `.flai-cache` are written. A missing, unreadable, or damaged file, or one with another key, is a miss that reads the log. A failed write never fails the call. Each call still gets its own deep copy.

In `flai/internal/storygit/committed_test.go`, with the counting runner S-0500's tests use, add tests that clear the in-memory store to stand for a fresh process: one with nothing moved starts no `git log`, one after a tip moved reads it again and rewrites the file, a damaged file reads the log, and an unwritable `.flai-cache` still answers. In `flai/internal/statsread/statsread_test.go`, add a fresh-process variant of `TestReadStartsNoLogWhileNothingChanged`.

This task waits for none.

## Done when

- A second process's `Committed`, with no branch tip moved, starts no `git log`, as a test shows.
- A moved tip, a damaged file, and an unwritable cache folder each still give the answer the log gives, as tests show.
- `flai test` on the three files passes.

## Notes
