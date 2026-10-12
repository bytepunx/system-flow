---
id: I-0144
title: "flai check finds `board.order` outside the story at close-out"
class: efficiency
status: open
count: 1
first_reported: 2026-10-12T20:48:59Z
last_reported: 2026-10-12T20:48:59Z
updated: 2026-10-12T20:48:59Z
---

# I-0144 flai check finds `board.order` outside the story at close-out

## Description
flai check finds `board.order` outside the story at close-out

## Instances

### 2026-10-12T20:48:59Z
Story: S-0238.
flai check found outside the story:
`wip/kanban/board.md`: order lists S-0353 which is in-progress; only ready and backlog stories belong in the pull order

## Remediation
