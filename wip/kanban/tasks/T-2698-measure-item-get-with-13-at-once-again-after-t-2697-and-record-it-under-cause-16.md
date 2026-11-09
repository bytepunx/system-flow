---
id: T-2698
type: task
nature: feature
title: Measure item.get with 13 at once again after T-2697, and record it under cause 16
status: done
parent: S-0606
owner: alex
created: 2026-11-08T01:41:10Z
updated: 2026-11-08T02:12:00Z
transitions:
  - to: ready
    at: 2026-11-08T02:07:19Z
    by: agent-S-0606
  - to: in-progress
    at: 2026-11-08T02:07:20Z
    by: agent-S-0606
  - to: done
    at: 2026-11-08T02:12:00Z
    by: agent-S-0606
stream: S-0606
tags: []
touches: [design/system/server-performance.md, design/system/flai-cli.md]
after: [T-2697]
usage:
  source: log
  seconds: 280
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 44
      output: 14473
      cache_read: 7203197
      cache_write: 300818
      cost: 3.3702
---

# T-2698 Measure item.get with 13 at once again after T-2697, and record it under cause 16

## Work

Measure `item.get` with 13 at once again, as T-2619 did, with the branch after T-2697 against `main`, and record the figures under cause 16 in `design/system/server-performance.md`; say in `design/system/flai-cli.md` that `item.get` copies only what it answers. If the operator's answer on TH-0908 moves more folders to the slow tier, measure the walk again too.

Waits for T-2697, whose change it measures.

## Done when

- [ ] `server-performance.md` records `item.get` with 13 at once after T-2697, with how it was measured, and whether it is under 50 ms.

## Notes
