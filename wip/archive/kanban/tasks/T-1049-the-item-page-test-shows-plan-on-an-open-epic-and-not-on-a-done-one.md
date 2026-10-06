---
id: T-1049
type: task
nature: improvement
title: The item page test shows Plan on an open epic and not on a done one
status: done
parent: S-0300
owner: alex
created: 2026-10-06T21:47:15Z
updated: 2026-10-06T23:07:08Z
transitions:
  - to: ready
    at: 2026-10-06T23:03:11Z
    by: agent-S-0300
  - to: in-progress
    at: 2026-10-06T23:03:11Z
    by: agent-S-0300
  - to: done
    at: 2026-10-06T23:07:08Z
    by: agent-S-0300
stream: S-0300
tags: [dashboard, planner]
touches: ["flaiover/src/routes/items/[id]/item.svelte.test.ts"]
usage:
  source: log
  seconds: 237
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 18
      output: 5853
      cache_read: 794705
      cache_write: 30155
      cost: 0.4719
---
# T-1049 The item page test shows Plan on an open epic and not on a done one

## Work

Criteria 1 and 2 are met already, and this task proves it.

- Criterion 1: S-0208's `PlanAction.svelte` shows on an epic's page through `canEdit` in `flaiover/src/routes/items/[id]/+page.svelte`.
- Criterion 2: S-0263's `cardMenu` offers Plan on an open epic, tested in `flaiover/src/lib/cardmenu.test.ts`.

No item page test covers an epic, so add one to `item.svelte.test.ts`. With a writable dashboard and `GET /api/items/E-nnnn/plan` answering `plan_enabled: true`, an open epic's page shows **Plan**. A done or archived epic's page does not. When you work this task, also open an epic on the running dashboard and right-click its card on the board with the epics checkbox ticked. Note what you saw in the narrative.

If the test finds Plan missing on an epic, fix `+page.svelte` too and add it to this task's touches.

This task waits for nothing and runs alongside the design and prompt tasks: it shares no path with them.

## Done when

- [ ] `item.svelte.test.ts` checks that Plan shows on an open epic's page and not on a done or archived one, and `npm test` in `flaiover` passes.
- [ ] The narrative records that the epic page and the epic card's menu each offered Plan.

## Notes

Drafted by the planner. TH-0201 tells the operator that criteria 1 and 2 are already met.
