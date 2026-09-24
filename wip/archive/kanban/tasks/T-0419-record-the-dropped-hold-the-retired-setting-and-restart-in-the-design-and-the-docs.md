---
id: T-0419
type: task
nature: remediation
title: Record the dropped hold, the retired setting, and restart in the design and the docs
status: done
parent: S-0116
owner: alex
created: 2026-09-24T09:00:33Z
updated: 2026-09-24T09:16:28Z
transitions:
  - to: ready
    at: 2026-09-24T09:15:03Z
    by: agent-S-0116
  - to: in-progress
    at: 2026-09-24T09:15:03Z
    by: agent-S-0116
  - to: done
    at: 2026-09-24T09:16:28Z
    by: agent-S-0116
stream: S-0116
tags: []
touches: [design/system, docs/operators, docs/users]
---
# T-0419 Record the dropped hold, the retired setting, and restart in the design and the docs

## Work

Update the launcher's section of `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `design/system/dashboard-host-channel.md`, `design/system/workflow.md`, `docs/operators/index.md`, `docs/users/flai.md`, `docs/users/flaiover.md`, and regenerate `docs/users/flai-reference.md`.

## Done when

- No document says attendance holds a ready story back or that `attended_minutes` does anything.
- Restart is documented for operators and users.
- `flai check --strict` and the markdown lint are clean.

## Notes
