---
id: I-0147
title: The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI
class: efficiency
status: open
count: 3
cost: 3m
first_reported: 2026-10-13T21:48:16Z
last_reported: 2026-10-14T21:40:37Z
updated: 2026-10-14T21:40:37Z
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

## Remediation
