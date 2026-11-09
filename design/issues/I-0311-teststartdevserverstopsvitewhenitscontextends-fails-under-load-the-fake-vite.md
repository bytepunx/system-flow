---
id: I-0311
title: "TestStartDevServerStopsViteWhenItsContextEnds fails under load: the fake vite wrote no pid"
class: defect
status: closed
count: 4
cost: 15m
first_reported: 2026-11-08T01:43:19Z
last_reported: 2026-11-09T23:09:40Z
updated: 2026-11-10T07:07:24Z
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

### 2026-11-09T22:25:39Z
Story: S-0622.
S-0622's close-out job S-0622-20261109T221518Z stopped at integration after 9m45s at 5714cb4 with flai/internal/components/devserver_test.go:430 TestStartDevServerStopsViteWhenItsContextEnds: "the fake vite wrote no pid". S-0622 changes only flai/cmd/orchestrate_permissions_test.go and flai/internal/serve/orchestrate_confirm_test.go, nothing under flai/internal/components; the host was running several stories' agents, planners, and the orchestrator at the time. The close-out was asked for again.

### 2026-11-09T23:09:40Z
Story: S-0638.
S-0638's close-out (job S-0638-20261109T230113Z, at 1928250) stopped at the integration tier after 7m57s: flai/internal/components/devserver_test.go:430 TestStartDevServerStopsViteWhenItsContextEnds, "the fake vite wrote no pid". S-0638 changes nothing in flai/internal/components. The host's load average was 90 to 112 on 24 cores.

## Remediation

Story S-0647 remediates this issue, created from it at 2026-11-10T05:12:04Z.
Closed 2026-11-10T07:07:24Z: S-0647 (T-2825, 8ceb994): the dev server tests no longer read the pid file the fake vite writes once running. A seam in flai/internal/components/devserver.go, launched, tells them the pid of each vite StartDevServer starts, which the fake's exec keeps, so a vite stopped before it ran is still waited on. TestStartDevServerStopsAViteStillStartingWhenItsContextEnds reproduces the case on every run with a start delay longer than its context.
