---
id: T-0230
type: task
nature: feature
title: "Dashboard endpoints: create an item through flai, and the template body"
status: done
parent: S-0059
owner: alex
created: 2026-09-19T09:21:36Z
updated: 2026-09-19T09:27:40Z
transitions:
  - to: ready
    at: 2026-09-19T09:25:01Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T09:25:02Z
    by: system-flow
  - to: done
    at: 2026-09-19T09:27:40Z
    by: system-flow
stream: S-0059
tags: []
---

# T-0230 Dashboard endpoints: create an item through flai, and the template body

## Work
`POST /api/items` `{ type, title, nature, parent?, tags?, touches?, body }` runs `flai <type> new` with the body on standard input, `--autocommit`, the designer as owner, and the dashboard's trailer; only `epic` and `story` are accepted; a refusal by the check is 409 with the findings; a rule or usage error is 400. `GET /api/items/template?type=` returns the template body from `flai <type> new --print-body`. Tests for the argument building and, through the real binary, a creation and a refusal.

## Done when
- The endpoints are tested, and nothing but item IDs, enumerated values, and flag-safe strings reach flai's arguments

## Notes
