---
id: TH-0754
title: The host's flai 1.77.0 is older than flai 1.77.1 in this project's history
anchor:
  path: system-flow.yaml
status: open
participants: [flai serve]
created: 2026-10-29T16:54:48Z
updated: 2026-10-29T16:54:48Z
---

# TH-0754 The host's flai 1.77.0 is older than flai 1.77.1 in this project's history

On system-flow.yaml.

## Entries

### 2026-10-29T16:54:48Z flai serve
flai serve, and the MCP server it has flai host keep for this project, run flai 1.77.0, older than flai 1.77.1, the newest release tagged in this project's history. They may lack rules and fields the project already uses, and flai task done refuses while the minor version differs (ADR-0142).

auto-upgrade is on for this project, and flai serve first asked flai host for flai 1.77.1 at 2026-10-29T15:53:51Z, but it is still not installed; the host's last answer: flai self-upgrade: resolve flai 1.77.1: bytepunx/system-flow has no published release tagged flai/v1.77.1 (drafts and prereleases are not published); published are 1.77.0, 1.76.1, 1.75.2, 1.74.1, 1.74.0, 1.73.0, 1.72.2, 1.72.1, 1.72.0, 1.71.5 and 191 older: choose one of them, or give no version for the newest.

Upgrade it in a shell on the host with: flai host upgrade --version 1.77.1.

flai serve asks once for each release: a newer one asks again.
