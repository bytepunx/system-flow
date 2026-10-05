---
id: I-0078
title: "flai check finds `item.archive` outside the story at close-out"
class: efficiency
status: open
count: 3
first_reported: 2026-10-05T05:52:49Z
last_reported: 2026-10-05T08:19:13Z
updated: 2026-10-05T08:19:13Z
---

# I-0078 flai check finds `item.archive` outside the story at close-out

## Description
flai check finds `item.archive` outside the story at close-out

## Instances

### 2026-10-05T05:52:49Z
Story: S-0276.
flai check found outside the story:
`wip/kanban/stories/S-0250-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md`: S-0250 is cancelled; run flai archive

### 2026-10-05T07:00:46Z
Story: S-0217.
flai check found outside the story:
`wip/kanban/stories/S-0250-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md`: S-0250 is cancelled; run flai archive

### 2026-10-05T08:19:13Z
Story: S-0218.
flai check found outside the story:
`wip/kanban/stories/S-0250-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md`: S-0250 is cancelled; run flai archive

## Remediation

Story S-0280 remediates this issue, created from it at 2026-10-05T07:09:05Z.
