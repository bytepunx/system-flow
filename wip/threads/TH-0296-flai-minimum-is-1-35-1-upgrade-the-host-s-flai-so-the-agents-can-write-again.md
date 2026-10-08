---
id: TH-0296
title: "flai.minimum is 1.35.1: upgrade the host's flai so the agents can write again"
anchor:
  path: system-flow.yaml
status: resolved
participants: [orchestrator, alex]
created: 2026-10-07T15:04:23Z
updated: 2026-10-08T04:31:35Z
---

# TH-0296 flai.minimum is 1.35.1: upgrade the host's flai so the agents can write again

On system-flow.yaml.

## Entries

### 2026-10-07T15:04:23Z orchestrator
Recommendation: run `flai host upgrade` once the flai/v1.35.1 binaries are built, then restart the orchestrate host action.

Publishing S-0293 as flai/v1.35.1 raised `flai.minimum` to 1.35.1, because S-0293 adds the `usage.turns` field. The flai on this host is below it: my MCP server runs 1.34.2 and the CLI on PATH is 1.34.6. Both now refuse this project, so I cannot log or act until the upgrade.

### 2026-10-08T04:31:35Z alex
Resolved.
