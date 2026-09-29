---
id: T-0578
type: task
nature: improvement
title: The dashboard asks for items without bodies and the archive only where it is shown
status: done
parent: S-0162
owner: alex
created: 2026-09-29T20:15:55Z
updated: 2026-09-29T20:23:56Z
transitions:
  - to: ready
    at: 2026-09-29T20:16:04Z
    by: agent-S-0162
  - to: in-progress
    at: 2026-09-29T20:19:16Z
    by: agent-S-0162
  - to: done
    at: 2026-09-29T20:23:56Z
    by: agent-S-0162
stream: S-0162
tags: []
touches: [flaiover/src]
usage:
  source: log
  seconds: 280
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 58
      output: 15449
      cache_read: 4349355
      cache_write: 57919
      cost: 1.6424
---
# T-0578 The dashboard asks for items without bodies and the archive only where it is shown

## Work

- `Repo.items()` asks `items.list` without bodies, by type and status, with the archive only when asked, each query kept under its own key.
- `/api/items` passes its query to flai; the overview, the document explorer, and the editor ask for the active items; the overview asks `items.count` for the archived number.
- `/_ready` asks for the manifest alone.
- `docs.tree`'s type loses `frontMatter`.

## Done when

- The overview's first answer from flai is under 200 KB on this repository, measured.
- `npm run check`, `npm run lint`, and `npm test` in `flaiover` pass.

## Notes
