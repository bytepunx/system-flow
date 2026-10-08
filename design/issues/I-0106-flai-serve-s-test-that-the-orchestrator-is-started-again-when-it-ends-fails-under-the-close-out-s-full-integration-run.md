---
id: I-0106
title: flai serve's test that the orchestrator is started again when it ends fails under the close-out's full integration run
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-10-07T07:31:02Z
last_reported: 2026-10-07T07:31:02Z
updated: 2026-10-08T04:44:25Z
---

# I-0106 flai serve's test that the orchestrator is started again when it ends fails under the close-out's full integration run

## Description
flai serve's test that the orchestrator is started again when it ends fails under the close-out's full integration run

## Instances

### 2026-10-07T07:31:02Z
Story: S-0212.
S-0212's second close-out stopped at the integration tier: TestTheOrchestratorIsStartedAgainWhenItEnds (flai/internal/serve/orchestrate_test.go:177) failed with 'started again within a minute of a failure', and once more run alone straight afterwards. S-0212 does not touch flai/internal/serve. It then passed 5 of 5 with -race -count=5, 3 of 3 without -race, and in a full make integration rerun: timing-dependent, like I-0090's sibling test and S-0310.

## Remediation

Story S-0320 remediates this issue, created from it at 2026-10-07T18:59:51Z.
Closed 2026-10-08T04:44:25Z: S-0320: TestTheOrchestratorIsStartedAgainWhenItEnds no longer reads the host's clock around the retry. `orchestrateLab` gained an `at` clock, and the test's looks right after the failure, one second before `Ended + orchestrateRetry`, and at it take their time from the failed run's recorded `Ended`. The cause was a host-relative shift: `Ended` is written to the second, so the truncated fraction plus the host's delay could make the run due early. The boundary is now exact and does not depend on the host (commit 5a37eb2f).
