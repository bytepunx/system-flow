---
id: I-0141
title: The MCP test tool of the installed flai 1.40.1 refuses main's manifest, so a story's agent runs flai test from the tree
class: efficiency
status: closed
count: 6
cost: 2m
first_reported: 2026-10-12T17:07:12Z
last_reported: 2026-10-15T01:46:39Z
updated: 2026-10-15T02:09:57Z
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

### 2026-10-14T21:01:28Z
Story: S-0305.
The MCP task_done of the installed flai 1.48.0 refused to close T-1151 because main's history has flai 1.49.1 (ADR-0142); S-0305's agent closed it with scripts/flai.sh task done from the tree.

### 2026-10-14T23:14:46Z
Story: S-0389.
the MCP task_done of the installed flai 1.48.0 refused to close T-1587 under ADR-0142, as main's history holds flai 1.49.6; closed it with scripts/flai.sh task done from the tree instead

### 2026-10-15T01:42:53Z
Story: S-0387.
The MCP task_done of the installed flai 1.48.0 refused to close T-1567 under ADR-0142, as main's history holds flai 1.50.0; closed it with scripts/flai.sh task done from the tree. This is the story that remediates the issue.

### 2026-10-15T01:46:39Z
Story: S-0382.
S-0382: the MCP task_done of the host's flai 1.48.0 refused to close T-1538 because the project's history has flai 1.50.0 (ADR-0142); closed T-1538 to T-1540 with scripts/flai.sh task done from the tree instead.

## Remediation

Story S-0387 remediates this issue, created from it at 2026-10-14T20:47:00Z.
Closed 2026-10-15T02:09:57Z: Fixed by S-0387 (ADR-0158). flai serve now checks every minute whether its flai is older than the newest flai release in each served project's history. With the shell-only auto-upgrade host action on, it has flai host install that release and restart on it. Otherwise it asks the operator on one thread per project and release, so the host no longer lags unnoticed. flai test, flai verify, and the MCP tools test and verify now say an outdated flai refused the tiers and name the upgrade, instead of asking for the manifest to be fixed. ADR-0142's task_done refusal stays: it holds only until the host catches up.
