---
id: I-0149
title: The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses
class: efficiency
status: open
count: 1
cost: 2m
first_reported: 2026-10-15T04:41:02Z
last_reported: 2026-10-15T04:41:02Z
updated: 2026-10-15T04:41:02Z
---

# I-0149 The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses

## Description
The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses

## Instances

### 2026-10-15T04:41:02Z
Story: S-0386.
The MCP server and the installed flai are both 1.48.0, and the project's history holds flai 1.51.2, so task_done refused T-1561 under ADR-0142. ADR-0160's fix, closing with the installed flai once it has caught up, does not help while the host is not upgraded. Every task was closed with the tree's CLI, scripts/flai.sh task done.

## Remediation
