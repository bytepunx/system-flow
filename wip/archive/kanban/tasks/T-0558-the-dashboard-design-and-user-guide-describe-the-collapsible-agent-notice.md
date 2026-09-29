---
id: T-0558
type: task
nature: feature
title: The dashboard design and user guide describe the collapsible agent notice
status: done
parent: S-0150
owner: alex
created: 2026-09-29T19:05:19Z
updated: 2026-09-29T19:09:12Z
transitions:
  - to: ready
    at: 2026-09-29T19:05:34Z
    by: agent-S-0150
  - to: in-progress
    at: 2026-09-29T19:08:40Z
    by: agent-S-0150
  - to: done
    at: 2026-09-29T19:09:12Z
    by: agent-S-0150
stream: S-0150
tags: []
usage:
  source: log
  seconds: 32
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 15
      output: 3996
      cache_read: 613616
      cache_write: 20338
      cost: 0.3654
---

# T-0558 The dashboard design and user guide describe the collapsible agent notice

## Work

- `design/system/flaiover-dashboard.md` § Agents started by flai: the notice collapses to its title, is remembered per browser, and opens again when its status changes.
- `docs/users/flaiover.md`: where the board's agent notice is described, say how to collapse it.

## Done when

- Both documents say what the notice does, with `updated` bumped, and `flai check --strict` passes.

## Notes
