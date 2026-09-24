---
id: T-0379
type: task
nature: improvement
title: The dashboard's /api/host route reads the host and asks it to start, stop, restart, check, and upgrade
status: done
parent: S-0107
owner: alex
created: 2026-09-24T01:35:32Z
updated: 2026-09-24T01:37:11Z
transitions:
  - to: ready
    at: 2026-09-24T01:37:10Z
    by: system-flow
  - to: in-progress
    at: 2026-09-24T01:37:11Z
    by: system-flow
  - to: done
    at: 2026-09-24T01:37:11Z
    by: system-flow
stream: S-0107
tags: []
touches: [flaiover/src/routes/api/host]
---
# T-0379 The dashboard's /api/host route reads the host and asks it to start, stop, restart, check, and upgrade

## Work

- `flaiover/src/routes/api/host/+server.ts`. GET returns `host.status` (the host's pid, version, and children) with `host_enabled` taken from `project.info`'s `host_actions.host`. No flai connected, or no host running, reads as `running: false` with a reason. It is not an error.
- POST `{action, process}`: `check` runs `host.check` as a read. `start`, `stop`, and `restart` for `serve`, `mcp`, or `all`, and `upgrade`, are `host.*` writes. Each gets a timeout that fits its Docker-free work, and a bad action or process gets a 400.
- Behaviour tests against the fake repo, following the existing `/api/dashboard` tests.

## Done when

- The route's tests pass and `pnpm check` and lint are clean.

## Notes
