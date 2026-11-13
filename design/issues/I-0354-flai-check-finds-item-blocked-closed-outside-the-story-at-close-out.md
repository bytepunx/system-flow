---
id: I-0354
title: "flai check finds `item.blocked-closed` outside the story at close-out"
class: efficiency
status: closed
count: 5
first_reported: 2026-11-12T17:07:15Z
last_reported: 2026-11-13T08:08:40Z
updated: 2026-11-13T09:05:57Z
---

# I-0354 flai check finds `item.blocked-closed` outside the story at close-out

## Description
flai check finds `item.blocked-closed` outside the story at close-out

## Instances

### 2026-11-12T17:07:15Z
Story: S-0723.
flai check found outside the story:
`wip/archive/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md`: done item still has an open blocked interval; close it at the item's close with flai unblock S-0678

### 2026-11-13T07:26:31Z
Story: S-0659.
flai check found outside the story:
`wip/archive/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md`: done item still has an open blocked interval; close it at the item's close with flai unblock S-0678

### 2026-11-13T07:27:08Z
Story: S-0762.
flai check found outside the story:
`wip/archive/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md`: done item still has an open blocked interval; close it at the item's close with flai unblock S-0678

### 2026-11-13T07:50:36Z
Story: S-0733.
flai check found outside the story:
`wip/archive/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md`: done item still has an open blocked interval; close it at the item's close with flai unblock S-0678

### 2026-11-13T08:08:40Z
Story: S-0734.
flai check found outside the story:
`wip/archive/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md`: done item still has an open blocked interval; close it at the item's close with flai unblock S-0678

## Remediation

Story S-0761 remediates this issue, created from it at 2026-11-13T07:23:41Z.
Closed 2026-11-13T09:05:57Z: S-0761: S-0678's blocked interval, the one every instance names, was closed at its close with the tree's flai unblock S-0678 in the main checkout; and a check scoped to a story now leaves out every item.blocked-closed outside it, so a close-out records none (ADR-0296, flai/internal/check/scope.go, with TestScopedCheckLeavesOutBlockedClosedOutsideTheStory and TestCheckRecordIssuesLeavesOutBlockedClosedOutsideTheStory reproducing it). The main checkout's check still reports one.
