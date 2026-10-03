---
id: S-0204
type: story
nature: feature
title: The new-item and edit forms take cost of delay inputs in a collapsed panel
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:12Z
updated: 2026-10-03T05:34:24Z
transitions:
  - to: ready
    at: 2026-10-03T05:34:24Z
    by: alex
tags: [dashboard, flai]
touches: [flaiover/src/routes/new, flaiover/src/routes/edit, flaiover/src/routes/items, flai/internal/hostapi]
after: [S-0199]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0204 The new-item and edit forms take cost of delay inputs in a collapsed panel

## Goal

The operator supplies what delay costs when creating or editing an epic or story: anticipated revenue per week, anticipated penalty per week, or time lost per cycle. The planner derives a cost of delay value from them. The panel is collapsed by default so the forms stay short.

## Acceptance criteria
- [ ] The new-item and edit forms for epics and stories have a **Cost of delay** panel, collapsed by default, with fields for revenue per week, penalty per week (both in the project's currency, from the manifest), and time lost per cycle (a duration); each is optional and the panel shows a summary line when collapsed and any value is set
- [ ] Saving writes `cost_of_delay.inputs` through hostapi with `by` the operator; clearing all three removes the inputs and leaves the planner's `value` to be recomputed
- [ ] The item page shows the inputs, the planner's `value` with its `by` and `at`, and a note when the value is stale (inputs newer than the value)
- [ ] `design/system/flaiover-dashboard.md` and the user guide describe the panel; tests cover saving, clearing, and the stale note

## Tasks

## Notes
