---
id: T-0094
type: task
nature: feature
title: GET /api/stats delegating to flai stats --json with a change-invalidated cache
status: done
parent: S-0014
owner: alex
created: 2026-09-17T04:50:13Z
updated: 2026-09-17T05:02:57Z
transitions:
  - to: ready
    at: 2026-09-17T05:02:57Z
    by: agent
  - to: in-progress
    at: 2026-09-17T05:02:57Z
    by: agent
  - to: done
    at: 2026-09-17T05:02:57Z
    by: agent
stream: S-0014
tags: [dashboard, charts]
---

# T-0094 GET /api/stats delegating to flai stats --json with a change-invalidated cache

## Work
src/lib/server/stats.ts validates since/type/by, runs flai stats --json through the wrapper, caches per argument set, and clears the cache on Repo change events; GET /api/stats passes the query through and maps bad arguments to 400.

## Done when
Unit tests for validation and a wrapper-backed test on a fixture copy.

## Notes
