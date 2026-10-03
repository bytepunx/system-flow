---
id: T-0767
type: task
nature: feature
title: flai story new and flai epic new take cost of delay inputs, and hostapi item.new passes them
status: done
parent: S-0204
owner: alex
created: 2026-10-03T18:35:38Z
updated: 2026-10-03T18:47:24Z
transitions:
  - to: ready
    at: 2026-10-03T18:36:18Z
    by: agent-S-0204
  - to: in-progress
    at: 2026-10-03T18:36:19Z
    by: agent-S-0204
  - to: done
    at: 2026-10-03T18:47:24Z
    by: agent-S-0204
stream: S-0204
tags: []
touches: [flai/cmd, flai/internal/hostapi, docs/users/flai.md, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 665
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 106
      output: 711
      cache_read: 4724647
      cache_write: 120262
      cost: 1.9633
    - model: claude-sonnet-5
      input: 32
      output: 76
      cache_read: 396281
      cache_write: 42557
      cost: 0.1742
---
# T-0767 flai story new and flai epic new take cost of delay inputs, and hostapi item.new passes them

## Work

Give `flai story new` and `flai epic new` the `--revenue-per-week`, `--penalty-per-week`, and `--time-lost-per-cycle` flags `flai edit` has, passed to `workitem.CreateOptions.CostOfDelay` (create.go takes one since S-0203), and give hostapi `item.new` a `cost_of_delay` object with those three keys as strings, turned into the flags, so the new-item form saves the inputs in the commit that makes the item, stamped with the operator (`--by` or the owner, as item.edit does). Waits for nothing: the first layer, beside the dashboard panel, which shares no path with it.

## Done when

- `flai story new --revenue-per-week 1200 ...` writes `cost_of_delay.inputs` with `by` and `at`; a bad amount or duration is refused with the edit's message.
- hostapi `item.new` with `cost_of_delay` builds those flags, empty keys are left out, and a task type is refused; tests in `flai/internal/hostapi` and `flai/cmd` cover it.
- `docs/users/flai.md` and the generated `docs/users/flai-reference.md` name the flags.

## Notes
