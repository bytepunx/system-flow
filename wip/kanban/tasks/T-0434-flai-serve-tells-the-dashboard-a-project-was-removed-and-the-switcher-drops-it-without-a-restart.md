---
id: T-0434
type: task
nature: feature
title: flai serve tells the dashboard a project was removed, and the switcher drops it without a restart
status: done
parent: S-0118
owner: alex
created: 2026-09-26T05:25:59Z
updated: 2026-09-26T05:30:23Z
transitions:
  - to: ready
    at: 2026-09-26T05:26:18Z
    by: agent-S-0118
  - to: in-progress
    at: 2026-09-26T05:28:34Z
    by: agent-S-0118
  - to: done
    at: 2026-09-26T05:30:23Z
    by: agent-S-0118
stream: S-0118
tags: []
---

# T-0434 flai serve tells the dashboard a project was removed, and the switcher drops it without a restart

## Work

Before `serve.Run` drops the client of a project that left the registry, it sends the notification `removed` with the project's key. flaiover's `AgentHub` marks itself removed on it, and the registry's `list()` leaves a removed hub out while it is not connected, as it does a gone candidate; a new connection for the key clears the mark. A dashboard older than this ignores the notification; a flai older than this never sends it.

## Done when

- a Go test sees `removed` sent before the connection ends when an entry is unregistered, and not when an entry changes
- a flaiover test sees the project leave `connectedProjects()` after `removed` and a close, and come back on a new connection
- `npm test` in flaiover and `scripts/flai-test.sh` pass

## Notes
