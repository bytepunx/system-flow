---
id: I-0141
title: The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree
class: efficiency
status: open
count: 1
cost: 2m
first_reported: 2026-10-12T17:07:12Z
last_reported: 2026-10-12T17:07:12Z
updated: 2026-10-12T17:07:12Z
---

# I-0141 The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree

## Description
The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree

## Instances

### 2026-10-12T17:07:12Z
Story: S-0359.
S-0359 called the MCP tool test after closing T-1415 and T-1417; it refused with 'tests[1].command (tier "vet") argument {importers} is not a placeholder flai fills', because the installed flai is 1.40.1 and main's manifest uses {importers} since S-0364 (TH-0411). scripts/flai.sh test ran the same paths and passed.

## Remediation
