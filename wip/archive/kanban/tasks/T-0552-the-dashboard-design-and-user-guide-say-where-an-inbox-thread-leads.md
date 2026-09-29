---
id: T-0552
type: task
nature: feature
title: The dashboard design and user guide say where an inbox thread leads
status: done
parent: S-0155
owner: alex
created: 2026-09-29T07:02:04Z
updated: 2026-09-29T07:05:45Z
transitions:
  - to: ready
    at: 2026-09-29T07:02:22Z
    by: agent-S-0155
  - to: in-progress
    at: 2026-09-29T07:05:23Z
    by: agent-S-0155
  - to: done
    at: 2026-09-29T07:05:45Z
    by: agent-S-0155
stream: S-0155
tags: []
usage:
  source: log
  seconds: 22
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 14
      output: 2819
      cache_read: 680722
      cache_write: 14002
      cost: 0.3046
---

# T-0552 The dashboard design and user guide say where an inbox thread leads

## Work

`design/system/flaiover-dashboard.md` (the `/inbox` and `/items/:id` rows) and `docs/users/flaiover.md` say that a thread entry opens its story's page on that thread.

## Done when

Both documents say it, with `updated` bumped; `flai check --strict` passes.

## Notes
