---
id: T-0595
type: task
nature: feature
title: "A right click on a board lane opens its menu: create, move stories forward or back, change the WIP limit"
status: done
parent: S-0167
owner: alex
created: 2026-09-29T23:33:31Z
updated: 2026-09-29T23:49:18Z
transitions:
  - to: ready
    at: 2026-09-29T23:33:34Z
    by: agent-S-0167
  - to: in-progress
    at: 2026-09-29T23:42:45Z
    by: agent-S-0167
  - to: done
    at: 2026-09-29T23:49:18Z
    by: agent-S-0167
stream: S-0167
tags: []
touches: [flaiover/src/routes/board, flaiover/src/lib/components, flaiover/src/lib, design/system/flaiover-dashboard.md, docs/users]
usage:
  source: log
  seconds: 393
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 85
      output: 27564
      cache_read: 7351852
      cache_write: 85563
      cost: 2.7065
---
# T-0595 A right click on a board lane opens its menu: create, move stories forward or back, change the WIP limit

## Work

- A context menu on each lane (right click, and the keyboard's context-menu key or Shift+F10), offering only what the lane allows:
  - create item: every lane, opening `/new?status=<lane>`;
  - move stories forward: backlog only, to ready;
  - move stories back: ready to backlog, in-progress to ready, review to in-progress (with the reason flai requires), cancelled to backlog;
  - change WIP limit: ready, in-progress, review.
- Moving opens a dialog listing the lane's stories to tick; a right click on a card ticks that card. Each move goes through the existing move route and the result names every story moved and every refusal.
- The WIP limit dialog posts to a new `/api/board/limit` route, which runs `board.limit` on the host.
- The menu is not offered on a read-only board.
- Component tests for which entries each lane shows, the move dialog, and the limit dialog; update `design/system/flaiover-dashboard.md` and the dashboard user docs.

## Done when

- Every acceptance criterion of S-0167 is met on the board.
- `npm run check`, lint, and the unit tests pass in `flaiover/`.

## Notes
