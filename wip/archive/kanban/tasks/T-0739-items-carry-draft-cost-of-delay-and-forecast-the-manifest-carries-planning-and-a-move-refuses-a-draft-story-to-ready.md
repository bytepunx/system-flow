---
id: T-0739
type: task
nature: feature
title: Items carry draft, cost_of_delay, and forecast, the manifest carries planning, and a move refuses a draft story to ready
status: done
parent: S-0199
owner: alex
created: 2026-10-03T05:41:02Z
updated: 2026-10-03T05:47:57Z
transitions:
  - to: ready
    at: 2026-10-03T05:41:45Z
    by: agent-S-0199
  - to: in-progress
    at: 2026-10-03T05:41:45Z
    by: agent-S-0199
  - to: done
    at: 2026-10-03T05:47:57Z
    by: agent-S-0199
stream: S-0199
tags: []
touches: [flai/internal/workitem, flai/internal/manifest, docs/operators/settings.md]
usage:
  source: log
  seconds: 372
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 53
      output: 19553
      cache_read: 3161596
      cache_write: 79759
      cost: 1.494
    - model: claude-sonnet-5
      input: 50
      output: 14308
      cache_read: 1474912
      cache_write: 117755
      cost: 0.7325
---
# T-0739 Items carry draft, cost_of_delay, and forecast, the manifest carries planning, and a move refuses a draft story to ready

## Work

Add to `workitem.Item`: `draft` (stories), `cost_of_delay` (epics and stories: `inputs` with `revenue_per_week`, `penalty_per_week`, `time_lost_per_cycle`; `value`; `by`; `at`), and `forecast` (stories: `duration`, `delivery`, `basis`, `by`, `at`). Write them in `Marshal` after `estimate`, only when set; validate their shapes in `Validate` (non-negative amounts, Go durations, timestamps, `by` and `at` present on a set block, the types that carry each). List them in `front-matter-fields.txt` with their nested keys and `item.types`. Add `planning` to the manifest: `currency` (ISO 4217, default USD), `hour_rate`, `cycle` (Go duration, default 168h), with accessors and a row each in `docs/operators/settings.md`. `workitem.Move` refuses a draft story to `ready` with "finalize it first" unless `MoveOptions.Finalize`, which clears `draft` in the same save; `Transition` and `TransitionAll` carry it. Waits for nothing: every other task builds on these types.

## Done when

The fields round-trip, invalid shapes are refused by `Validate`, the fields file test passes, and a draft story's move to ready is refused without `Finalize` and clears `draft` with it, each with a test.

## Notes
