---
id: I-0373
title: "flai check finds `board.order` outside the story at close-out"
class: efficiency
status: open
count: 1
first_reported: 2026-11-15T00:30:19Z
last_reported: 2026-11-15T00:30:19Z
updated: 2026-11-15T00:30:19Z
---

# I-0373 flai check finds `board.order` outside the story at close-out

## Description
flai check finds `board.order` outside the story at close-out

## Instances

### 2026-11-15T00:30:19Z
Story: S-0738.
flai check found outside the story:
`wip/kanban/board.md`: order lists S-0738 which is in-progress; only ready and backlog stories belong in the pull order

## Remediation
