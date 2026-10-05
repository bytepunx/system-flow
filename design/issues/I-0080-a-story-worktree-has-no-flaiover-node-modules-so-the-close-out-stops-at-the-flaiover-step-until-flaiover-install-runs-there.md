---
id: I-0080
title: A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-10-05T07:07:27Z
last_reported: 2026-10-05T07:07:27Z
updated: 2026-10-05T07:07:27Z
---

# I-0080 A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there

## Description
A story worktree has no flaiover/node_modules, so the close-out stops at the flaiover step until flaiover-install runs there

## Instances

### 2026-10-05T07:07:27Z
Story: S-0217.
S-0217's first close-out stopped at flaiover lint, types, and unit tests with prettier: not found; scripts/flaiover-install.sh in the worktree fixed it. Some worktrees have node_modules, others do not: stream open does not install them.

## Remediation
