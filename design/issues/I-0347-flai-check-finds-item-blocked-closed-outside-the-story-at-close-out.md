---
id: I-0347
title: "flai check finds `item.blocked-closed` outside the story at close-out"
class: efficiency
status: open
count: 4
first_reported: 2026-11-11T07:11:49Z
last_reported: 2026-11-11T08:08:57Z
updated: 2026-11-11T08:08:57Z
---

# I-0347 flai check finds `item.blocked-closed` outside the story at close-out

## Description
flai check finds `item.blocked-closed` outside the story at close-out

## Instances

### 2026-11-11T07:11:49Z
Story: S-0722.
flai check found outside the story:
`wip/archive/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md`: done item still has an open blocked interval

### 2026-11-11T07:50:11Z
Story: S-0691.
flai check found outside the story:
`wip/archive/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md`: done item still has an open blocked interval

### 2026-11-11T07:51:03Z
Story: S-0686.
flai check found outside the story:
`wip/archive/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md`: done item still has an open blocked interval

### 2026-11-11T08:08:57Z
Story: S-0697.
flai check found outside the story:
`wip/archive/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md`: done item still has an open blocked interval

## Remediation

Story S-0723 remediates this issue, created from it at 2026-11-11T07:23:14Z.
