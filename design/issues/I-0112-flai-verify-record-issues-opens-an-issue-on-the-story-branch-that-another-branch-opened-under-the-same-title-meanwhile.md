---
id: I-0112
title: flai verify --record-issues opens an issue on the story branch that another branch opened under the same title meanwhile
class: defect
status: open
count: 2
cost: 8m
first_reported: 2026-10-07T09:09:31Z
last_reported: 2026-10-08T04:54:40Z
updated: 2026-10-08T04:54:40Z
---

# I-0112 flai verify --record-issues opens an issue on the story branch that another branch opened under the same title meanwhile

## Description
flai verify --record-issues opens an issue on the story branch that another branch opened under the same title meanwhile

## Instances

### 2026-10-07T09:09:31Z
Story: S-0275.
S-0275's close-out opened I-0112 for its narrative.state notes; S-0265 had opened I-0111 with the same title on main; after a sync the close-out bumped I-0111, leaving I-0112 a duplicate, removed by hand.

### 2026-10-08T04:54:40Z
Story: S-0318.
S-0318's close-out bumped I-0118 (markdown.MD034 on wip/agents/orchestrator.md) on its branch while S-0316's bump of the same issue reached main, so the close-out stopped at its sync check and flai stream sync stopped on I-0118 and summary.md. Resolved by hand, keeping both instances (count 3) and regenerating summary.md; the trial merge then found the same conflict with S-0320 (MS-0004).

## Remediation

Story S-0326 remediates this issue, created from it at 2026-10-07T18:59:58Z.
