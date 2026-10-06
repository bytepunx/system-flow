---
id: I-0090
title: flai serve's orchestrator test that waits for the stopped run's activity fails under the close-out's full parallel test run
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-10-06T11:05:36Z
last_reported: 2026-10-06T11:05:36Z
updated: 2026-10-06T11:35:28Z
---

# I-0090 flai serve's orchestrator test that waits for the stopped run's activity fails under the close-out's full parallel test run

## Description
flai serve's orchestrator test that waits for the stopped run's activity fails under the close-out's full parallel test run

## Instances

### 2026-10-06T11:05:36Z
Story: S-0285.
S-0285's close-out stopped at the full Go tests: TestTheOrchestratorIsStoppedWhenTheActionIsTurnedOff/a_run_this_flai_serve_waits_for (flai/internal/serve/orchestrate_test.go) failed with 'never happened: its activity logged' after 5.07s. S-0285 does not touch flai/internal/serve; run alone the test passed three times with -race -count=3, and the package passed on its own. The close-out had to be run again.

## Remediation

Story S-0292 remediates this issue, created from it at 2026-10-06T11:14:34Z.
Closed 2026-10-06T11:35:28Z: S-0292: the orchestrator stub in flai/internal/serve/orchestrate_test.go wrote what it was given before it printed its stream, so under load the stop's TERM could land before its usage reached the run's log, and no activity was logged. The stub now prints its stream first, so given returning means the usage is in the log. With a 0.5 s delay before the stream the test failed in the old order with the same message and passes in the new one.
