---
id: T-0444
type: task
nature: feature
title: Document flai serve project in the user guide, the CLI design, the channel design, and the operators' runbooks
status: done
parent: S-0121
owner: alex
created: 2026-09-26T05:26:00Z
updated: 2026-09-26T05:39:18Z
transitions:
  - to: ready
    at: 2026-09-26T05:26:18Z
    by: agent-S-0118
  - to: in-progress
    at: 2026-09-26T05:35:52Z
    by: agent-S-0118
  - to: done
    at: 2026-09-26T05:39:18Z
    by: agent-S-0118
stream: S-0121
tags: []
---

# T-0444 Document flai serve project in the user guide, the CLI design, the channel design, and the operators' runbooks

## Work

`docs/users/flai.md` (flai serve project, and how it differs from flai serve import), `design/system/flai-cli.md`, `design/system/dashboard-host-channel.md` (the `removed` notification), `design/system/flaiover-dashboard.md` (a removed hub), the operators' runbooks where a project is added or removed, and the generated flag reference (`make flai-reference`).

## Done when

- every file above says what the commands do
- `flai check --strict` and the docs tests pass

## Notes
