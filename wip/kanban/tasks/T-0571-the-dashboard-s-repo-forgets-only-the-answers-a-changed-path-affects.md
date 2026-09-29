---
id: T-0571
type: task
nature: feature
title: The dashboard's Repo forgets only the answers a changed path affects
status: done
parent: S-0161
owner: alex
created: 2026-09-29T19:40:45Z
updated: 2026-09-29T19:45:14Z
transitions:
  - to: ready
    at: 2026-09-29T19:41:10Z
    by: agent-S-0161
  - to: in-progress
    at: 2026-09-29T19:41:11Z
    by: agent-S-0161
  - to: done
    at: 2026-09-29T19:45:14Z
    by: agent-S-0161
stream: S-0161
tags: []
touches: [flaiover/src/lib, design/system/flaiover-dashboard.md]
usage:
  source: log
  seconds: 243
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 44
      output: 326
      cache_read: 2971080
      cache_write: 28843
      cost: 1.1941
---

# T-0571 The dashboard's Repo forgets only the answers a changed path affects

## Work

- `src/lib/changes.ts`: a changed path's kind (`project`, `item`, `narrative`, `thread`, `adr`, `document`, `other`) from the manifest's layout; importable by pages and the server.
- `Repo.changed` forgets only the answers whose inputs include that kind, by a table of what each remembered answer is read from in flai; an answer not in the table, and any path of kind `other` or a change before the layout is known, forgets everything.
- `/api/events` sends the kind with the path.
- `design/system/flaiover-dashboard.md` says which paths forget which answers.

## Done when

- Behaviour tests in `flaiover/src/lib/server` show a narrative, thread, or document change keeps `board.get` and the other unaffected answers, and a work item or `board.md` change forgets them.
- `npm run check`, lint, and the unit tests pass in `flaiover/`.

## Notes
