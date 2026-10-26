---
id: I-0184
title: Every stats.get and flai stats runs one git log for claims drift, about 450 ms of a 900 ms warm answer
class: efficiency
status: open
count: 1
cost: 450ms
first_reported: 2026-10-21T12:50:20Z
last_reported: 2026-10-21T12:50:20Z
updated: 2026-10-26T04:01:29Z
---

# I-0184 Every stats.get and flai stats runs one git log for claims drift, about 450 ms of a 900 ms warm answer

## Description
Every stats.get and flai stats runs one git log for claims drift, about 450 ms of a 900 ms warm answer

## Instances

### 2026-10-21T12:50:20Z
Story: S-0471.
S-0471 measured stats.get with since=30d, type=story, bucket=day on this repository on 2026-10-21 after fixing the hold replay (I-0155): 894 ms warm, of which exec.git.log 453 ms, under a load average of 11 to 21. statsread.Read runs storygit.Committed, one git log --name-only over main and every story branch for claims.drift, at every call, even in flai serve's long-running process where nothing else is re-read. The rest of a warm call is about 330 ms of CPU in metrics.Compute (160 ms of it replayHolds) and 170 ms in statsread.Read. The 2026-09-29 figure was 9.1 to 10.6 ms warm. Figures in design/system/server-performance.md, section A regression in the chart data.

## Remediation

Story S-0500 remediates this issue, created from it at 2026-10-26T04:01:29Z.
