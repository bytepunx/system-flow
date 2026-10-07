---
id: I-0108
title: flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-10-07T08:06:26Z
last_reported: 2026-10-07T08:06:26Z
updated: 2026-10-07T08:06:26Z
---

# I-0108 flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit

## Description
flai task done commits every uncommitted file in the story worktree, so closing one task of a layer sweeps its siblings' work into its commit

## Instances

### 2026-10-07T08:06:26Z
Story: S-0274.
S-0274 ran layer 3 and layer 4 tasks in one worktree (T-1109 with T-1111; T-1115 with T-1117 and T-1118). Closing T-1109, then T-1115, with flai task done ran git add -A and committed the siblings' finished files under the closing task's message; each needed a git commit --amend or reset --soft to give every task its own commit. A --paths option, or committing only the task's touches, would close one task of a layer at a time.

## Remediation
