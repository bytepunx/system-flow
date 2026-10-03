---
id: S-0204
type: story
nature: feature
title: The new-item and edit forms take cost of delay inputs in a collapsed panel
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:12Z
updated: 2026-10-03T20:06:55Z
transitions:
  - to: ready
    at: 2026-10-03T05:34:24Z
    by: alex
  - to: in-progress
    at: 2026-10-03T18:31:22Z
    by: agent-S-0204
  - to: review
    at: 2026-10-03T19:42:52Z
    by: agent-S-0204
  - to: done
    at: 2026-10-03T20:06:55Z
    by: alex
tags: [dashboard, flai]
touches: [flaiover/src/routes/new, flaiover/src/routes/edit, flaiover/src/routes/items, flaiover/src/routes/api/items/+server.ts, flaiover/src/lib/components/CostOfDelayPanel.svelte, flaiover/src/lib/components/NewItemForm.svelte, flaiover/src/lib/components/ItemEditor.svelte, flaiover/src/lib/costofdelay.ts, flaiover/src/lib/planning.ts, flaiover/src/lib/planning.test.ts, flai/internal/hostapi, flai/cmd, flai/internal/workitem, flai/internal/itemedit, design/adrs, design/system/work-hierarchy.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md, docs/users/flai-reference.md, "flaiover/src/routes/api/items/[id]/edit", flaiover/src/lib/costofdelay.test.ts, flaiover/src/lib/components/CostOfDelayPanel.svelte.test.ts, flaiover/src/lib/components/NewItemForm.svelte.test.ts, flaiover/src/lib/components/ItemEditor.svelte.test.ts, design/issues/I-0066-an-acceptance-merge-committed-conflict-markers-to-a-design-document-on-main-and-nothing-caught-it.md, design/issues/summary.md, docs/operators/settings.md, flai/internal/issues, flai/internal/mcpserver, design/issues/I-0063-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md]
after: [S-0199]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4322
  models:
    - model: claude-haiku-4-5-20251001
      input: 170
      output: 9502
      cache_read: 1222239
      cache_write: 99386
      cost: 0.2941
    - model: claude-opus-5-5
      input: 598
      output: 159081
      cache_read: 34004920
      cache_write: 634355
      cost: 13.8927
    - model: claude-sonnet-5
      input: 104
      output: 25113
      cache_read: 2977580
      cache_write: 166362
      cost: 1.2628
---
# S-0204 The new-item and edit forms take cost of delay inputs in a collapsed panel

## Goal

The operator supplies what delay costs when creating or editing an epic or story: anticipated revenue per week, anticipated penalty per week, or time lost per cycle. The planner derives a cost of delay value from them. The panel is collapsed by default so the forms stay short.

## Acceptance criteria
- [x] The new-item and edit forms for epics and stories have a **Cost of delay** panel, collapsed by default, with fields for revenue per week, penalty per week (both in the project's currency, from the manifest), and time lost per cycle (a duration); each is optional and the panel shows a summary line when collapsed and any value is set
- [x] Saving writes `cost_of_delay.inputs` through hostapi with `by` the operator; clearing all three removes the inputs and leaves the planner's `value` to be recomputed
- [x] The item page shows the inputs, the planner's `value` with its `by` and `at`, and a note when the value is stale (inputs newer than the value)
- [x] `design/system/flaiover-dashboard.md` and the user guide describe the panel; tests cover saving, clearing, and the stale note

## Tasks
- T-0767 flai story new and flai epic new take cost of delay inputs, and hostapi item.new passes them
- T-0768 The new-item and edit forms have a collapsed Cost of delay panel
- T-0769 Cost of delay inputs carry their own by and at, apart from the value's
- T-0770 The item page shows the inputs, the planner's value with who set it and when, and a stale note
- T-0771 The dashboard design and the user guide describe the Cost of delay panel and the item page

## Notes
