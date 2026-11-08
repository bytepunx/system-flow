---
id: I-0311
title: "TestStartDevServerStopsViteWhenItsContextEnds fails under load: the fake vite wrote no pid"
class: defect
status: open
count: 1
cost: 15m
first_reported: 2026-11-08T01:43:19Z
last_reported: 2026-11-08T01:43:19Z
updated: 2026-11-08T01:43:19Z
---

# I-0311 TestStartDevServerStopsViteWhenItsContextEnds fails under load: the fake vite wrote no pid

## Description
TestStartDevServerStopsViteWhenItsContextEnds fails under load: the fake vite wrote no pid

## Instances

### 2026-11-08T01:43:19Z
Story: S-0599.
S-0599's close-out at fd34a59 stopped at the integration tier (11m37s) on `flai/internal/components/devserver_test.go:430` `TestStartDevServerStopsViteWhenItsContextEnds`: "the fake vite wrote no pid: open /tmp/TestStartDevServerStopsViteWhenItsContextEnds.../002/pid: no such file or directory". S-0599 changes nothing under `flai/internal/components`. The host's load average was about 100 to 130 on 24 cores. `flai test flai/internal/components` passed right after, so the test looks load-dependent: it seems to read the fake vite's pid file before the process has had time to write it. The cost is one more close-out run, about 14 minutes.

## Remediation
