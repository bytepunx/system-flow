---
id: I-0143
title: "flai check finds `board.order` outside the story at close-out"
class: efficiency
status: open
count: 4
first_reported: 2026-10-12T20:46:42Z
last_reported: 2026-10-12T21:18:49Z
updated: 2026-10-12T21:18:49Z
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

### 2026-10-12T21:18:49Z
Story: S-0353.
flai check found outside the story:
`wip/kanban/board.md`: order lists S-0353 which is in-progress; only ready and backlog stories belong in the pull order

## Remediation
