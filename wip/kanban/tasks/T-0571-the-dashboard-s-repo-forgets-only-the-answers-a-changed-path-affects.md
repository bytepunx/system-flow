---
id: T-0571
type: task
nature: feature
title: The dashboard's Repo forgets only the answers a changed path affects
status: backlog
parent: S-0161
owner: alex
created: 2026-09-29T19:40:45Z
updated: 2026-09-29T19:40:45Z
transitions: []
stream: S-0161
tags: []
touches: [flaiover/src/lib, design/system/flaiover-dashboard.md]
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
