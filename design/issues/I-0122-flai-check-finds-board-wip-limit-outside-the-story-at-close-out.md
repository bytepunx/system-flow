---
id: I-0122
title: "flai check finds `board.wip-limit` outside the story at close-out"
class: efficiency
status: closed
count: 1
first_reported: 2026-10-08T08:30:23Z
last_reported: 2026-10-08T08:30:23Z
updated: 2026-10-12T14:21:13Z
---

# I-0122 flai check finds `board.wip-limit` outside the story at close-out

## Description
flai check finds `board.wip-limit` outside the story at close-out

## Instances

### 2026-10-08T08:30:23Z
Story: S-0291.
flai check found outside the story:
`wip/kanban/board.md`: 6 stories in in-progress, limit 5

## Remediation

Story S-0361 remediates this issue, created from it at 2026-10-12T14:12:27Z.
Closed 2026-10-12T14:21:13Z: Fixed by S-0348 (ADR-0133), released in flai 1.39.15: a check scoped to a story, as the close-out runs it, leaves out every board.wip-limit outside the story. I-0122's one instance, S-0291's close-out reporting 6 stories in in-progress over a limit of 5, came before that fix and is the "in progress over its limit" case of TestScopedCheckLeavesOutABoardWIPLimitOutsideTheStory in flai/internal/check/scope_test.go, which S-0361 names I-0122 in. Duplicate of I-0123.
