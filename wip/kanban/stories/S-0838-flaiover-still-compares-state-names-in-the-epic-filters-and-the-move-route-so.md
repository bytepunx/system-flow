---
id: S-0838
type: story
nature: remediation
title: flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there
status: ready
owner: alex
created: 2026-11-17T08:32:49Z
updated: 2026-11-17T12:35:58Z
transitions:
  - to: ready
    at: 2026-11-17T12:35:53Z
    by: orchestrator
tags: [flaiover, workflow]
topics: [dashboard]
touches: ["flaiover/src/routes/api/items/[id]/move/+server.ts", "flaiover/src/routes/api/items/[id]/move/kinds.test.ts", flaiover/src/lib/lanes.fixture.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts, flaiover/src/lib/components/design/ProposeChange.svelte, flaiover/src/lib/components/design/ProposeChange.svelte.test.ts, design/components/flaiover/src/lib/components/design/ProposeChange.examples.json, flaiover/src/lib/components/DeploySettings.svelte, flaiover/src/lib/components/DeploySettings.svelte.test.ts, flaiover/src/lib/components/NewItemForm.svelte, flaiover/src/lib/components/NewItemForm.svelte.test.ts, design/system/flaiover-dashboard.md, design/issues/I-0399-flaiover-still-compares-state-names-in-the-epic-filters-and-the-move-route-so.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 14
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 17
          cache_read: 124457
          cache_write: 4105
          cost: 0.0402
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-17T08:32:49Z
  value: 25
  by: planner-S-0838
  at: 2026-11-17T12:32:32Z
forecast:
  duration: 34m
  delivery: 2026-11-17T13:19:00Z
  basis: "Its own forecast of 34m; 2nd in the pull order with an in-progress limit of 3, behind S-0806, S-0834 and S-0811."
  by: flai
  at: 2026-11-17T12:35:58Z
finalized:
  by: orchestrator
  at: 2026-11-17T12:35:32Z
---
# S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there

## Goal

This story remediates [I-0399](../../../design/issues/I-0399-flaiover-still-compares-state-names-in-the-epic-filters-and-the-move-route-so.md), "flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0399 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0399 is closed with `flai issue close I-0399 --reason` saying what fixed it

## Tasks
- T-3568 The move route judges a cancellation's preview and an acceptance's timeout by the target lane's kind
- T-3571 The design page offers only open epics by their lane's kind, and a requested change moves to the workflow's pullable lane
- T-3575 The new-item form and the Deployments page's story form offer only open epics by their lane's kind
- T-3579 flaiover-dashboard.md names the epic filters and the move route among what reads a state's kind, and I-0399 is closed

## Notes

Cost of delay inputs set by flai from I-0399. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-17T08:25:17Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0399's one instance: each site reads the board's flow and asks `$lib/workflow` what a state is, as S-0758's views do. The three epic filters use `isOpen(flow, e.status)` over the flow from `/api/board`. The move route reads the flow from the repo's board view on the server and uses `isCancelled` for the dry-run preview and `isAccepted` for the acceptance's timeout. `ProposeChange` moves a requested change through `movesTo` to the first start lane that `isPullable` names, in place of `ready`. A fixture flow with renamed accepted, cancelled, and pullable lanes in `lanes.fixture.ts` reproduces the cause in each test.

Touches, where each came from:

| Touch | Source |
|-------|--------|
| `flaiover/src/routes/api/items/[id]/move/+server.ts`, `src/routes/design/+page.svelte`, `src/lib/components/DeploySettings.svelte`, `src/lib/components/NewItemForm.svelte`, `src/lib/components/design/ProposeChange.svelte` | design: the files I-0399's instance names |
| `flaiover/src/routes/api/items/[id]/move/kinds.test.ts` | layout: a new test beside the route's `designer.test.ts` |
| `flaiover/src/lib/lanes.fixture.ts` | layout: where S-0758's baseline and added-state flows live |
| `design.svelte.test.ts`, `ProposeChange.svelte.test.ts`, `DeploySettings.svelte.test.ts`, `NewItemForm.svelte.test.ts` | layout: each changed view's existing test |
| `design/components/flaiover/src/lib/components/design/ProposeChange.examples.json` | layout: `ProposeChange` gains a required `flow` prop, and its examples must keep in step |
| `design/system/flaiover-dashboard.md` | design: § A state's kind (S-0758) says no view compares a state's name, and must list these |
| `design/issues/I-0399-…md`, `design/issues/summary.md` | design: the second criterion's `flai issue close` writes them |

`flai touches suggest` co-change lists `docs/users/flaiover.md` (14%) and `design/system/flai-cli.md` (10%) first. Neither is a touch: the user guide describes the epic lists as open epics, which the fix makes true again, and flai is not changed. No folder touch is kept. The story names no `after`: it reads nothing that an open story writes.

Forecast 34m, delivery 2026-11-17T20:20:00Z, kept as `flai forecast` gives it. Three of its 15 touches are files flai writes and an examples file, which inflate the size; four test files that need a `/api/board` answer under-count it. The two offset, so no adjustment.

Cost of delay value 25 USD a week, kept as `flai cod` works it out from the inputs flai set from I-0399: 10m lost per 168h cycle at 150 USD an hour. One instance gives no reason to adjust.

Tasks form three layers: T-3568 first, since it adds the fixture the others' tests read; T-3571 and T-3575 together, sharing no file; T-3579 last, since it documents and closes what those changed. Topic `dashboard` is added, since the tasks change flaiover and its design document.
