---
id: T-0575
type: task
nature: feature
title: Measure the reads and record cause 4 in the design
status: done
parent: S-0159
owner: alex
created: 2026-09-29T19:58:35Z
updated: 2026-09-29T20:12:03Z
transitions:
  - to: ready
    at: 2026-09-29T19:58:49Z
    by: agent-S-0159
  - to: in-progress
    at: 2026-09-29T20:09:00Z
    by: agent-S-0159
  - to: done
    at: 2026-09-29T20:12:03Z
    by: agent-S-0159
stream: S-0159
tags: []
touches: [design/system, docs]
usage:
  source: log
  seconds: 183
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 46
      output: 20779
      cache_read: 4506125
      cache_write: 65812
      cost: 1.8435
---

# T-0575 Measure the reads and record cause 4 in the design

## Work

Time the seven reads in one process with `flai hostapi --timing` and through the table flai serve answers with, warm, on this repository, and record the result under cause 4 in `design/system/server-performance.md`; update `design/system/flai-cli.md` and any operator or user doc that says the reads start flai.

## Done when

- [x] `publish.preview` and `push.pending` measured under 50 ms warm on this repository, and recorded.
- [x] The design and docs say the reads are answered in flai serve's process.

## Notes
