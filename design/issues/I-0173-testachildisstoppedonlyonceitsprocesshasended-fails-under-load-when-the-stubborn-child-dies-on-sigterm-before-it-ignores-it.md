---
id: I-0173
title: TestAChildIsStoppedOnlyOnceItsProcessHasEnded fails under load when the stubborn child dies on SIGTERM before it ignores it
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-21T06:07:23Z
last_reported: 2026-10-21T06:07:23Z
updated: 2026-10-21T06:07:23Z
---

# I-0173 TestAChildIsStoppedOnlyOnceItsProcessHasEnded fails under load when the stubborn child dies on SIGTERM before it ignores it

## Description
TestAChildIsStoppedOnlyOnceItsProcessHasEnded fails under load when the stubborn child dies on SIGTERM before it ignores it

## Instances

### 2026-10-21T06:07:23Z
Story: S-0455.
S-0455's close-out stopped at the integration tier (5m56s): flai/internal/host/host_test.go:355 TestAChildIsStoppedOnlyOnceItsProcessHasEnded reported "stopping names the process that is still running: {Name:mcp Root:/a State:stopping PID:3654399 ...}, alive false". The test sleeps 100 ms after the child is running and takes it that the "stubborn" stub has by then set itself to ignore SIGTERM; under the integration run's load it had not, so the stop's SIGTERM ended it before the test saw it stopping. S-0455 changes nothing under flai/internal/host. A fix would have the stub signal that it ignores SIGTERM (a line on stdout or a file) and the test wait for that rather than sleep. TestAChildThatWillNotStopIsKilled shares the same sleep.

## Remediation
