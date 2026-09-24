---
id: T-0380
type: task
nature: improvement
title: The host panel shows serve and the MCP servers with start, stop, and restart, and one check and one upgrade for flai
status: done
parent: S-0107
owner: alex
created: 2026-09-24T01:35:33Z
updated: 2026-09-24T01:40:31Z
transitions:
  - to: ready
    at: 2026-09-24T01:37:10Z
    by: system-flow
  - to: in-progress
    at: 2026-09-24T01:37:11Z
    by: system-flow
  - to: done
    at: 2026-09-24T01:40:31Z
    by: system-flow
stream: S-0107
tags: []
touches: [flaiover/src/lib/components]
---
# T-0380 The host panel shows serve and the MCP servers with start, stop, and restart, and one check and one upgrade for flai

## Work

- A `HostProcesses` component, placed on `/host` beside the Dashboard area and styled like it. The heading is `flai host`, with the host's version and pid.
- One row for `serve` and one for `MCP`. Each row shows its state, its version, and its restarts, and the MCP row lists the projects served. Each row has Start, Stop, and Restart; a button is disabled when it would do nothing.
- One `Check for upgrade` shows whether a newer flai is available. One `Upgrade` installs it and has the host restart serve and MCP. Stopping or restarting serve, and upgrading, all end the connection that asked. So, like the Dashboard area, these treat a network error or a timeout as the work going ahead, and poll `/api/host` until it answers again.
- While the `host` action is off, the controls say how to turn it on. With no host running, the area says how to start one (`flai host start`).
- Component tests: the rows, disabled states, the check result, the reconnect after an upgrade, and the off and no-host states.

## Done when

- Component tests pass. `pnpm check`, lint, and the flaiover build are clean.

## Notes
