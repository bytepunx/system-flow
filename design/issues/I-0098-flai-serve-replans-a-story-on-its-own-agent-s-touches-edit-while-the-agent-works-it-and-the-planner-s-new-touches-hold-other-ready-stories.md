---
id: I-0098
title: flai serve replans a story on its own agent's touches edit while the agent works it, and the planner's new touches hold other ready stories
class: efficiency
status: closed
count: 1
cost: 5m
first_reported: 2026-10-06T23:28:35Z
last_reported: 2026-10-06T23:28:35Z
updated: 2026-10-14T21:03:30Z
---

# I-0098 flai serve replans a story on its own agent's touches edit while the agent works it, and the planner's new touches hold other ready stories

## Description
flai serve replans a story on its own agent's touches edit while the agent works it, and the planner's new touches hold other ready stories

## Instances

### 2026-10-06T23:28:35Z
Story: S-0300.
agent-S-0300 added template/template.yaml to T-1047's and S-0300's touches at 23:03Z; flai serve started planner-S-0300 at 23:06Z (ADR-0084, plan on an edit). It added T-1123 and four touches while the story's agent was mid-layer, including flai/cmd/plan.go and flai/internal/mcpserver/plan.go, which held S-0302 and S-0303 on flai/cmd and overlapped S-0261 on flai/internal/mcpserver. T-1123 was right; the timing cost a re-read of the story and a re-plan of what was left.

## Remediation

Story S-0305 remediates this issue, created from it at 2026-10-07T01:07:10Z.
Closed 2026-10-14T21:03:30Z: Fixed by S-0305: flai serve's replanner gives no trigger for an edit by the story's own agent, as for a planner's, and drops and logs a queued planner whose story is in progress or in review or whose agent runs, whatever queued it. The operator's flai plan is unchanged. Tests in flai/internal/serve/replan_test.go reproduce the instance.
