---
id: T-0249
type: task
nature: improvement
title: "workitem: plan and apply a cancellation with everything open under it, all or nothing, and the ADR for stories in review"
status: done
parent: S-0070
owner: alex
created: 2026-09-20T06:16:43Z
updated: 2026-09-20T06:23:29Z
transitions:
  - to: ready
    at: 2026-09-20T06:17:24Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T06:17:24Z
    by: system-flow
  - to: done
    at: 2026-09-20T06:23:29Z
    by: system-flow
stream: S-0070
tags: []
---
# T-0249 workitem: plan and apply a cancellation with everything open under it, all or nothing, and the ADR for stories in review

## Work
In `flai/internal/workitem`: a plan that lists the open descendants of an epic or story (stories, then their tasks) with their current states, and an apply that validates every transition before writing any file, gives each child a `cancelled` transition with the parent's actor and timestamp and the note `<parent> cancelled: <why>`, takes cancelled stories out of the board order, and writes the index once. Decide what happens to a story in `review` and record it in an ADR refining ADR-0004. `Repo.Transition` to `cancelled` goes through it, so every caller cascades.

## Done when
- Tests cover an epic with stories in every state and tasks under them, a story with tasks, closed and archived children left alone, and a failure that leaves every file unchanged
- The ADR is recorded with `flai adr new`

## Notes
