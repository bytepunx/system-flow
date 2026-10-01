---
id: T-0649
type: task
nature: feature
title: flai host restarts a dashboard that is gone or not answering, with a back-off, and dashboard.no_restart turns it off
status: done
parent: S-0184
owner: arobson
created: 2026-10-01T08:55:54Z
updated: 2026-10-01T09:02:50Z
transitions:
  - to: ready
    at: 2026-10-01T08:56:22Z
    by: agent-S-0184
  - to: in-progress
    at: 2026-10-01T08:58:25Z
    by: agent-S-0184
  - to: done
    at: 2026-10-01T09:02:50Z
    by: agent-S-0184
stream: S-0184
tags: []
usage:
  source: log
  seconds: 265
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 41
      output: 13875
      cache_read: 3129875
      cache_write: 55981
      cost: 1.2869
---

# T-0649 flai host restarts a dashboard that is gone or not answering, with a back-off, and dashboard.no_restart turns it off

## Work

- `flai dashboard` (start, restart, upgrade) records the container it started: name, image reference, and publish; `flai dashboard stop` forgets it, so a stopped dashboard is not brought back.
- `flai/internal/host`: a watch loop that probes the recorded dashboard, restarts it when it is gone or has not answered for two probes in a row, waits a back-off that doubles between restarts and resets once it answers, and logs each restart with why.
- `dashboard.no_restart` in the config turns the watch off; the host reads it at every look.

## Done when

- Behaviour tests in `flai/internal/host` cover a restart of a gone and of a not-answering dashboard, the back-off, and the setting turning it off; `flai/cmd` tests cover the record.

## Notes
