---
id: I-0141
title: The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree
class: efficiency
status: open
count: 2
cost: 2m
first_reported: 2026-10-12T17:07:12Z
last_reported: 2026-10-14T20:40:07Z
updated: 2026-10-14T20:47:00Z
---

# I-0141 The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree

## Description
The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree

## Instances

### 2026-10-12T17:07:12Z
Story: S-0359.
S-0359 called the MCP tool test after closing T-1415 and T-1417; it refused with 'tests[1].command (tier "vet") argument {importers} is not a placeholder flai fills', because the installed flai is 1.40.1 and main's manifest uses {importers} since S-0364 (TH-0411). scripts/flai.sh test ran the same paths and passed.

### 2026-10-14T20:40:07Z
Story: S-0380.
The MCP server's installed flai, 1.48.0, refused task_done for T-1530 under ADR-0142, since the project's history has flai 1.49.1, a minor version ahead. The host had not been upgraded while S-0380 waited a day for the operator's answer. The agent closed the task with the tree's flai, `scripts/flai.sh task done`, instead.

## Remediation

Story S-0387 remediates this issue, created from it at 2026-10-14T20:47:00Z.
