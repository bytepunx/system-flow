---
id: I-0369
title: TestContainerStrategicStartUnderALayerWaitsForTheHarnessImageAlone fails under load in the race run
class: defect
status: closed
count: 2
cost: 13m
first_reported: 2026-11-13T20:34:57Z
last_reported: 2026-11-15T22:29:21Z
updated: 2026-11-16T01:23:39Z
---

# I-0369 TestContainerStrategicStartUnderALayerWaitsForTheHarnessImageAlone fails under load in the race run

## Description
TestContainerStrategicStartUnderALayerWaitsForTheHarnessImageAlone fails under load in the race run

## Instances

### 2026-11-13T20:34:57Z
Story: S-0777.
S-0777's scripts/flai-test.sh on Linux (integration.sh --race on every core, the host shared with two other stories' agents) failed in flai/internal/serve on this test alone: container_test.go:1618 "the held run is kept as a past run", the planner's run held on its image job (ImageHeld true, Why "waiting for its agent image ... stale (built for flai 1.0.0, not dev)") found among PastPlans. The test passed in two whole short-suite runs the same hour, and 10 of 10 with go test -race -count=10 -run on it alone, so it fails only under load: the check of PastPlans races the held run's move to the past runs. S-0777 changes no file under flai/internal/serve.

### 2026-11-15T22:29:21Z
Story: S-0790.
S-0790's race tier (flai test race), run after its board lock change, failed twice at container_test.go:1618 with 'the held run is kept as a past run'. Main at ffe1e0838, in a scratch worktree, fails the same way under -race, so it is not S-0790's change. The close-out's integration tier, without -race, passed.

## Remediation

Story S-0795 remediates this issue, created from it at 2026-11-16T00:01:08Z.
Closed 2026-11-16T01:23:39Z: S-0795: T-3392 (587b9d4) made a held run's replacement one write: keep in flai/internal/serve/agents_state.go leaves out a replaced run that is imageHeld(), so put never records the PID-0 held run among a role's past runs, and the second write, dropPast, is gone. TestAHeldRunItsStartReplacesIsNoPastRun fails on the code as it was. Under load (one CPU, GOMAXPROCS=1, 40 busy loops, -race) main's code failed this test 22 of 100 runs at container_test.go:1618; the fix passed 100 of 100, the 16 held-start container tests passed 50 counts each, and flai test race passed (T-3393's Notes).
