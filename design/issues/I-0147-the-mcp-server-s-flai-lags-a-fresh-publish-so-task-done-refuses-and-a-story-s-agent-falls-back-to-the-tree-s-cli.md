---
id: I-0147
title: The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI
class: efficiency
status: closed
count: 5
cost: 3m
first_reported: 2026-10-13T21:48:16Z
last_reported: 2026-10-15T03:20:04Z
updated: 2026-10-15T03:20:12Z
---

# I-0147 The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI

## Description
The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI

## Instances

### 2026-10-13T21:48:16Z
Story: S-0358.
S-0358 started minutes after flai 1.48.0 was published; the MCP server ran the installed flai 1.47.0, whose task_done refused T-1411 under ADR-0142 as a minor version behind. Every write (task done, move, criteria, issue) then went through scripts/flai.sh, whose output carries the whole inbox, and the MCP tools stayed usable only for reads and flai test.

### 2026-10-14T21:17:19Z
Story: S-0304.
task_done refused T-1572 because the MCP server's flai was 1.48.0 and the project's history had flai 1.49.2 (published minutes earlier with S-0305's acceptance); closed it with scripts/flai.sh task done from the worktree instead.

### 2026-10-14T21:40:37Z
Story: S-0372.
The MCP server's flai is 1.48.0 and the project's history has flai 1.49.1, so task_done refused to close T-1482 (ADR-0142). Every task of S-0372 was closed with scripts/flai.sh task done from the worktree instead.

### 2026-10-15T02:50:44Z
Story: S-0385.
S-0385: the installed flai 1.48.0 behind the MCP server refused task_done for T-1556 (ADR-0142, 1.50.1 in the project's history); every task was closed with the tree's flai, ../../../scripts/flai.sh task done, from the worktree instead.

### 2026-10-15T02:52:47Z
Story: S-0383.
The MCP task_done of the installed flai 1.48.0 refused to close T-1544 (ADR-0142), because the project's history has flai 1.50.0. Every task of S-0383 was closed with the tree's scripts/flai.sh task done instead.

### 2026-10-15T03:20:04Z
Story: S-0390.
S-0390 started on the installed flai 1.48.0 while the project's history had flai 1.50.1; MCP task_done refused T-1595 (ADR-0142), so every task of S-0390 was closed with the tree's scripts/flai.sh task done.

## Remediation

Story S-0390 remediates this issue, created from it at 2026-10-15T02:55:10Z.
Closed 2026-10-15T03:20:12Z: S-0390 (ADR-0160): when ADR-0142 refuses an outdated MCP server's task_done, it now closes the task with the flai installed at the server's own path once that flai has caught up, answering closed_by; TestTaskDoneClosesWithTheInstalledFlaiOnceItHasCaughtUp reproduces the instances.
