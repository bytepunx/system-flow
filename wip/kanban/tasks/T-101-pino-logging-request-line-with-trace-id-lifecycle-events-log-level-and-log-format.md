---
id: T-101
type: task
nature: feature
title: "pino logging: request line with trace_id, lifecycle events, LOG_LEVEL and LOG_FORMAT"
status: done
parent: S-032
owner: alex
created: 2026-09-17T05:52:19Z
updated: 2026-09-17T06:13:51Z
transitions:
  - to: ready
    at: 2026-09-17T06:13:51Z
    by: alex
  - to: in-progress
    at: 2026-09-17T06:13:51Z
    by: alex
  - to: done
    at: 2026-09-17T06:13:51Z
    by: alex
stream: S-032
tags: [dashboard, logging]
---

# T-101 pino logging: request line with trace_id, lifecycle events, LOG_LEVEL and LOG_FORMAT

## Work
src/lib/server/log.ts builds a pino logger (JSON by default, pino-pretty text when LOG_FORMAT=text outside production, LOG_LEVEL, ts and level keys per the convention); hooks.server.ts logs one request line per request with trace_id, method, route, path, status, duration_ms, a lifecycle event at start, and errors once in handleError; health and metrics paths are quiet.

## Done when
Log lines are JSON events with the convention's keys; verified on the served build.

## Notes
