---
id: I-0076
title: "flai check finds `wip.overlap` outside the story at close-out"
class: efficiency
status: open
count: 2
first_reported: 2026-10-05T03:24:33Z
last_reported: 2026-10-05T03:24:33Z
updated: 2026-10-05T03:24:33Z
---

# I-0076 flai check finds `wip.overlap` outside the story at close-out

## Description
flai check finds `wip.overlap` outside the story at close-out

## Instances

### 2026-10-05T03:24:33Z
Story: S-0260.
flai check found outside the story:
`wip/kanban/stories/S-0253-an-acceptance-merge-committed-conflict-markers-to-a-design-document-on-main-and-nothing-caught-it.md`: S-0253 touches docs/users/flai.md, which S-0260 (in progress) also touches as docs/users/flai.md
### 2026-10-05T03:23:20Z
Story: S-0253.
flai check found outside the story:
`wip/kanban/stories/S-0253-an-acceptance-merge-committed-conflict-markers-to-a-design-document-on-main-and-nothing-caught-it.md`: S-0253 touches docs/users/flai.md, which S-0260 (in progress) also touches as docs/users/flai.md


## Remediation
