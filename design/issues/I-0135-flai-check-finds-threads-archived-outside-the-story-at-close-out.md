---
id: I-0135
title: "flai check finds `threads.archived` outside the story at close-out"
class: efficiency
status: open
count: 2
first_reported: 2026-10-09T18:33:15Z
last_reported: 2026-10-12T14:23:39Z
updated: 2026-10-12T14:23:39Z
---

# I-0135 flai check finds `threads.archived` outside the story at close-out

## Description
flai check finds `threads.archived` outside the story at close-out

## Instances

### 2026-10-09T18:33:15Z
Story: S-0289.
flai check found outside the story:
`wip/threads/TH-0389-s-0338-s-acceptance-stopped-before-its-commit-s-0337-s-narrative-worktree-and-branch-are-empty-files.md`: TH-0389 is open but S-0338 is archived; resolve it with flai thread resolve TH-0389

### 2026-10-12T14:23:39Z
Story: S-0342.
flai check found TH-0389's threads.archived outside the story at S-0342's close-out on 2026-10-08T21:20:32Z, first recorded as I-0131 on the story's branch and folded in here when the branch was rebased onto main

## Remediation

Story S-0366 remediates this issue, created from it at 2026-10-12T14:12:35Z.
