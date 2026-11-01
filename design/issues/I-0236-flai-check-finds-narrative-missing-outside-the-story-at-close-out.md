---
id: I-0236
title: "flai check finds `narrative.missing` outside the story at close-out"
class: efficiency
status: open
count: 1
first_reported: 2026-10-31T10:07:40Z
last_reported: 2026-10-31T10:07:40Z
updated: 2026-10-31T10:07:40Z
---

# I-0236 flai check finds `narrative.missing` outside the story at close-out

## Description
flai check finds `narrative.missing` outside the story at close-out

## Instances

### 2026-10-31T10:07:40Z
Story: S-0547.
flai check found outside the story:
`wip/kanban/stories/S-0549-the-merge-queue-checks-the-main-checkout-s-uncommitted-files-and-the-lint-of.md`: S-0549 is in-progress but has no narrative; run flai stream open S-0549

## Remediation
