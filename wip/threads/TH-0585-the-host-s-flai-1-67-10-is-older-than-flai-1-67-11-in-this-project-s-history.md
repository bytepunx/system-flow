---
id: TH-0585
title: The host's flai 1.67.10 is older than flai 1.67.11 in this project's history
anchor:
  path: system-flow.yaml
status: open
participants: [flai serve]
created: 2026-10-21T14:27:31Z
updated: 2026-10-21T14:27:31Z
---

# TH-0585 The host's flai 1.67.10 is older than flai 1.67.11 in this project's history

On system-flow.yaml.

## Entries

### 2026-10-21T14:27:31Z flai serve
flai serve, and the MCP server it has flai host keep for this project, run flai 1.67.10, older than flai 1.67.11, the newest release tagged in this project's history. They may lack rules and fields the project already uses, and flai task done refuses while the minor version differs (ADR-0142).

auto-upgrade is on for this project, and flai serve first asked flai host for flai 1.67.11 at 2026-10-21T13:26:53Z, but it is still not installed; the host's last answer: flai self-upgrade: resolve flai 1.67.11: bytepunx/system-flow has no published release tagged flai/v1.67.11 (drafts and prereleases are not published); published are 1.67.10, 1.67.9, 1.67.8, 1.67.7, 1.67.6, 1.67.5, 1.67.4, 1.67.3, 1.67.2, 1.67.1 and 171 older: choose one of them, or give no version for the newest.

Upgrade it in a shell on the host with: flai host upgrade --version 1.67.11.

flai serve asks once for each release: a newer one asks again.
