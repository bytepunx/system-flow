---
id: TH-1272
title: The host's flai 1.101.4 is older than flai 1.102.0 in this project's history
anchor:
  path: system-flow.yaml
status: open
participants: [flai serve]
created: 2026-11-17T04:28:35Z
updated: 2026-11-17T04:28:35Z
---

# TH-1272 The host's flai 1.101.4 is older than flai 1.102.0 in this project's history

On system-flow.yaml.

## Entries

### 2026-11-17T04:28:35Z flai serve
flai serve, and the MCP server it has flai host keep for this project, run flai 1.101.4, older than flai 1.102.0, the newest release tagged in this project's history. They may lack rules and fields the project already uses, and flai task done refuses while the minor version differs (ADR-0142).

auto-upgrade is on for this project, and flai serve first asked flai host for flai 1.102.0 at 2026-11-17T03:27:34Z, but it is still not installed; the host's last answer: flai self-upgrade: resolve flai 1.102.0: bytepunx/system-flow has no published release tagged flai/v1.102.0 (drafts and prereleases are not published); published are 1.101.4, 1.101.3, 1.101.2, 1.101.0, 1.100.0, 1.99.0, 1.98.0, 1.97.0, 1.96.0, 1.95.6 and 240 older: choose one of them, or give no version for the newest.

Upgrade it in a shell on the host with: flai host upgrade --version 1.102.0.

flai serve asks once for each release: a newer one asks again.
