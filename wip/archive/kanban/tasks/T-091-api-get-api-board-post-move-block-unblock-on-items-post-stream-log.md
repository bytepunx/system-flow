---
id: T-091
type: task
nature: feature
title: "API: GET /api/board; POST move, block, unblock on items; POST stream log"
status: done
parent: S-013
owner: alex
created: 2026-09-17T04:32:27Z
updated: 2026-09-17T04:38:32Z
transitions:
  - to: ready
    at: 2026-09-17T04:38:32Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:38:32Z
    by: agent
  - to: done
    at: 2026-09-17T04:38:32Z
    by: agent
stream: S-013
tags: [dashboard]
---

# T-091 API: GET /api/board; POST move, block, unblock on items; POST stream log

## Work
GET /api/board returns wip_limits, order, and columns of cards (age in column, blocked, nature, parent) from the reader; POST /api/items/[id]/move {to, reason?, by?}, /block {reason}, /unblock, and POST /api/streams/[id]/log {entry} delegate to flai and return its JSON; errors carry the rule.

## Done when
Handler tests: a valid move changes status, an invalid one returns 400 with the rule named.

## Notes
