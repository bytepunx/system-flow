---
id: T-092
type: task
nature: feature
title: "Board page with columns, cards, WIP counts, drag to transition; item page with history, children, actions, narrative link"
status: done
parent: S-013
owner: alex
created: 2026-09-17T04:32:27Z
updated: 2026-09-17T04:38:33Z
transitions:
  - to: ready
    at: 2026-09-17T04:38:32Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:38:32Z
    by: agent
  - to: done
    at: 2026-09-17T04:38:33Z
    by: agent
stream: S-013
tags: [dashboard]
---

# T-092 Board page with columns, cards, WIP counts, drag to transition; item page with history, children, actions, narrative link

## Work
/board: one column per state with WIP count against the limit, story cards (toggle to include tasks and epics) showing ID, title, nature, age, blocked flag; HTML5 drag and drop posts a move and shows the rule on refusal; refreshes on SSE change events. /items/[id]: front matter, rendered body, children table, transitions timeline, blocked intervals, narrative link, and action buttons for allowed moves, block, unblock, and a log entry.

## Done when
Board matches flai board for this repository; a drag on the running app performs a real move.

## Notes
