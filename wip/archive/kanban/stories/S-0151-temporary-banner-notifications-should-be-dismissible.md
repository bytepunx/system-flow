---
id: S-0151
type: story
nature: improvement
title: Temporary banner notifications should be dismissible
status: done
parent: E-0013
owner: alex
created: 2026-09-29T05:45:39Z
updated: 2026-09-29T19:18:03Z
transitions:
  - to: ready
    at: 2026-09-29T06:44:33Z
    by: alex
  - to: in-progress
    at: 2026-09-29T19:10:36Z
    by: agent-S-0151
  - to: review
    at: 2026-09-29T19:16:42Z
    by: agent-S-0151
  - to: done
    at: 2026-09-29T19:18:03Z
    by: alex
tags: [dashboard]
topics: [client-side]
touches: [flaiover/src, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 398
  models:
    - model: claude-opus-5-5
      input: 112
      output: 25704
      cache_read: 5316658
      cache_write: 123091
      cost: 2.5626
---
# S-0151 Temporary banner notifications should be dismissible

## Goal

A right aligned X should appear in all temporary notifications so that operators can dismiss them

## Acceptance criteria
- [x] Clicking the X dismisses the banner

## Tasks
- T-0559 A dismissible notice component with a right-aligned X
- T-0560 Every temporary banner uses the dismissible notice
- T-0561 The dashboard's user guide and design say temporary banners can be dismissed

## Notes

The temporary banners are the board's action notice, the item page's notice, the `/adrs` notice, the document editor's saved notice, and the Publish result above the done column. Each uses `flaiover/src/lib/components/DismissibleNotice.svelte`. Banners that show state stay without an X. Component tests cover the X on the component, the item page, the document editor, and the Publish result. The board and `/adrs` have no page tests; they use the same component with the same `ondismiss` that clears their notice.
