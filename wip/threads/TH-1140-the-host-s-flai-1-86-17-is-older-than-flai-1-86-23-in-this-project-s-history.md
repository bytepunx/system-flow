---
id: TH-1140
title: The host's flai 1.86.17 is older than flai 1.86.23 in this project's history
anchor:
  path: system-flow.yaml
status: resolved
participants: [flai serve, alex]
created: 2026-11-13T08:28:23Z
updated: 2026-11-13T20:05:30Z
---

# TH-1140 The host's flai 1.86.17 is older than flai 1.86.23 in this project's history

On system-flow.yaml.

## Entries

### 2026-11-13T08:28:23Z flai serve
flai serve, and the MCP server it has flai host keep for this project, run flai 1.86.17, older than flai 1.86.23, the newest release tagged in this project's history. They may lack rules and fields the project already uses, and flai task done refuses while the minor version differs (ADR-0142).

auto-upgrade is on for this project, and flai serve first asked flai host for flai 1.86.23 at 2026-11-13T07:27:27Z, but it is still not installed; the host's last answer: flai self-upgrade: resolve flai 1.86.23: bytepunx/system-flow has no published release tagged flai/v1.86.23 (drafts and prereleases are not published); published are 1.86.17, 1.86.16, 1.86.15, 1.86.14, 1.86.12, 1.86.11, 1.86.6, 1.86.5, 1.86.4, 1.86.3 and 225 older: choose one of them, or give no version for the newest.

Upgrade it in a shell on the host with: flai host upgrade --version 1.86.23.

flai serve asks once for each release: a newer one asks again.

### 2026-11-13T20:05:30Z alex
Resolved.
