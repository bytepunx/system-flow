---
id: T-0651
type: task
nature: feature
title: "Design and operator guide describe the watchdog and the refused connection; close I-0025 and I-0029"
status: in-progress
parent: S-0184
owner: arobson
created: 2026-10-01T08:55:54Z
updated: 2026-10-01T09:06:58Z
transitions:
  - to: ready
    at: 2026-10-01T08:56:22Z
    by: agent-S-0184
  - to: in-progress
    at: 2026-10-01T09:06:58Z
    by: agent-S-0184
stream: S-0184
tags: []
---

# T-0651 Design and operator guide describe the watchdog and the refused connection; close I-0025 and I-0029

## Work

- `design/system/flaiover-dashboard.md` and `design/system/dashboard-host-channel.md`: the HEALTHCHECK, the host's watch, and the refused connection with its close code.
- `docs/operators/` (settings index, guide) and `docs/users/flai.md`: `flai dashboard status`'s states and `dashboard.no_restart`.
- `flai issue close I-0025` and `I-0029` with what fixed them.

## Done when

- The documents describe the behaviour as built, `flai check --strict` is clean, and both issues are closed.

## Notes
