---
id: I-0157
title: The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close
class: efficiency
status: open
count: 1
cost: 3m
first_reported: 2026-10-17T01:51:51Z
last_reported: 2026-10-17T01:51:51Z
updated: 2026-10-17T01:51:51Z
---

# I-0157 The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close

## Description
The host's installed flai is again behind the project's releases after I-0149 closed, so the MCP task_done refuses every close

## Instances

### 2026-10-17T01:51:51Z
Story: S-0439.
S-0439, 2026-10-17: the MCP server's flai 1.52.0 and the installed ~/.flai/bin/flai it falls back to (ADR-0160) are both behind 1.55.0 in the project's history, so task_done refused T-1867. All five tasks were closed with the tree's scripts/flai.sh task done in the worktree. The upgrade threads TH-0469, TH-0485, and TH-0487 are still open. I-0147 and I-0149 recorded the same failure and are closed.

## Remediation
