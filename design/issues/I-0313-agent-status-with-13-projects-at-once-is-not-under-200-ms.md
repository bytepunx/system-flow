---
id: I-0313
title: agent.status with 13 projects at once is not under 200 ms
class: efficiency
status: open
count: 1
cost: 1s
first_reported: 2026-11-08T02:32:22Z
last_reported: 2026-11-08T02:32:22Z
updated: 2026-11-10T08:34:55Z
---

# I-0313 agent.status with 13 projects at once is not under 200 ms

## Description
agent.status with 13 projects at once is not under 200 ms

## Instances

### 2026-11-08T02:32:22Z
Story: S-0603.
S-0603 fixed cause 12 of design/system/server-performance.md: agents.json is kept parsed, the agentsMu queue is gone, and Activity reads the conversations once a call. Measured with 13 projects at once at a load average of 95 to 101 on 24 cores, a call's p50 fell from 369-463 ms to 221-283 ms, but only 14 to 26 of 65 calls answered under 200 ms. This repository's call still takes about 1.1 s: alone, its Activity takes 540 to 600 ms, of which reading 909 threads is 180 to 210 ms, 381 conversations 70 to 180 ms, and the holds with their trial merges the rest. The thread and conversation reads are cause 13 (S-0602), whose criteria time threads.list, inbox.designer, and stats.get but not agent.status, so nothing measures agent.status against 200 ms once both are in. Raised on TH-0911.

## Remediation

Story S-0678 remediates this issue, created from it at 2026-11-10T08:34:55Z.
