---
id: I-0248
title: "TestTheOperatorsStopIsOneStoppedAgentEvent's waiting-for-an-answer case fails under a loaded close-out: the ended event is missing before stopped"
class: defect
status: closed
count: 3
cost: 5m
first_reported: 2026-11-02T02:50:06Z
last_reported: 2026-11-07T05:00:15Z
updated: 2026-11-07T08:29:09Z
---

# I-0248 TestTheOperatorsStopIsOneStoppedAgentEvent's waiting-for-an-answer case fails under a loaded close-out: the ended event is missing before stopped

## Description
TestTheOperatorsStopIsOneStoppedAgentEvent's waiting-for-an-answer case fails under a loaded close-out: the ended event is missing before stopped

## Instances

### 2026-11-02T02:50:06Z
Story: S-0421.
S-0421's close-out (job S-0421-20261102T024624Z, at c264f49) stopped at the integration tier on flai/internal/serve/stop_test.go:321, "an agent waiting for an answer": event 2 was the operator's stop where the test wants the agent's ended event, waiting for an answer to TH-0001. S-0421 changes nothing under flai/internal/serve, so the order of the two events depends on timing on a loaded host. The test needs to wait for the ended event before it stops the agent.

### 2026-11-02T03:51:56Z
Story: S-0422.
S-0422's close-out at 30a6b40 stopped at integration on TestTheOperatorsStopIsOneStoppedAgentEvent/an_agent_waiting_for_an_answer (stop_test.go:321: the stopped event came where ended was wanted), with the host's load average at 38.

### 2026-11-07T05:00:15Z
Story: S-0570.
S-0570's close-out (job S-0570-20261107T045512Z, head ae84860) stopped at integration on the waiting-for-an-answer case at stop_test.go:321: the event was `stopped by the operator`, not `ended`. The branch does not touch flai/internal/serve. With the host's load average at 45, the test passed 5 of 5 alone under -race -count=5. The close-out was asked for again.

## Remediation

Story S-0561 remediates this issue, created from it at 2026-11-02T03:56:56Z.
Closed 2026-11-07T08:29:09Z: S-0561 (T-2550, 375012db1): the waiting-for-an-answer case in flai/internal/serve/stop_test.go now waits for the ended event before it stops the agent. The launcher's endRun records the run's outcome, then appends its ended event; the case stopped the agent as soon as the outcome read asked, so on a loaded host the stopped event could be appended in between. With a 300ms sleep held between those two steps, the old case failed exactly as I-0248's instances record and the new one passed 3 of 3. No product code changed: the event log is news, not state (ADR-0211 §4), and a stop in that gap only reorders two events a reader still gets.
