---
id: T-3249
type: task
nature: feature
title: flaiover's board, lane and card menus, drag and drop, reason prompt, reordering, collapsing, and home page counts read the columns and moves from board.get
status: done
parent: S-0753
owner: alex
created: 2026-11-12T16:09:32Z
updated: 2026-11-17T04:14:52Z
transitions:
  - to: ready
    at: 2026-11-17T03:54:51Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T03:54:51Z
    by: system-flow
  - to: done
    at: 2026-11-17T04:14:52Z
    by: agent-S-0753
stream: S-0753
tags: [dashboard]
touches: [flaiover/src/lib/server/board.ts, flaiover/src/lib/lanes.ts, flaiover/src/lib/lanes.test.ts, flaiover/src/lib/lanecollapse.svelte.ts, flaiover/src/lib/lanecollapse.svelte.test.ts, flaiover/src/lib/reorder.ts, flaiover/src/lib/reorder.test.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, flaiover/src/routes/board/lanemenu.svelte.test.ts, flaiover/src/routes/+page.svelte, flaiover/src/lib/components/CardMenu.svelte, flaiover/src/lib/components/LaneMenu.svelte, flaiover/src/lib/components/LaneMoveDialog.svelte, flaiover/src/lib/components/NewItemForm.svelte, flaiover/src/lib/components/NewItemForm.svelte.test.ts, flaiover/src/lib/lanes.fixture.ts, flaiover/src/lib/server/repo-channel.test.ts, flaiover/src/lib/server/writes.test.ts, flaiover/src/routes/board/cardmenu.svelte.test.ts, flaiover/src/routes/new/+page.svelte, flaiover/src/routes/overview.svelte.test.ts]
after: [T-3244]
usage:
  source: log
  seconds: 1201
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 160
      output: 54000
      cache_read: 13509057
      cache_write: 328293
      cost: 5.7751
---
# T-3249 flaiover's board, lane and card menus, drag and drop, reason prompt, reordering, collapsing, and home page counts read the columns and moves from board.get

## Work

Waits for T-3244: it reads the columns and transitions that task adds to `board.get`. Remove `STATES` from `board.ts`, `LANES`, `START_LANES`, `back`, `forwardOf`, and `needsReason` from `lanes.ts`, `laneStates` from `lanecollapse.svelte.ts`, `ORDERED` from `reorder.ts`, and the lists in the two pages, and draw everything from the answer: a lane per column in order, its WIP limit, the moves the operator may make, the reason prompt where a transition needs one, and the ordered lanes. A move flai refuses still shows flai's words. Test a board for a workflow with a state added between `in-progress` and `review`.

## Done when

None of these files holds a list of states, and the flaiover tests pass on the baseline and on the other workflow.

## Notes
