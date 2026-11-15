---
id: I-0375
title: "flai check finds `board.order` outside the story at close-out"
class: efficiency
status: closed
count: 2
first_reported: 2026-11-15T00:30:19Z
last_reported: 2026-11-15T00:46:56Z
updated: 2026-11-15T22:08:07Z
---

# I-0375 flai check finds `board.order` outside the story at close-out

## Description
flai check finds `board.order` outside the story at close-out

## Instances

### 2026-11-15T00:30:19Z
Story: S-0738.
flai check found outside the story:
`wip/kanban/board.md`: order lists S-0738 which is in-progress; only ready and backlog stories belong in the pull order

### 2026-11-15T00:46:56Z
Story: S-0746.
flai check found outside the story:
`wip/kanban/board.md`: order lists S-0738 which is in-progress; only ready and backlog stories belong in the pull order

## Remediation

Story S-0790 remediates this issue, created from it at 2026-11-15T01:56:43Z.
Closed 2026-11-15T22:08:07Z: Fixed by S-0790. The cause: S-0738, S-0746, and S-0739 started in the same second. Each start's TransitionAll read board.md, took its story out of the order, and saved it without a lock, so one save wrote back S-0738 after another start had removed it. ADR-0306: every write of board.md, a move of a story or an epic, flai order, flai accept, and flai board limit, now takes the board lock (.flai-cache/board.lock) before the items lock, reads the board again inside it, and drops order names that are no longer active ready or backlog stories before it writes. TestStoriesStartedAtOnceAllLeaveTheOrder in flai/internal/workitem/board_update_test.go reproduces the lost update with six concurrent starts. TestUpdateBoardDropsStaleNamesFromTheOrder covers the normalisation, and TestOrderByApplyReadsTheBoardInsideTheLock covers flai order --apply.
