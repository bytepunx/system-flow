---
id: S-0753
type: story
nature: feature
title: The dashboard's board draws the columns, their kinds, and the moves between them from what flai gives, and holds no list of states of its own
status: done
parent: E-0034
owner: alex
created: 2026-11-12T16:01:41Z
updated: 2026-11-17T04:46:04Z
transitions:
  - to: ready
    at: 2026-11-17T03:29:40Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-17T03:29:48Z
    by: agent-S-0753
  - to: review
    at: 2026-11-17T04:44:22Z
    by: agent-S-0753
  - to: done
    at: 2026-11-17T04:46:04Z
    by: orchestrator
tags: [dashboard, cli]
topics: [dashboard, cli]
touches: [flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flaiover/src/lib/server/board.ts, flaiover/src/lib/lanes.ts, flaiover/src/lib/lanes.test.ts, flaiover/src/lib/lanecollapse.svelte.ts, flaiover/src/lib/lanecollapse.svelte.test.ts, flaiover/src/lib/reorder.ts, flaiover/src/lib/reorder.test.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, flaiover/src/routes/board/lanemenu.svelte.test.ts, flaiover/src/routes/+page.svelte, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/components/CardMenu.svelte, flaiover/src/lib/components/LaneMenu.svelte, flaiover/src/lib/components/LaneMoveDialog.svelte, flaiover/src/lib/components/NewItemForm.svelte, flaiover/src/lib/components/NewItemForm.svelte.test.ts, flaiover/src/lib/lanes.fixture.ts, flaiover/src/lib/server/repo-channel.test.ts, flaiover/src/lib/server/writes.test.ts, flaiover/src/routes/board/cardmenu.svelte.test.ts, flaiover/src/routes/new/+page.svelte, flaiover/src/routes/overview.svelte.test.ts, design/components/flaiover/src/lib/components/CardMenu.examples.json, design/components/flaiover/src/lib/components/LaneMenu.examples.json, flai/internal/hostapi/writes_items.go, flai/internal/hostapi/writes_test.go]
after: [S-0739]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4489
  priced_by: harness
  turns:
    - day: 2026-11-17
      ceremony: 2
      hand_edits: 1
      work: 73
  models:
    - model: claude-opus-5-5
      input: 446
      output: 150477
      cache_read: 37644679
      cache_write: 914829
      cost: 16.093
  strategic:
    - kind: orchestrator
      seconds: 131
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 33
          output: 238
          cache_read: 1359306
          cache_write: 163524
          cost: 0.4657
        - model: claude-sonnet-5-5
          input: 8
          output: 55
          cache_read: 88828
          cache_write: 40686
          cost: 0.1037
cost_of_delay:
  value: 74.82
  by: system-flow
  at: 2026-11-14T17:44:13Z
forecast:
  duration: 42m
  delivery: 2026-11-17T04:20:00Z
  basis: "Its own forecast of 42m; 3rd in the pull order with an in-progress limit of 3, behind S-0750, S-0805, S-0241 and S-0752."
  by: flai
  at: 2026-11-17T03:29:09Z
finalized:
  by: orchestrator
  at: 2026-11-14T17:45:57Z
---
# S-0753 The dashboard's board draws the columns, their kinds, and the moves between them from what flai gives, and holds no list of states of its own

## Goal

The dashboard draws whatever board the workflow defines. `board.get` gives the columns in order with each one's kind, WIP limit, and order flag, and the moves between them with who may make each and whether it needs a reason, and flaiover draws the lanes, the card and lane menus, the drag targets, the reason prompt, and the collapsed lanes from that answer instead of its own lists in `board.ts`, `lanes.ts`, `lanecollapse.svelte.ts`, `reorder.ts`, and the two pages.

## Acceptance criteria

- [x] `board.get` answers the workflow's columns, each with its name, kind, WIP limit, and whether it is ordered, and its transitions, each with from, to, who may make it, and whether it needs a reason
- [x] `board.ts`, `lanes.ts`, `lanecollapse.svelte.ts`, `reorder.ts`, `routes/board/+page.svelte`, and `routes/+page.svelte` hold no list of the six states: the board, the lane and card menus, drag and drop, the reason prompt, reordering, collapsing, and the home page's counts read the columns and transitions from `board.get`
- [x] A move the operator may not make from the board is not offered, and one flai refuses still shows flai's words
- [x] Tests draw a board for a workflow with a state renamed and one added between `in-progress` and `review`, with its lane, its WIP limit, and its moves
- [x] With the baseline, the board looks and behaves as today, and every existing flaiover test passes, its fixtures updated only for the new fields
- [x] `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` say the board's columns come from the workflow

## Tasks

- T-3244 board.get answers the workflow's columns with their kinds, WIP limits, and order, and its transitions with who may make each and its reason
- T-3249 flaiover's board, lane and card menus, drag and drop, reason prompt, reordering, collapsing, and home page counts read the columns and moves from board.get
- T-3250 flaiover-dashboard.md and the dashboard's user guide say the board's columns and moves come from the workflow
- T-3509 board.limit takes the workflow's limited columns, and the lane and card menus' examples carry the flow

## Notes

- The lists today: `flaiover/src/lib/server/board.ts` (`STATES`), `lanes.ts` (`LANES`, `START_LANES`, `back`, `forwardOf`, `needsReason`), `lanecollapse.svelte.ts` (`laneStates`), `reorder.ts` (`ORDERED`), `routes/board/+page.svelte`, and `routes/+page.svelte`. `board.get` is in `flai/internal/hostapi/hostapi.go` and returns `workitem`'s board view.

### Planning

- Touches, all files:
  - `flai/internal/workitem/boardview.go`, `boardview_test.go`, `flai/internal/hostapi/hostapi.go`, `hostapi_test.go`: layout. `board.get` returns the board view.
  - The flaiover files named in the Notes, with their tests: layout, from the explorer's search for the six states' lists.
  - `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`: design (criterion 6).
  - Left out from `flai touches suggest`: `flai-cli.md`, `docs/users/flai.md`, and `docs/operators/index.md` (24% to 36%); `flaiover/src/lib/server/agent.ts` and `routes/items/[id]/+page.svelte` (7% to 8%), the second of which is S-0758's.
  - Criterion 2 names the board's files, since the planning search found about fifteen more flaiover files comparing a status with a state's name; S-0758 takes those.
- `after`: S-0739. `board.get` adds the transitions beside the columns, kinds, and WIP limits S-0739's board view gives.
- Tasks: T-3244 (`board.get`) first; then T-3249 (the board) and T-3250 (the documents) together.
- Forecast 42m, as `flai forecast` gives it for size 23.
- Cost of delay: 74.82 USD a week, as `flai cod` gives it: this story's share of E-0034's 1500.00 USD a week, the operator's `time_lost_per_cycle: 10h` (TH-1115) at 150 USD an hour, by its forecast of 42m over the 14h2m of the epic's 19 open stories. Kept: the epic's outcome needs every story, so a share by size is the split that fits.

### Accepted by the orchestrator

- Verified: 38043c359c1be2f15f3f65d7965600713a48e6ab
- At: 2026-11-17T04:46:04Z

Verdict: meets, all six criteria met at 38043c359c1be2f15f3f65d7965600713a48e6ab; close-out passed every step there.
- 1: flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go
- 2: flaiover/src/lib/server/board.ts, flaiover/src/lib/lanes.ts, flaiover/src/lib/lanecollapse.svelte.ts, flaiover/src/lib/reorder.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/+page.svelte, flaiover/src/lib/components/CardMenu.svelte, flaiover/src/lib/components/LaneMenu.svelte, flaiover/src/lib/components/LaneMoveDialog.svelte, flai/internal/hostapi/writes_items.go
- 3: flaiover/src/routes/board/board.svelte.test.ts, flaiover/src/lib/lanes.ts, flai/internal/workitem/boardview.go
- 4: flaiover/src/lib/lanes.fixture.ts, flaiover/src/routes/board/board.svelte.test.ts, flaiover/src/lib/lanes.test.ts, flaiover/src/lib/lanecollapse.svelte.test.ts, flaiover/src/lib/reorder.test.ts
- 5: flaiover/src/lib/lanes.fixture.ts, flaiover/src/lib/server/repo-channel.test.ts, flaiover/src/lib/server/writes.test.ts, flaiover/src/routes/board/cardmenu.svelte.test.ts, flaiover/src/routes/board/lanemenu.svelte.test.ts, flaiover/src/routes/overview.svelte.test.ts
- 6: design/system/flaiover-dashboard.md, docs/users/flaiover.md
