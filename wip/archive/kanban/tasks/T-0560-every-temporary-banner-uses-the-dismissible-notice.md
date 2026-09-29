---
id: T-0560
type: task
nature: feature
title: Every temporary banner uses the dismissible notice
status: done
parent: S-0151
owner: alex
created: 2026-09-29T19:11:39Z
updated: 2026-09-29T19:15:25Z
transitions:
  - to: ready
    at: 2026-09-29T19:11:56Z
    by: agent-S-0151
  - to: in-progress
    at: 2026-09-29T19:12:52Z
    by: agent-S-0151
  - to: done
    at: 2026-09-29T19:15:25Z
    by: agent-S-0151
stream: S-0151
tags: []
usage:
  source: log
  seconds: 153
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 42
      output: 9729
      cache_read: 2012372
      cache_write: 46590
      cost: 0.9699
---

# T-0560 Every temporary banner uses the dismissible notice

## Work

Put the temporary banners on `DismissibleNotice`, each dismissed by clearing the state it shows: the board's action notice, the item page's notice, the ADRs page's notice, the document editor's saved notice, and the Publish result above the done column. State notices (flai not connected, unpushed, host agent, pending publish) are not temporary and stay as they are.

## Done when

Clicking the X on each of those banners removes it; component tests cover the item page, the document editor, and the Publish result; `npm run check`, `npm run lint`, and the unit tests pass.

## Notes
