---
id: I-0080
title: A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there
class: efficiency
status: open
count: 2
cost: 3m
first_reported: 2026-10-05T07:07:27Z
last_reported: 2026-10-06T19:09:21Z
updated: 2026-10-06T19:09:21Z
---

# I-0080 A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there

## Description
A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there

## Instances

### 2026-10-05T07:07:27Z
Story: S-0217.
S-0217's first close-out stopped at flaiover lint, types, and unit tests with prettier: not found; scripts/flaiover-install.sh in the worktree fixed it. Some worktrees have node_modules, others do not: stream open does not install them.

### 2026-10-06T19:09:21Z
Story: S-0295.
S-0295's worktree had no flaiover/node_modules before T-1032; make flaiover-install fixed it in 3s

## Remediation

Story S-0281 remediates this issue, created from it at 2026-10-05T07:09:06Z.
