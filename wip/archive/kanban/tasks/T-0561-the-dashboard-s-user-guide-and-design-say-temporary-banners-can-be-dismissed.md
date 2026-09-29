---
id: T-0561
type: task
nature: feature
title: The dashboard's user guide and design say temporary banners can be dismissed
status: done
parent: S-0151
owner: alex
created: 2026-09-29T19:11:39Z
updated: 2026-09-29T19:16:05Z
transitions:
  - to: ready
    at: 2026-09-29T19:11:56Z
    by: agent-S-0151
  - to: in-progress
    at: 2026-09-29T19:15:25Z
    by: agent-S-0151
  - to: done
    at: 2026-09-29T19:16:05Z
    by: agent-S-0151
stream: S-0151
tags: []
usage:
  source: log
  seconds: 40
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 18
      output: 4178
      cache_read: 864279
      cache_write: 20010
      cost: 0.4166
---

# T-0561 The dashboard's user guide and design say temporary banners can be dismissed

## Work

Say in `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md` that a banner an action leaves has an X that dismisses it, and which banners are not dismissible because they show state.

## Done when

Both documents say it, with `updated` bumped, and `flai check --strict` passes.

## Notes
