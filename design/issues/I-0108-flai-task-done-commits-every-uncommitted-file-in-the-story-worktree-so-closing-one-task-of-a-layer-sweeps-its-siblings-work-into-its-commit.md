---
id: I-0108
title: flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit
class: efficiency
status: open
count: 2
cost: 4m
first_reported: 2026-10-07T08:06:26Z
last_reported: 2026-10-08T06:02:11Z
updated: 2026-10-08T06:02:11Z
---

# I-0108 flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit

## Description
flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit

## Instances

### 2026-10-07T08:06:26Z
Story: S-0274.
S-0274 ran layer 3 and layer 4 tasks in one worktree (T-1109 with T-1111; T-1115 with T-1117 and T-1118). Closing T-1109, then T-1115, with flai task done ran git add -A and committed the siblings' finished files under the closing task's message; each needed a git commit --amend or reset --soft to give every task its own commit. A --paths option, or committing only the task's touches, would close one task of a layer at a time.

### 2026-10-08T06:02:11Z
Story: S-0309.
S-0309's layer 2 ran T-1275, T-1276, and T-1277 together in the story's worktree. Closing T-1275 with flai task done committed all three tasks' edits in one commit (1df150d5). T-1276 and T-1277 then closed with nothing to commit, and their log entries had to name T-1275's commit.

## Remediation

Story S-0322 remediates this issue, created from it at 2026-10-07T18:59:53Z.
