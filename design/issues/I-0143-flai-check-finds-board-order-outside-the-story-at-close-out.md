---
id: I-0143
title: "flai check finds `board.order` outside the story at close-out"
class: efficiency
status: closed
count: 5
first_reported: 2026-10-12T20:46:42Z
last_reported: 2026-10-12T21:18:49Z
updated: 2026-10-21T07:24:54Z
---

# I-0143 flai check finds `board.order` outside the story at close-out

## Description
flai check finds `board.order` outside the story at close-out

## Instances

### 2026-10-12T20:46:42Z
Story: S-0370.
flai check found outside the story:
`wip/kanban/board.md`: order lists S-0353 which is in-progress; only ready and backlog stories belong in the pull order

### 2026-10-12T20:48:59Z
Story: S-0238.
flai check found outside the story:
`wip/kanban/board.md`: order lists S-0353 which is in-progress; only ready and backlog stories belong in the pull order

### 2026-10-12T20:55:37Z
Story: S-0360.
flai check found outside the story:
`wip/kanban/board.md`: order lists S-0353 which is in-progress; only ready and backlog stories belong in the pull order

### 2026-10-12T21:09:35Z
Story: S-0352.
flai check found outside the story:
`wip/kanban/board.md`: order lists S-0353 which is in-progress; only ready and backlog stories belong in the pull order

### 2026-10-12T21:18:49Z
Story: S-0353.
flai check found outside the story:
`wip/kanban/board.md`: order lists S-0353 which is in-progress; only ready and backlog stories belong in the pull order

## Remediation
Closed 2026-10-21T07:24:54Z: Not seen since 2026-10-12: all five reports are one cause on one day, S-0353 left in board.order while in progress. Every flai move out of ready or backlog now drops the story from the order (flai/internal/workitem/rules.go, RemoveFromOrder), board.md's order is empty, and 69 stories have started since with no recurrence. Reopen it with a new issue if flai check finds board.order again.
