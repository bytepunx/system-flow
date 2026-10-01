---
id: T-0632
type: task
nature: remediation
title: A thread's entries in one second by one author share a heading
status: done
parent: S-0179
owner: arobson
created: 2026-10-01T08:23:14Z
updated: 2026-10-01T08:24:02Z
transitions:
  - to: ready
    at: 2026-10-01T08:23:32Z
    by: claude-opus-5-5
  - to: in-progress
    at: 2026-10-01T08:23:32Z
    by: claude-opus-5-5
  - to: done
    at: 2026-10-01T08:24:02Z
    by: claude-opus-5-5
stream: S-0179
tags: []
touches: [flai/internal/threads]
usage:
  source: log
  seconds: 30
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 6
      output: 4003
      cache_read: 698771
      cache_write: 7586
      cost: 0.2805
---
# T-0632 A thread's entries in one second by one author share a heading

## Work

`threads.Reply` and `threads.Resolve` append to the last entry when it carries the same stamp and author, as narratives and issues do (I-0043). Tests cover reply then resolve in one second, and two authors in one second.

## Done when

- [x] A reply then a resolution in one second by one author leave one entry heading, and `Entries` returns the merged text
- [x] `make test` passes

## Notes

Part of S-0179.
