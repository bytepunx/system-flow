---
id: T-0669
type: task
nature: feature
title: flai serve, the MCP inbox, and the dashboard's host badge say when the running flai is older than the newest flai tag reachable from HEAD
status: done
parent: S-0181
owner: arobson
created: 2026-10-01T10:43:14Z
updated: 2026-10-01T10:53:42Z
transitions:
  - to: ready
    at: 2026-10-01T10:50:53Z
    by: agent-S-0181
  - to: in-progress
    at: 2026-10-01T10:50:53Z
    by: agent-S-0181
  - to: done
    at: 2026-10-01T10:53:42Z
    by: agent-S-0181
stream: S-0181
tags: []
usage:
  source: log
  seconds: 169
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 82
      output: 23328
      cache_read: 5717090
      cache_write: 97545
      cost: 2.2526
---

# T-0669 flai serve, the MCP inbox, and the dashboard's host badge say when the running flai is older than the newest flai tag reachable from HEAD

## Work

Find the newest `flai/v*` tag reachable from the project's HEAD and compare it with the running flai's version. `flai serve` logs a warning at start, the MCP `inbox` carries it, and the dashboard's host badge shows it, each with the upgrade command.

## Done when

- The comparison is tested against a repository with tags above and below the running version
- `flai serve`, `inbox`, and the host badge carry the warning, tested where each is tested

## Notes
