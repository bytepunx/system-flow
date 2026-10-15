---
id: I-0149
title: The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses
class: efficiency
status: closed
count: 1
cost: 2m
first_reported: 2026-10-15T04:41:02Z
last_reported: 2026-10-15T04:41:02Z
updated: 2026-10-15T22:54:39Z
---

# I-0149 The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses

## Description
The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses

## Instances

### 2026-10-15T04:41:02Z
Story: S-0386.
The MCP server and the installed flai are both 1.48.0, and the project's history holds flai 1.51.2, so task_done refused T-1561 under ADR-0142. ADR-0160's fix, closing with the installed flai once it has caught up, does not help while the host is not upgraded. Every task was closed with the tree's CLI, scripts/flai.sh task done.

## Remediation

Story S-0391 remediates this issue, created from it at 2026-10-15T05:02:18Z.
Closed 2026-10-15T22:54:39Z: S-0391. The host ran flai 1.48.0, which predates ADR-0158's minute check (1.50.1) and ADR-0160's close through the installed flai (1.51.1), so nothing upgraded it or asked. The operator ran flai host upgrade once (TH-0461). At 2026-10-15T22:54Z, flai host, flai serve, and the MCP servers run flai 1.52.0, the newest release in the project's history, and MCP task_done closed T-1601 with no refusal. From now on ADR-0158 keeps the host current. With auto-upgrade off, its thread now names a release only once it is published (T-1599, ADR-0161, TestTheUpgradeThreadWaitsForTheReleaseToBePublished), so the command it gives always installs. The operator docs name the one manual upgrade that a host older than 1.50.1 needs (T-1600).
