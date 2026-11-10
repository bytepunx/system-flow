---
id: TH-0970
title: The host's flai 1.86.6 is older than flai 1.86.8 in this project's history
anchor:
  path: system-flow.yaml
status: open
participants: [flai serve]
created: 2026-11-10T06:57:14Z
updated: 2026-11-10T06:57:14Z
---

# TH-0970 The host's flai 1.86.6 is older than flai 1.86.8 in this project's history

On system-flow.yaml.

## Entries

### 2026-11-10T06:57:14Z flai serve
flai serve, and the MCP server it has flai host keep for this project, run flai 1.86.6, older than flai 1.86.8, the newest release tagged in this project's history. They may lack rules and fields the project already uses, and flai task done refuses while the minor version differs (ADR-0142).

auto-upgrade is on for this project, and flai serve first asked flai host for flai 1.86.8 at 2026-11-10T05:56:33Z, but it is still not installed; the host's last answer: flai self-upgrade: resolve flai 1.86.8: bytepunx/system-flow has no published release tagged flai/v1.86.8 (drafts and prereleases are not published); published are 1.86.6, 1.86.5, 1.86.4, 1.86.3, 1.86.2, 1.86.1, 1.86.0, 1.85.1, 1.85.0, 1.84.1 and 219 older: choose one of them, or give no version for the newest.

Upgrade it in a shell on the host with: flai host upgrade --version 1.86.8.

flai serve asks once for each release: a newer one asks again.
