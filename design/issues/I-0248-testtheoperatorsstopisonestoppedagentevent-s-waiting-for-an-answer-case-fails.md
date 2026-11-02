---
id: I-0248
title: "TestTheOperatorsStopIsOneStoppedAgentEvent's waiting-for-an-answer case fails under a loaded close-out: the ended event is missing before stopped"
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-11-02T02:50:06Z
last_reported: 2026-11-02T02:50:06Z
updated: 2026-11-02T03:56:56Z
---

# I-0248 TestTheOperatorsStopIsOneStoppedAgentEvent's waiting-for-an-answer case fails under a loaded close-out: the ended event is missing before stopped

## Description
TestTheOperatorsStopIsOneStoppedAgentEvent's waiting-for-an-answer case fails under a loaded close-out: the ended event is missing before stopped

## Instances

### 2026-11-02T02:50:06Z
Story: S-0421.
S-0421's close-out (job S-0421-20261102T024624Z, at c264f49) stopped at the integration tier on flai/internal/serve/stop_test.go:321, "an agent waiting for an answer": event 2 was the operator's stop where the test wants the agent's ended event, waiting for an answer to TH-0001. S-0421 changes nothing under flai/internal/serve, so the order of the two events depends on timing on a loaded host. The test needs to wait for the ended event before it stops the agent.

## Remediation

Story S-0561 remediates this issue, created from it at 2026-11-02T03:56:56Z.
