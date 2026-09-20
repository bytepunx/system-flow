---
id: I-0030
title: The watcher's test for a file still being written fails on a slow runner
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-09-20T13:50:06Z
last_reported: 2026-09-20T13:50:06Z
updated: 2026-09-20T13:50:06Z
---

# I-0030 The watcher's test for a file still being written fails on a slow runner

## Description
The watcher's test for a file still being written fails on a slow runner

## Instances

### 2026-09-20T13:50:06Z
2026-09-20: TestAFileStillBeingWrittenIsReportedOnce (internal/watch, from S-0073) failed in the flai workflow on main after the S-0078 acceptance and passed on a rerun. It writes every 15 ms against a 40 ms tick and expects one report; a runner that pauses for two ticks between writes makes the watcher, correctly, report the file as settled midway. The test should drive the clock or the ticks itself instead of sleeping.

## Remediation
