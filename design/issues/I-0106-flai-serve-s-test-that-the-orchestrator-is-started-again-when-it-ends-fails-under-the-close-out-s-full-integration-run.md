---
id: I-0106
title: flai serve's test that the orchestrator is started again when it ends fails under the close-out's full integration run
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-07T07:31:02Z
last_reported: 2026-10-07T07:31:02Z
updated: 2026-10-07T18:59:51Z
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
