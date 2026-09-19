---
id: I-0021
title: flai accept merged a story's branch and removed its worktree before refusing it for unticked criteria or open tasks
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-09-19T04:19:48Z
last_reported: 2026-09-19T04:19:48Z
updated: 2026-09-19T04:19:48Z
---

# I-0021 flai accept merged a story's branch and removed its worktree before refusing it for unticked criteria or open tasks

## Description
flai accept merged a story's branch and removed its worktree before refusing it for unticked criteria or open tasks

## Instances

### 2026-09-19T04:19:48Z
S-0041: found verifying the review page in a scratch container. Accepting a story with an unticked acceptance criterion streamed 'merged' and then failed with 'S-0002 has unchecked acceptance criteria': the story stayed in review, but its commit was already on main, its worktree gone, and its branch deleted. The preflight checked identity and the worktree, not the workflow's rules for done, which ran after the merge. Not seen on this repository because agents tick criteria before review. Fixed in S-0041: the rules are checked on a copy of the item before anything changes and appear as a blocker in the preview; TestAcceptChecksTheRulesBeforeItMergesAnything reproduces it.

## Remediation
