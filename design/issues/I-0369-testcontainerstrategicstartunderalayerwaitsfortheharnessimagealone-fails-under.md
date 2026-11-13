---
id: I-0369
title: TestContainerStrategicStartUnderALayerWaitsForTheHarnessImageAlone fails under load in the race run
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-13T20:34:57Z
last_reported: 2026-11-13T20:34:57Z
updated: 2026-11-13T20:34:57Z
---

# I-0369 TestContainerStrategicStartUnderALayerWaitsForTheHarnessImageAlone fails under load in the race run

## Description
TestContainerStrategicStartUnderALayerWaitsForTheHarnessImageAlone fails under load in the race run

## Instances

### 2026-11-13T20:34:57Z
Story: S-0777.
S-0777's scripts/flai-test.sh on Linux (integration.sh --race on every core, the host shared with two other stories' agents) failed in flai/internal/serve on this test alone: container_test.go:1618 "the held run is kept as a past run", the planner's run held on its image job (ImageHeld true, Why "waiting for its agent image ... stale (built for flai 1.0.0, not dev)") found among PastPlans. The test passed in two whole short-suite runs the same hour, and 10 of 10 with go test -race -count=10 -run on it alone, so it fails only under load: the check of PastPlans races the held run's move to the past runs. S-0777 changes no file under flai/internal/serve.

## Remediation
