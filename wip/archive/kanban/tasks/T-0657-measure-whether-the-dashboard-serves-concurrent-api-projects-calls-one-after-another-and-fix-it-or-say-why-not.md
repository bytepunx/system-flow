---
id: T-0657
type: task
nature: feature
title: Measure whether the dashboard serves concurrent /api/projects calls one after another, and fix it or say why not
status: done
parent: S-0186
owner: alex
created: 2026-10-01T09:18:12Z
updated: 2026-10-01T09:23:06Z
transitions:
  - to: ready
    at: 2026-10-01T09:20:07Z
    by: agent-S-0186
  - to: in-progress
    at: 2026-10-01T09:20:07Z
    by: agent-S-0186
  - to: done
    at: 2026-10-01T09:23:06Z
    by: agent-S-0186
stream: S-0186
tags: []
touches: [flaiover/src/routes/api/projects, flaiover/src/lib/server]
usage:
  source: log
  seconds: 179
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 57
      output: 12674
      cache_read: 2861044
      cache_write: 55485
      cost: 1.2002
---
# T-0657 Measure whether the dashboard serves concurrent /api/projects calls one after another, and fix it or say why not

## Work

- Read how `hub.ask` and flai's side of the dashboard connection handle several requests at once, and measure concurrent `/api/projects` calls against a project that answers slowly.
- If they are served one after another, fix the cause with a test that fails without the fix; if not, record the measurement and why.

## Done when

- The measurement and its result are in the narrative, and either a fix with its test is committed or the reason there is none is ready for the design.

## Notes
