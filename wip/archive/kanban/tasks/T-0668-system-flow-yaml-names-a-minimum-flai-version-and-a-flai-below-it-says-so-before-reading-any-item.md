---
id: T-0668
type: task
nature: feature
title: system-flow.yaml names a minimum flai version, and a flai below it says so before reading any item
status: done
parent: S-0181
owner: arobson
created: 2026-10-01T10:43:14Z
updated: 2026-10-01T10:50:53Z
transitions:
  - to: ready
    at: 2026-10-01T10:46:36Z
    by: agent-S-0181
  - to: in-progress
    at: 2026-10-01T10:46:36Z
    by: agent-S-0181
  - to: done
    at: 2026-10-01T10:50:53Z
    by: agent-S-0181
stream: S-0181
tags: []
usage:
  source: log
  seconds: 257
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 67
      output: 18869
      cache_read: 4624439
      cache_write: 78902
      cost: 1.8221
---

# T-0668 system-flow.yaml names a minimum flai version, and a flai below it says so before reading any item

## Work

Add a minimum flai version to `system-flow.yaml`. Every command that opens the project checks it before reading any item, and a flai below it stops with the version needed and the command that upgrades the host.

## Done when

- A manifest minimum above the running version stops listing with that message before any item is read, tested
- A dev build and a manifest with no minimum are not stopped

## Notes
