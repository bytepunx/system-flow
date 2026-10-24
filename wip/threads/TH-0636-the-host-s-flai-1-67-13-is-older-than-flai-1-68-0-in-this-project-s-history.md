---
id: TH-0636
title: The host's flai 1.67.13 is older than flai 1.68.0 in this project's history
anchor:
  path: system-flow.yaml
status: open
participants: [flai serve]
created: 2026-10-24T00:41:22Z
updated: 2026-10-24T00:41:22Z
---

# TH-0636 The host's flai 1.67.13 is older than flai 1.68.0 in this project's history

On system-flow.yaml.

## Entries

### 2026-10-24T00:41:22Z flai serve
flai serve, and the MCP server it has flai host keep for this project, run flai 1.67.13, older than flai 1.68.0, the newest release tagged in this project's history. They may lack rules and fields the project already uses, and flai task done refuses while the minor version differs (ADR-0142).

auto-upgrade is on for this project, and flai serve first asked flai host for flai 1.68.0 at 2026-10-23T23:40:49Z, but it is still not installed; the host's last answer: flai self-upgrade: resolve flai 1.68.0: bytepunx/system-flow has no published release tagged flai/v1.68.0 (drafts and prereleases are not published); published are 1.67.13, 1.67.12, 1.67.10, 1.67.9, 1.67.8, 1.67.7, 1.67.6, 1.67.5, 1.67.4, 1.67.3 and 173 older: choose one of them, or give no version for the newest.

Upgrade it in a shell on the host with: flai host upgrade --version 1.68.0.

flai serve asks once for each release: a newer one asks again.
