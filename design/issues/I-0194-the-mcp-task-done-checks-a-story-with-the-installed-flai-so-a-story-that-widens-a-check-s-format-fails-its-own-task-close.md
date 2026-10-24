---
id: I-0194
title: The MCP task_done checks a story with the installed flai, so a story that widens a check's format fails its own task close
class: efficiency
status: open
count: 5
cost: 4m
first_reported: 2026-10-22T22:06:51Z
last_reported: 2026-10-24T14:47:11Z
updated: 2026-10-24T14:47:11Z
---

# I-0194 The MCP task_done checks a story with the installed flai, so a story that widens a check's format fails its own task close

## Description
The MCP task_done checks a story with the installed flai, so a story that widens a check's format fails its own task close

## Instances

### 2026-10-22T22:06:51Z
Story: S-0476.
S-0476's T-2048 added answers to fourteen examples files, a key the branch's flai accepts (ADR-0187). task_done, run by the MCP server's installed flai 1.67.13, committed, synced, and moved the task, then stopped at its check step on 17 components.examples findings, 'key "answers" is not one an example has', since that flai predates the key. The tree's scripts/flai.sh check --strict --story S-0476 passed with no finding. The close had to be confirmed by hand with the tree's CLI. A check step of task_done could run the flai the story's branch builds, or the close could say the finding is from a flai older than the branch's format.

### 2026-10-22T22:08:57Z
Story: S-0476.
S-0476's T-2051 close stopped the same way: the MCP task_done's check, by the installed flai 1.67.13, reported the 17 components.examples findings on answers again, after it had committed, synced, and moved the task; the tree's check passed.

### 2026-10-23T18:23:07Z
Story: S-0429.
T-1817's task_done ran flai check --strict with the installed flai 1.67.13, which predates S-0428 (accepted, unpublished) and its exemption of digest entries from doc.updated, so it reported 382 doc.updated findings on the seeded entries; the tree's flai check --strict --story S-0429 reports 0 errors. The commit, sync, and move had already happened; only the check and the inbox were skipped.

### 2026-10-23T18:42:07Z
Story: S-0429.
T-2089's task_done stopped again at the installed flai 1.67.13's check on the 382 digest entries (doc.updated), after its commit, sync, and move had gone through.

### 2026-10-24T14:47:11Z
Story: S-0432.
S-0432: the MCP task_done for T-1834 ran flai check --strict with the installed flai 1.68.1. That flai does not know the `curate` role S-0431 added, which is accepted but not yet published, so once the story touched design/conventions/strategic-agents.md it reported `conventions.roles` as a finding in the story. The tree's flai passed. T-1835 and T-1836 were closed with scripts/flai.sh task done instead.

## Remediation
