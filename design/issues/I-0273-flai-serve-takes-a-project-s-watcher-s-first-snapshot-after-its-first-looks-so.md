---
id: I-0273
title: flai serve takes a project's watcher's first snapshot after its first looks, so a change made in between waits up to a minute for the next look
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-11-07T05:23:04Z
last_reported: 2026-11-07T05:23:04Z
updated: 2026-11-11T02:33:20Z
---

# I-0273 flai serve takes a project's watcher's first snapshot after its first looks, so a change made in between waits up to a minute for the next look

## Description
flai serve takes a project's watcher's first snapshot after its first looks, so a change made in between waits up to a minute for the next look

## Instances

### 2026-11-07T05:23:04Z
Story: S-0569.
Found by S-0569's T-2528 while reproducing I-0247's TestServeStartsADefinedRoleOnAnEventAndOnASchedule failure. watch.Watcher.Run takes its first snapshot inside its own goroutine, which serve starts after its first looks at the project. A change written in that window, about 20 ms after "project served" in the reproduction, is already in the snapshot. So the watcher never reports it, and only serve's once-a-minute look (time.After(time.Minute) in serve.go) acts on it: the run for a story moved to review was queued at 60.15 s, in 5 of 25 runs under a pinned CPU load. The test now waits for the watcher to report a probe file. flai serve is unchanged. Taking the snapshot before the first looks would close the window.

## Remediation

Story S-0690 remediates this issue, created from it at 2026-11-10T08:49:01Z.
Closed 2026-11-11T02:33:20Z: S-0690 fixed it. watch.Watcher.Mark (T-2981) takes the watcher's first snapshot in the caller's goroutine, and Run compares its first tick with it. flai serve's reconcile now calls it before the item snapshot and the project's first looks (T-2984), so a change made while those looks run is reported within a tick, not at the minute's look. TestAChangeMadeDuringServesFirstLooksIsActedOnWithinSeconds in flai/internal/serve/serve_test.go reproduces it. It moves a story to review from a first look's Agent hook and needs the watcher role's run within 10 s. It failed 20 of 20 runs on the old serve.go and passes 25 of 25 under -race. TestServeStartsADefinedRoleOnAnEventAndOnASchedule no longer needs its watcherReports probe wait.
