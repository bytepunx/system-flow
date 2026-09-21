---
id: I-0030
title: The watcher's test for a file still being written fails on a slow runner
class: defect
status: closed
count: 2
cost: 15m
first_reported: 2026-09-20T13:50:06Z
last_reported: 2026-09-21T22:48:43Z
updated: 2026-09-21T22:54:28Z
---

# I-0030 The watcher's test for a file still being written fails on a slow runner

## Description
The watcher's test for a file still being written fails on a slow runner

## Instances

### 2026-09-20T13:50:06Z
2026-09-20: TestAFileStillBeingWrittenIsReportedOnce (internal/watch, from S-0073) failed in the flai workflow on main after the S-0078 acceptance and passed on a rerun. It writes every 15 ms against a 40 ms tick and expects one report; a runner that pauses for two ticks between writes makes the watcher, correctly, report the file as settled midway. The test should drive the clock or the ticks itself instead of sleeping.

### 2026-09-21T22:48:43Z
2026-09-21: failed again in the flai workflow on main after S-0087 was pushed (seen [wip/a.md x7], want [wip/a.md]; 3.64s, so the writer loop took ~0.6s instead of ~0.12s). Not reproducible locally, not even with 48 spinning processes on 24 cores, and unrelated to S-0087. The 2026-09-20 remediation note (drive the ticks instead of sleeping) was not acted on; the second failure is the trigger. Fixed by S-0093.

## Remediation
Closed 2026-09-21T22:54:28Z: S-0093: the debounce is a unit fed snapshots and the test helper writes in one step; no watcher test asserts a count a stalled goroutine can change
