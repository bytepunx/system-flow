---
id: T-2619
type: task
nature: improvement
title: Measure a watcher walk and item.get with 13 at once before and after on the operator's host, and record them with the snapshot and the slow tier in the design
status: done
parent: S-0606
owner: alex
created: 2026-11-07T22:49:36Z
updated: 2026-11-08T01:39:29Z
transitions:
  - to: ready
    at: 2026-11-08T01:21:00Z
    by: agent-S-0606
  - to: in-progress
    at: 2026-11-08T01:21:04Z
    by: agent-S-0606
  - to: done
    at: 2026-11-08T01:39:29Z
    by: agent-S-0606
stream: S-0606
tags: [cli]
touches: [design/system/server-performance.md, design/system/flai-cli.md, design/tech/go-libraries.md]
after: [T-2617, T-2618]
usage:
  source: log
  seconds: 1105
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 71
      output: 17848
      cache_read: 15741017
      cache_write: 423448
      cost: 5.9065
  strategic:
    - kind: orchestrator
      seconds: 8
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 12
          cache_read: 62744
          cache_write: 12650
          cost: 0.0192
---
# T-2619 Measure a watcher walk and item.get with 13 at once before and after on the operator's host, and record them with the snapshot and the slow tier in the design

## Work

Build two flai binaries the same hour, `main` and the story branch, and run each from the main checkout as earlier rounds in `design/system/server-performance.md` did. Time one walk of this project's watched folders, and `item.get` with 13 asked at once through the host API's method table in one process with the snapshot on. Record the method, the host's load, the counts of files and items, and the figures before and after under cause 16 of S-0596's section. Then describe the snapshot and the watcher's slow tier where `design/system/flai-cli.md` describes the store (S-0156) and the watcher, and correct `design/tech/go-libraries.md`, which says the watcher stats the three folders every 300 ms.

Waits for T-2617 and T-2618: it measures what they deliver together.

## Done when

- [ ] `design/system/server-performance.md` records the walk and `item.get` with 13 at once, before and after, with how they were measured; the walk is under 30 ms and `item.get` under 50 ms, or the shortfall is stated with its cause.
- [ ] `design/system/flai-cli.md` and `design/tech/go-libraries.md` describe the snapshot and the slow tier.

## Notes
