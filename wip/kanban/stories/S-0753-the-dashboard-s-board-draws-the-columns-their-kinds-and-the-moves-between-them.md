---
id: S-0753
type: story
nature: feature
title: The dashboard's board draws the columns, their kinds, and the moves between them from what flai gives, and holds no list of states of its own
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:01:41Z
updated: 2026-11-12T16:01:41Z
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
draft: true
---
# S-0753 The dashboard's board draws the columns, their kinds, and the moves between them from what flai gives, and holds no list of states of its own

## Goal

The dashboard draws whatever board the workflow defines. `board.get` gives the columns in order with each one's kind, WIP limit, and order flag, and the moves between them with who may make each and whether it needs a reason, and flaiover draws the lanes, the card and lane menus, the drag targets, the reason prompt, and the collapsed lanes from that answer instead of its own lists in `board.ts`, `lanes.ts`, `lanecollapse.svelte.ts`, `reorder.ts`, and the two pages.

## Acceptance criteria

- [ ] `board.get` answers the workflow's columns, each with its name, kind, WIP limit, and whether it is ordered, and its transitions, each with from, to, who may make it, and whether it needs a reason
- [ ] `flaiover/src` holds no list of the six states: the board, the lane and card menus, drag and drop, the reason prompt, reordering, collapsing, and the home page's counts read the columns and transitions from `board.get`
- [ ] A move the operator may not make from the board is not offered, and one flai refuses still shows flai's words
- [ ] Tests draw a board for a workflow with a state renamed and one added between `in-progress` and `review`, with its lane, its WIP limit, and its moves
- [ ] With the baseline, the board looks and behaves as today, and every existing flaiover test passes, its fixtures updated only for the new fields
- [ ] `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` say the board's columns come from the workflow

## Tasks

## Notes

- The lists today: `flaiover/src/lib/server/board.ts` (`STATES`), `lanes.ts` (`LANES`, `START_LANES`, `back`, `forwardOf`, `needsReason`), `lanecollapse.svelte.ts` (`laneStates`), `reorder.ts` (`ORDERED`), `routes/board/+page.svelte`, and `routes/+page.svelte`. `board.get` is in `flai/internal/hostapi/hostapi.go` and returns `workitem`'s board view.
