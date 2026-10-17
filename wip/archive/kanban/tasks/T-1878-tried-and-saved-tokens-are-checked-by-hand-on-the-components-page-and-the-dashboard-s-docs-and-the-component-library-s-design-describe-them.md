---
id: T-1878
type: task
nature: feature
title: Tried and saved tokens are checked by hand on the Components page, and the dashboard's docs and the component library's design describe them
status: done
parent: S-0441
owner: alex
created: 2026-10-16T21:36:40Z
updated: 2026-10-17T03:21:53Z
transitions:
  - to: ready
    at: 2026-10-17T03:07:36Z
    by: agent-S-0441
  - to: in-progress
    at: 2026-10-17T03:07:37Z
    by: agent-S-0441
  - to: done
    at: 2026-10-17T03:21:53Z
    by: agent-S-0441
stream: S-0441
tags: [dashboard]
touches: [docs/users/flaiover.md, design/system/component-library.md, design/system/flaiover-dashboard.md]
after: [T-1876, T-1877]
usage:
  source: log
  seconds: 856
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 145
      output: 55695
      cache_read: 6509406
      cache_write: 292984
      cost: 4.0671
---
# T-1878 Tried and saved tokens are checked by hand on the Components page, and the dashboard's docs and the component library's design describe them

## Work

- On a fixture project with a token file and a stylesheet, or on this repository once its design system exists: try a colour on the Design page and see it on the Components page in another tab, then save it and see it stay after the tried values are dropped. Record both in the story's notes.
- Describe tried and saved tokens in `docs/users/flaiover.md` and `design/system/component-library.md`, and the Design page's kept draft and the `tokens` message in `design/system/flaiover-dashboard.md` § Editing design tokens (S-0383) and the Components page's entry.

Waits for T-1876 and T-1877, the two halves the check runs through.

## Done when

The story's notes record both checks, and the three files describe the behaviour, passing `flai test` on them.

## Notes

None.
