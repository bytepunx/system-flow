---
id: T-0212
type: task
nature: feature
title: "Dashboard server: backlog and ready columns in pull order, and an endpoint that calls flai order"
status: done
parent: S-0057
owner: alex
created: 2026-09-19T06:41:25Z
updated: 2026-09-19T06:50:27Z
transitions:
  - to: ready
    at: 2026-09-19T06:48:26Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T06:48:26Z
    by: system-flow
  - to: done
    at: 2026-09-19T06:50:27Z
    by: system-flow
stream: S-0057
tags: []
---

# T-0212 Dashboard server: backlog and ready columns in pull order, and an endpoint that calls flai order

## Work
`flaiover/src/lib/server/board.ts` returns the backlog and ready columns with their stories in the pull sequence, the same rule as flai's, leaving epics and tasks where they are; other columns unchanged. `POST /api/items/:id/order` with `{ before | after | top | bottom }` runs `flai order` (ADR-0016: the dashboard never edits `board.md`), returns the new order, passes a rule refusal through as 409 with flai's message, and answers 503 when flai is not available. Tests for the sorting and the endpoint's argument building.

## Done when
- The board API returns backlog and ready in pull order, with tests
- The endpoint calls flai and reports refusals, with tests

## Notes
