---
id: T-0835
type: task
nature: feature
title: The design and the user guide say the card menu offers Plan
status: done
parent: S-0263
owner: alex
created: 2026-10-04T23:39:24Z
updated: 2026-10-04T23:44:05Z
transitions:
  - to: ready
    at: 2026-10-04T23:39:45Z
    by: agent-S-0263
  - to: in-progress
    at: 2026-10-04T23:43:00Z
    by: agent-S-0263
  - to: done
    at: 2026-10-04T23:44:05Z
    by: agent-S-0263
stream: S-0263
tags: []
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-0834]
usage:
  source: log
  seconds: 65
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 21
      output: 5634
      cache_read: 838530
      cache_write: 32772
      cost: 0.4954
---

# T-0835 The design and the user guide say the card menu offers Plan

## Work

Add Plan to the card menu in `design/system/flaiover-dashboard.md` (the `/board` row's card menu entries, and `plan_enabled` on `/api/host-agent`) and to the card menu table in `docs/users/flaiover.md`. Waits for T-0834, so that the documents describe what it built.

## Done when

- Both documents name Plan, where it is offered, and what it does, and the markdown lint passes on them.

## Notes
