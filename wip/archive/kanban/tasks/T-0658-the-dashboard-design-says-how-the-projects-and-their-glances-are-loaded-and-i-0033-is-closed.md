---
id: T-0658
type: task
nature: feature
title: The dashboard design says how the projects and their glances are loaded, and I-0033 is closed
status: done
parent: S-0186
owner: alex
created: 2026-10-01T09:18:14Z
updated: 2026-10-01T09:24:19Z
transitions:
  - to: ready
    at: 2026-10-01T09:23:06Z
    by: agent-S-0186
  - to: in-progress
    at: 2026-10-01T09:23:06Z
    by: agent-S-0186
  - to: done
    at: 2026-10-01T09:24:19Z
    by: agent-S-0186
stream: S-0186
tags: []
touches: [design/system/flaiover-dashboard.md, design/issues/I-0033-the-root-page-asks-api-projects-twice-on-one-load-and-the-dashboard-serves-the-two-calls-serially.md, design/issues/summary.md]
usage:
  source: log
  seconds: 73
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 38
      output: 8478
      cache_read: 1913823
      cache_write: 37115
      cost: 0.8028
---
# T-0658 The dashboard design says how the projects and their glances are loaded, and I-0033 is closed

## Work

- `design/system/flaiover-dashboard.md`: how `projectState` loads the projects and their glances once per load, who asks again and when, and what the measurement of concurrent calls found.
- Close I-0033 with what fixed it and update `design/issues/summary.md`.

## Done when

- The design and the issue say what was done, and `flai check --strict` is clean.

## Notes
