---
id: I-0151
title: A pending recommendation waits a day or more for the operator's confirmation while its story idles in progress
class: efficiency
status: open
count: 1
cost: 29h
first_reported: 2026-10-15T22:48:37Z
last_reported: 2026-10-15T22:48:37Z
updated: 2026-10-15T22:48:37Z
usage:
  strategic:
    - kind: analyzer
      seconds: 114
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 98
          output: 6305
          cache_read: 370755
          cache_write: 41688
          cost: 1.2427
        - model: claude-haiku-5-5
          input: 8
          output: 6718
          cache_read: 204809
          cache_write: 23259
          cost: 0.0083
---

# I-0151 A pending recommendation waits a day or more for the operator's confirmation while its story idles in progress

## Description
A pending recommendation waits a day or more for the operator's confirmation while its story idles in progress

## Instances

### 2026-10-15T22:48:37Z
Report: design/analysis/2026-10-15-bottlenecks.md.
The orchestrator recommends within seconds, but a recommendation ends no wait (ADR-0091); the operator's confirmation took 25 to 113 hours on the four W42 waits, and S-0391 is waiting now on TH-0461 since 2026-10-15T05:09Z.

## Impact
flai stats W42 threads.total_seconds 1,003,576, threads.confirmed count 4 total_seconds 421,270 (mean 105,317 s each); waiting.longest: S-0232 TH-0324 405,614 s, S-0337 TH-0390 325,914 s, S-0380 TH-0445 94,433 s, S-0378 TH-0436 89,637 s, all awaited alex, each with the orchestrator's recommendation posted within a minute. S-0232 cycle_time_seconds 407,617 of which wait_threads_seconds 405,614. aging: S-0391 in progress 63,454 s over p85 on TH-0461's pending recommendation.

- time_lost_per_cycle: 117h

## Remediation

Found by the analysis in [2026-10-15-bottlenecks.md](../analysis/2026-10-15-bottlenecks.md).
