---
id: I-0155
title: The chart data takes about 4 s warm, against 10 ms measured on 2026-09-29, almost all of it in metrics.replayHolds
class: defect
status: closed
count: 1
cost: 4s
first_reported: 2026-10-16T06:09:48Z
last_reported: 2026-10-16T06:09:48Z
updated: 2026-10-21T12:50:23Z
---

# I-0155 The chart data takes about 4 s warm, against 10 ms measured on 2026-09-29, almost all of it in metrics.replayHolds

## Description
The chart data takes about 4 s warm, against 10 ms measured on 2026-09-29, almost all of it in metrics.replayHolds

## Instances

### 2026-10-16T06:09:48Z
Story: S-0427.
S-0427's baseline (design/experiments/S-0427-sqlite-and-postgres-read-models.md, What was done) measured stats.get with since=30d, type=story, bucket=day on this repository: 4715 ms cold, 4022 ms warm, and flai stats --json 4701 ms. design/system/server-performance.md measured the same call at 125 ms cold and 9.1 to 10.6 ms warm on 2026-09-29. Warm, statsread.Read is 450 ms (370 ms of it one git log) and metrics.Compute 3612 ms; a CPU profile puts 67 percent of the samples under metrics.replayHolds, through workitem.Holds.Claim, Holds.Overlaps, and PathsOverlap, with the garbage collector a fifth of all samples. Every charts page load and every agent's flai stats waits on it.

## Remediation

Story S-0471 remediates this issue, created from it at 2026-10-21T11:58:35Z.
Closed 2026-10-21T12:50:23Z: S-0471: replayHolds judges the holds again only at an event time that can change one, over only the open stories, with a HoldsMemo of claims between judgements, and PathsOverlap no longer allocates. One replay at this repository's scale went from 2,726 ms to 71 ms with the same result. On this repository on 2026-10-21, stats.get warm went from 5,695 ms to 894 ms and replayHolds from 73 to 22 percent of the samples. Half of what remains is a git log per call, recorded as I-0184.
