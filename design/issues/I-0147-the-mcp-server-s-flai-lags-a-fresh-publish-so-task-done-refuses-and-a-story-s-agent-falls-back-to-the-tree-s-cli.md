---
id: I-0147
title: The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-10-13T21:48:16Z
last_reported: 2026-10-13T21:48:16Z
updated: 2026-10-13T21:48:16Z
---

# I-0147 The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI

## Description
The MCP server's flai lags a fresh publish, so task_done refuses and a story's agent falls back to the tree's CLI

## Instances

### 2026-10-13T21:48:16Z
Story: S-0358.
S-0358 started minutes after flai 1.48.0 was published; the MCP server ran the installed flai 1.47.0, whose task_done refused T-1411 under ADR-0142 as a minor version behind. Every write (task done, move, criteria, issue) then went through scripts/flai.sh, whose output carries the whole inbox, and the MCP tools stayed usable only for reads and flai test.

## Remediation
