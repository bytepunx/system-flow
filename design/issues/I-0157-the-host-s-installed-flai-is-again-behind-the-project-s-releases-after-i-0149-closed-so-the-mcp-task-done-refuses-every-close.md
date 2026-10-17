---
id: I-0157
title: The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close
class: efficiency
status: open
count: 5
cost: 3m
first_reported: 2026-10-17T01:51:51Z
last_reported: 2026-10-17T07:57:59Z
updated: 2026-10-17T07:57:59Z
---

# I-0157 The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close

## Description
The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close

## Instances

### 2026-10-17T01:51:51Z
Story: S-0439.
S-0439, 2026-10-17: the MCP server's flai 1.52.0 and the installed ~/.flai/bin/flai it falls back to (ADR-0160) are both behind 1.55.0 in the project's history, so task_done refused T-1867. All five tasks were closed with the tree's scripts/flai.sh task done in the worktree. The upgrade threads TH-0469, TH-0485, and TH-0487 are still open. I-0147 and I-0149 recorded the same failure and are closed.

### 2026-10-17T02:21:16Z
Story: S-0440.
S-0440's agent found the installed flai at 1.52.0 against 1.55.0 in the project's history, so it closed every task with the tree's flai, scripts/flai.sh task done, instead of the MCP task_done.

### 2026-10-17T02:53:57Z
Story: S-0442.
The MCP task_done refused T-1879: the server and the installed flai are 1.52.0, the project's history 1.55.1. Closed it with the worktree's scripts/flai.sh task done instead.

### 2026-10-17T04:04:28Z
Story: S-0443.
MCP task_done refused T-1884: the server and the installed flai at ~/.flai/bin/flai are 1.52.0, a minor version behind 1.55.3 in the project's history; closed with scripts/flai.sh task done instead.

### 2026-10-17T07:57:59Z
Story: S-0393.
S-0393: the MCP task_done refused T-1610 because the installed flai 1.57.0 is behind flai 1.58.0 in the project's history. Tasks were closed with the tree's flai through scripts/flai.sh task done instead.

## Remediation
