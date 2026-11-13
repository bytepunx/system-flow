---
id: I-0354
title: "flai check finds `item.blocked-closed` outside the story at close-out"
class: efficiency
status: open
count: 1
first_reported: 2026-11-12T17:07:15Z
last_reported: 2026-11-12T17:07:15Z
updated: 2026-11-13T07:23:41Z
---

# I-0354 flai check finds `item.blocked-closed` outside the story at close-out

## Description
flai check finds `item.blocked-closed` outside the story at close-out

## Instances

### 2026-11-12T17:07:15Z
Story: S-0723.
flai check found outside the story:
`wip/archive/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md`: done item still has an open blocked interval; close it at the item's close with flai unblock S-0678

## Remediation

Story S-0761 remediates this issue, created from it at 2026-11-13T07:23:41Z.
