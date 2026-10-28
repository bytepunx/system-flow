---
id: I-0183
title: "flai check finds `item.chronology` outside the story at close-out"
class: efficiency
status: open
count: 1
first_reported: 2026-10-21T12:18:49Z
last_reported: 2026-10-21T12:18:49Z
updated: 2026-10-28T12:53:36Z
---

# I-0183 flai check finds `item.chronology` outside the story at close-out

## Description
flai check finds `item.chronology` outside the story at close-out

## Instances

### 2026-10-21T12:18:49Z
Story: S-0473.
flai check found outside the story:
`wip/kanban/stories/S-0482-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-node-modules-without-it-so-every-other-close-out-fails-testeverylistedflaiovercomponentcompilesintheharness.md`: updated is before the last transition

## Remediation

Story S-0522 remediates this issue, created from it at 2026-10-28T12:53:36Z.
