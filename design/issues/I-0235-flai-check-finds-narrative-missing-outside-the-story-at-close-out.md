---
id: I-0235
title: "flai check finds `narrative.missing` outside the story at close-out"
class: efficiency
status: closed
count: 7
first_reported: 2026-10-31T09:53:32Z
last_reported: 2026-11-01T10:59:26Z
updated: 2026-11-07T04:14:53Z
---

# I-0235 flai check finds `narrative.missing` outside the story at close-out

## Description
flai check finds `narrative.missing` outside the story at close-out

## Instances

### 2026-10-31T09:53:32Z
Story: S-0548.
flai check found outside the story:
`wip/kanban/stories/S-0549-the-merge-queue-checks-the-main-checkout-s-uncommitted-files-and-the-lint-of.md`: S-0549 is in-progress but has no narrative; run flai stream open S-0549

### 2026-10-31T10:07:40Z
Story: S-0547.
flai check found outside the story:
`wip/kanban/stories/S-0549-the-merge-queue-checks-the-main-checkout-s-uncommitted-files-and-the-lint-of.md`: S-0549 is in-progress but has no narrative; run flai stream open S-0549

### 2026-10-31T10:29:52Z
Story: S-0550.
flai check found outside the story:
`wip/kanban/stories/S-0549-the-merge-queue-checks-the-main-checkout-s-uncommitted-files-and-the-lint-of.md`: S-0549 is in-progress but has no narrative; run flai stream open S-0549

### 2026-11-01T06:00:39Z
Story: S-0546.
flai check found outside the story:
`wip/kanban/stories/S-0549-the-merge-queue-checks-the-main-checkout-s-uncommitted-files-and-the-lint-of.md`: S-0549 is in-progress but has no narrative; run flai stream open S-0549

### 2026-11-01T10:59:26Z
Story: S-0551.
flai check found outside the story:
`wip/kanban/stories/S-0549-the-merge-queue-checks-the-main-checkout-s-uncommitted-files-and-the-lint-of.md`: S-0549 is in-progress but has no narrative; run flai stream open S-0549

## Remediation

Story S-0555 remediates this issue, created from it at 2026-11-01T05:45:09Z.
Closed 2026-11-07T04:14:53Z: S-0555: ScopeToStories (flai/internal/check/scope.go) now leaves a narrative.missing outside the story out of a check scoped to a story (ADR-0256), so a close-out neither notes it nor records it here. TestScopedCheckLeavesOutANarrativeMissingOutsideTheStory and TestCheckRecordIssuesLeavesOutANarrativeMissingOutsideTheStory reproduce the S-0549 case: another story in progress with no narrative. A story's own missing narrative stays inside its check, an error, and flai check in the main checkout still reports every story in progress or in review with no narrative.
