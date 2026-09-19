---
id: T-0227
type: task
nature: feature
title: The standing notice on the board and the item page
status: done
parent: S-0063
owner: alex
created: 2026-09-19T08:34:59Z
updated: 2026-09-19T08:42:01Z
transitions:
  - to: ready
    at: 2026-09-19T08:37:34Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T08:37:34Z
    by: system-flow
  - to: done
    at: 2026-09-19T08:42:01Z
    by: system-flow
stream: S-0063
tags: []
---

# T-0227 The standing notice on the board and the item page

## Work
The board API carries `unpushed` (the dashboard reads it through `flai board --json`, or the same git questions, whichever the reader's design allows without shelling out on every request more than it does now). The board page and the item page of an accepted story show a standing notice while it is true: the counts, the tags, the command, and that a push made from another clone is not seen until someone fetches here. It clears by itself. No retry button. Component tests; tried in a browser against a scratch project with a scratch remote.

## Done when
- The notice shows on every load while true and disappears after a push from the clone
- `make flaiover-test` and `make flaiover-build` pass

## Notes
