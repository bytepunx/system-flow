---
id: I-0123
title: "flai check finds `board.wip-limit` outside the story at close-out"
class: efficiency
status: closed
count: 2
first_reported: 2026-10-08T08:22:06Z
last_reported: 2026-10-08T08:34:29Z
updated: 2026-10-08T09:29:49Z
---

# I-0123 flai check finds `board.wip-limit` outside the story at close-out

## Description
flai check finds `board.wip-limit` outside the story at close-out

## Instances

### 2026-10-08T08:22:06Z
Story: S-0321.
flai check found outside the story:
`wip/kanban/board.md`: 6 stories in in-progress, limit 5

### 2026-10-08T08:34:29Z
Story: S-0339.
flai check found outside the story:
`wip/kanban/board.md`: 6 stories in in-progress, limit 5

## Remediation

Story S-0348 remediates this issue, created from it at 2026-10-08T08:37:06Z.
Closed 2026-10-08T09:29:49Z: S-0348 (ADR-0133): a check scoped to a story leaves out every board.wip-limit outside it, so a close-out neither notes nor records a column of the main checkout's board over its limit; the main checkout's own flai check still reports it. TestCheckRecordIssuesLeavesOutABoardWIPLimitOutsideTheStory reproduces S-0321's and S-0339's close-outs.
