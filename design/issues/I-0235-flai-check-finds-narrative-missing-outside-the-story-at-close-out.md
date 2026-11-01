---
id: I-0235
title: "flai check finds `narrative.missing` outside the story at close-out"
class: efficiency
status: open
count: 2
first_reported: 2026-10-31T09:53:32Z
last_reported: 2026-10-31T10:29:52Z
updated: 2026-11-01T05:22:01Z
---

# I-0235 flai check finds `narrative.missing` outside the story at close-out

## Description
flai check finds `narrative.missing` outside the story at close-out

## Instances

### 2026-10-31T09:53:32Z
Story: S-0548.
flai check found outside the story:
`wip/kanban/stories/S-0549-the-merge-queue-checks-the-main-checkout-s-uncommitted-files-and-the-lint-of.md`: S-0549 is in-progress but has no narrative; run flai stream open S-0549

### 2026-10-31T10:29:52Z
Story: S-0550.
flai check found outside the story:
`wip/kanban/stories/S-0549-the-merge-queue-checks-the-main-checkout-s-uncommitted-files-and-the-lint-of.md`: S-0549 is in-progress but has no narrative; run flai stream open S-0549

## Remediation
