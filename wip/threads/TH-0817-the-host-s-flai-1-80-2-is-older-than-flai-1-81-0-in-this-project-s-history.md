---
id: TH-0817
title: The host's flai 1.80.2 is older than flai 1.81.0 in this project's history
anchor:
  path: system-flow.yaml
status: resolved
participants: [flai serve, alex]
created: 2026-11-02T22:53:16Z
updated: 2026-11-10T09:22:55Z
---

# TH-0817 The host's flai 1.80.2 is older than flai 1.81.0 in this project's history

On system-flow.yaml.

## Entries

### 2026-11-02T22:53:16Z flai serve
flai serve, and the MCP server it has flai host keep for this project, run flai 1.80.2, older than flai 1.81.0, the newest release tagged in this project's history. They may lack rules and fields the project already uses, and flai task done refuses while the minor version differs (ADR-0142).

auto-upgrade is on for this project, and flai serve first asked flai host for flai 1.81.0 at 2026-11-02T21:52:15Z, but it is still not installed; the host's last answer: flai self-upgrade: resolve flai 1.81.0: bytepunx/system-flow has no published release tagged flai/v1.81.0 (drafts and prereleases are not published); published are 1.80.2, 1.80.1, 1.80.0, 1.79.0, 1.78.13, 1.78.12, 1.78.11, 1.78.6, 1.78.5, 1.78.4 and 206 older: choose one of them, or give no version for the newest.

Upgrade it in a shell on the host with: flai host upgrade --version 1.81.0.

flai serve asks once for each release: a newer one asks again.

### 2026-11-10T09:22:55Z alex
Resolved.
