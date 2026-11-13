---
id: S-0753
type: story
nature: feature
title: The dashboard's board draws the columns, their kinds, and the moves between them from what flai gives, and holds no list of states of its own
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:01:41Z
updated: 2026-11-13T15:47:58Z
transitions: []
tags: [dashboard, cli]
topics: [dashboard, cli]
touches: [flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flaiover/src/lib/server/board.ts, flaiover/src/lib/lanes.ts, flaiover/src/lib/lanes.test.ts, flaiover/src/lib/lanecollapse.svelte.ts, flaiover/src/lib/lanecollapse.svelte.test.ts, flaiover/src/lib/reorder.ts, flaiover/src/lib/reorder.test.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, flaiover/src/routes/board/lanemenu.svelte.test.ts, flaiover/src/routes/+page.svelte, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0739]
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
      seconds: 13
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 15
          cache_read: 121381
          cache_write: 20507
          cost: 0.0405
draft: true
forecast:
  duration: 42m
  delivery: 2026-11-13T21:42:00Z
  basis: "Its own forecast of 42m; 20th in the pull order with an in-progress limit of 3, behind S-0751, S-0760, S-0759, S-0757, S-0241, S-0730, S-0736, S-0737, S-0738, S-0739, S-0741, S-0742, S-0743, S-0744, S-0746, S-0747, S-0748, S-0749, S-0750 and S-0752."
  by: flai
  at: 2026-11-13T15:47:58Z
---
# S-0753 The dashboard's board draws the columns, their kinds, and the moves between them from what flai gives, and holds no list of states of its own

## Goal

The dashboard draws whatever board the workflow defines. `board.get` gives the columns in order with each one's kind, WIP limit, and order flag, and the moves between them with who may make each and whether it needs a reason, and flaiover draws the lanes, the card and lane menus, the drag targets, the reason prompt, and the collapsed lanes from that answer instead of its own lists in `board.ts`, `lanes.ts`, `lanecollapse.svelte.ts`, `reorder.ts`, and the two pages.

## Acceptance criteria

- [ ] `board.get` answers the workflow's columns, each with its name, kind, WIP limit, and whether it is ordered, and its transitions, each with from, to, who may make it, and whether it needs a reason
- [ ] `board.ts`, `lanes.ts`, `lanecollapse.svelte.ts`, `reorder.ts`, `routes/board/+page.svelte`, and `routes/+page.svelte` hold no list of the six states: the board, the lane and card menus, drag and drop, the reason prompt, reordering, collapsing, and the home page's counts read the columns and transitions from `board.get`
- [ ] A move the operator may not make from the board is not offered, and one flai refuses still shows flai's words
- [ ] Tests draw a board for a workflow with a state renamed and one added between `in-progress` and `review`, with its lane, its WIP limit, and its moves
- [ ] With the baseline, the board looks and behaves as today, and every existing flaiover test passes, its fixtures updated only for the new fields
- [ ] `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` say the board's columns come from the workflow

## Tasks

- T-3244 board.get answers the workflow's columns with their kinds, WIP limits, and order, and its transitions with who may make each and its reason
- T-3249 flaiover's board, lane and card menus, drag and drop, reason prompt, reordering, collapsing, and home page counts read the columns and moves from board.get
- T-3250 flaiover-dashboard.md and the dashboard's user guide say the board's columns and moves come from the workflow

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
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
