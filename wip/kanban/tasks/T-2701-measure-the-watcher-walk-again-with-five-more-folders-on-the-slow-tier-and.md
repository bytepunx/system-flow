---
id: T-2701
type: task
nature: feature
title: Measure the watcher walk again with five more folders on the slow tier, and record them and the reason under cause 16
status: done
parent: S-0606
owner: alex
created: 2026-11-09T07:36:13Z
updated: 2026-11-09T08:09:59Z
transitions:
  - to: ready
    at: 2026-11-09T08:02:42Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T08:02:43Z
    by: system-flow
  - to: done
    at: 2026-11-09T08:09:59Z
    by: agent-S-0606
stream: S-0606
tags: []
touches: [design/system/server-performance.md, design/system/flai-cli.md, design/tech/go-libraries.md]
after: [T-2700]
usage:
  source: log
  seconds: 436
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 32
      output: 178
      cache_read: 9160728
      cache_write: 21875
      cost: 4.4521
---

# T-2701 Measure the watcher walk again with five more folders on the slow tier, and record them and the reason under cause 16

## Work

Measure one watcher walk of this project's watched folders again, as T-2619 and T-2698 did, with `main` and the branch built the same hour, against the main checkout and read only: the fast look and the slow look, CPU and wall, median and p90, with the load. Record it under cause 16 in `design/system/server-performance.md`, with the five folders added to the slow tier and why (TH-0908), and say whether the walk is under 30 ms. Describe the watcher's tiers in `design/system/flai-cli.md` and `design/tech/go-libraries.md`.

Waits for T-2700, whose change it measures.

## Done when

- [ ] `server-performance.md` records the walk after T-2700 with how it was measured, and says whether it is under 30 ms.
- [ ] `flai-cli.md` and `go-libraries.md` name the slow tier's six folders.

## Notes
