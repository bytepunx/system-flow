---
id: I-0195
title: A thread with a pending orchestrator recommendation waits on the operator however long it takes, with no setting to confirm it after a threshold
class: efficiency
status: open
count: 1
first_reported: 2026-10-23T04:56:41Z
last_reported: 2026-10-23T04:56:41Z
updated: 2026-10-26T04:01:40Z
---

# I-0195 A thread with a pending orchestrator recommendation waits on the operator however long it takes, with no setting to confirm it after a threshold

## Description
A thread with a pending orchestrator recommendation waits on the operator however long it takes, with no setting to confirm it after a threshold

## Instances

### 2026-10-23T04:56:41Z
On TH-0598 (2026-10-23T04:56:05Z) the operator confirmed this as its own story: a new permission, orchestration.permissions.confirm_after, off by default, a duration such as 48h. When an orchestrator recommendation has waited that long with no operator entry, flai serve confirms it on a schedule, marked as a timeout confirmation in the thread, the activity log, and flai stats beside the operator's confirmations (ADR-0091). Never for money (cost of delay inputs, estimates, spend), host changes (installs, keys, host actions), permission prompts, .claude/ writes, or a thread the operator wrote on since the recommendation. It needs an ADR refining ADR-0090. Instances: S-0424's planner waited 30h on TH-0588 and stalled story planning, and S-0423 on TH-0583.

## Remediation

Story S-0502 remediates this issue, created from it at 2026-10-26T04:01:40Z.
