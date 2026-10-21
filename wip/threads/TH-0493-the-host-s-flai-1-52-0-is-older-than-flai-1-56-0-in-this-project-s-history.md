---
id: TH-0493
title: The host's flai 1.52.0 is older than flai 1.56.0 in this project's history
anchor:
  path: system-flow.yaml
status: resolved
participants: [flai serve, alex]
created: 2026-10-17T05:08:40Z
updated: 2026-10-21T12:29:17Z
---

# TH-0493 The host's flai 1.52.0 is older than flai 1.56.0 in this project's history

On system-flow.yaml.

## Entries

### 2026-10-17T05:08:40Z flai serve
flai serve, and the MCP server it has flai host keep for this project, run flai 1.52.0, older than flai 1.56.0, the newest release tagged in this project's history. They may lack rules and fields the project already uses, and flai task done refuses while the minor version differs (ADR-0142).

Upgrade it in a shell on the host with: flai host upgrade --version 1.56.0

Or have flai serve upgrade the host whenever this project's history gets ahead of it: run flai serve enable auto-upgrade in the project (ADR-0158).

flai serve asks once for each release: a newer one asks again.

### 2026-10-21T12:29:17Z alex
Resolved.
