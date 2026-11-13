---
id: I-0370
title: scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-11-13T20:39:18Z
last_reported: 2026-11-13T20:39:18Z
updated: 2026-11-13T20:39:18Z
---

# I-0370 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged

## Description
scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged

## Instances

### 2026-11-13T20:39:18Z
Story: S-0777.
S-0777 ran scripts/flai-test.sh on Linux, as its criterion 4 asks. Every Go tier, the lint, and the template smoke passed. The smoke's repository check (flai check --strict, unscoped) then failed on four story.unaccepted warnings: S-0001, S-0002, S-0003, and S-0005 are done and archived, but their branches story/S-0001 to story/S-0005 were never merged. Each branch is still checked out in a worktree of its own (git branch shows them with +). --strict does not pass over story.unaccepted, so make flai-test.sh cannot pass on this host until the operator merges or deletes those branches and their worktrees. A close-out scopes the check to its story and is not stopped by them. S-0777 deleted nothing.

## Remediation
