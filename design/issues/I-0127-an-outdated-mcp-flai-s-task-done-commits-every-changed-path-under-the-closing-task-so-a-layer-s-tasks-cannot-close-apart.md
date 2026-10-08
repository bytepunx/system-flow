---
id: I-0127
title: An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-10-08T10:41:45Z
last_reported: 2026-10-08T10:41:45Z
updated: 2026-10-08T10:41:45Z
---

# I-0127 An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart

## Description
An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart

## Instances

### 2026-10-08T10:41:45Z
Story: S-0338.
S-0338 ran three layer-1 tasks together (T-1324, T-1325, T-1326, no path in common). The MCP server runs the installed flai 1.39.13, older than the tree's 1.40.0 (inbox says flai_outdated). Its task_done for T-1324 committed every changed path in the worktree, the two other tasks' flaiover files included, in one commit (9f67200a), and widened T-1324's touches to them. ADR-0128 says task_done leaves another open task's paths for that task's close. T-1324 was done by then and refuses an edit, so its touches stay too wide; T-1325 and T-1326 closed with nothing to commit. Fix: upgrade the host's flai (flai host upgrade), or have task_done refuse, naming the upgrade, while the running flai is older than the project's history.

## Remediation
