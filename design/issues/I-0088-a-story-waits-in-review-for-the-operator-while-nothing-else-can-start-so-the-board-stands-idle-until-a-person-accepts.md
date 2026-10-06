---
id: I-0088
title: A story waits in review for the operator while nothing else can start, so the board stands idle until a person accepts
class: blocker
status: open
count: 2
cost: 1h24m
first_reported: 2026-10-06T10:32:27Z
last_reported: 2026-10-06T10:32:27Z
updated: 2026-10-06T11:44:51Z
---

# I-0088 A story waits in review for the operator while nothing else can start, so the board stands idle until a person accepts

## Description

Only the operator accepts a story. While it waits in review it keeps its claim, and when every ready story overlaps it (I-0087) nothing else can start. The board then does no work until a person acts, however long that is. The agents were finished and idle; the wait was for a person.

## Instances

### 2026-10-06T10:32:27Z
Story: S-0220.
S-0220 moved to review at 07:11:57Z with every criterion ticked and its close-out passed. The operator accepted it at about 09:56:50Z. For those 2h45m nothing was in progress, no agent ran, and ten ready stories were held behind it. The story's own work took about 1h13m.

Story: S-0283.
S-0283 moved to review at 10:30:28Z and the operator, at the dashboard, accepted it at about 10:31:58Z. Nine ready stories waited 90 seconds. The wait is as long as the operator is away.

## Remediation

S-0221, already in ready, lets the orchestrator accept a story in review when the operator turns that permission on. S-0286 keeps the operator's acceptance for a story that changes a path Claude Code protects. Between them, the wait would remain only for those stories and for the ones the operator chooses to review. Until S-0221 is built, an agent in the operator's session can run the acceptance on the operator's word.

Story S-0296 remediates this issue, created from it at 2026-10-06T11:44:51Z.
