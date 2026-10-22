---
id: I-0194
title: The MCP task_done checks a story with the installed flai, so a story that widens a check's format fails its own task close
class: efficiency
status: open
count: 2
cost: 3m
first_reported: 2026-10-22T22:06:51Z
last_reported: 2026-10-22T22:08:57Z
updated: 2026-10-22T22:08:57Z
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

## Remediation
