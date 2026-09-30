---
id: T-0609
type: task
nature: feature
title: The user guide and the dashboard design describe the two-tier menu
status: done
parent: S-0172
owner: alex
created: 2026-09-30T01:18:30Z
updated: 2026-09-30T01:25:02Z
transitions:
  - to: ready
    at: 2026-09-30T01:18:44Z
    by: agent-S-0172
  - to: in-progress
    at: 2026-09-30T01:24:07Z
    by: agent-S-0172
  - to: done
    at: 2026-09-30T01:25:02Z
    by: agent-S-0172
stream: S-0172
tags: []
usage:
  source: log
  seconds: 55
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 24
      output: 6973
      cache_read: 1201241
      cache_write: 26474
      cost: 0.5916
---

# T-0609 The user guide and the dashboard design describe the two-tier menu

## Work

- `docs/users/flaiover.md` § The header: the groups, their pages, and how the menu opens.
- `design/system/flaiover-dashboard.md`: the menu's components and the Settings and Host rows' place in it.

## Done when

- Both documents describe the menu as built, with `updated` bumped, and `flai check --strict` passes.

## Notes
