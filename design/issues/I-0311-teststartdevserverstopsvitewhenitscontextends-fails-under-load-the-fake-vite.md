---
id: I-0311
title: "TestStartDevServerStopsViteWhenItsContextEnds fails under load: the fake vite wrote no pid"
class: defect
status: open
count: 2
cost: 18m
first_reported: 2026-11-08T01:43:19Z
last_reported: 2026-11-09T09:28:38Z
updated: 2026-11-09T09:28:38Z
---

# I-0311 TestStartDevServerStopsViteWhenItsContextEnds fails under load: the fake vite wrote no pid

## Description
TestStartDevServerStopsViteWhenItsContextEnds fails under load: the fake vite wrote no pid

## Instances

### 2026-11-08T01:43:19Z
Story: S-0599.
S-0599's close-out at fd34a59 stopped at the integration tier (11m37s) on `flai/internal/components/devserver_test.go:430` `TestStartDevServerStopsViteWhenItsContextEnds`: "the fake vite wrote no pid: open /tmp/TestStartDevServerStopsViteWhenItsContextEnds.../002/pid: no such file or directory". S-0599 changes nothing under `flai/internal/components`. The host's load average was about 100 to 130 on 24 cores. `flai test flai/internal/components` passed right after, so the test looks load-dependent: it seems to read the fake vite's pid file before the process has had time to write it. The cost is one more close-out run, about 14 minutes.

### 2026-11-09T09:28:38Z
Story: S-0613.
S-0613's close-out (job S-0613-20261109T090752Z, head cf61e14) stopped at integration after 20m31s on TestStartDevServerStopsViteWhenItsContextEnds: "the fake vite wrote no pid" (devserver_test.go:430), while other stories' agents and close-outs ran on the host. S-0613 changes only flai/internal/guard and docs; the close-out had to run again.

## Remediation
