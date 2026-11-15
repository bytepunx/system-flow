---
id: I-0375
title: "flai check finds `board.order` outside the story at close-out"
class: efficiency
status: open
count: 1
first_reported: 2026-11-15T00:46:56Z
last_reported: 2026-11-15T00:46:56Z
updated: 2026-11-15T01:56:43Z
---

# I-0375 flai check finds `board.order` outside the story at close-out

## Description
flai check finds `board.order` outside the story at close-out

## Instances

### 2026-11-15T00:46:56Z
Story: S-0746.
flai check found outside the story:
`wip/kanban/board.md`: order lists S-0738 which is in-progress; only ready and backlog stories belong in the pull order

## Remediation

Story S-0790 remediates this issue, created from it at 2026-11-15T01:56:43Z.
