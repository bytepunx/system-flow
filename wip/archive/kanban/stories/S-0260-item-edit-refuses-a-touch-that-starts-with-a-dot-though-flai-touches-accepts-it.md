---
id: S-0260
type: story
nature: remediation
title: item_edit refuses a touch that starts with a dot, though flai touches accepts it
status: done
owner: alex
created: 2026-10-04T06:29:48Z
updated: 2026-10-05T03:46:30Z
transitions:
  - to: ready
    at: 2026-10-04T23:14:26Z
    by: alex
  - to: in-progress
    at: 2026-10-05T03:14:11Z
    by: agent-S-0260
  - to: review
    at: 2026-10-05T03:26:42Z
    by: agent-S-0260
  - to: done
    at: 2026-10-05T03:46:30Z
    by: alex
tags: []
touches: [flai/internal/workitem/create.go, flai/internal/workitem/create_test.go, flai/internal/itemedit/itemedit.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/items_write_test.go, flai/cmd/edit_test.go, flai/cmd/touches.go, flai/cmd/touches_test.go, design/system/work-hierarchy.md, design/issues, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 809
  models:
    - model: claude-opus-5-5
      input: 182
      output: 45584
      cache_read: 8040530
      cache_write: 245325
      cost: 4.1502
    - model: claude-sonnet-5
      input: 36
      output: 7572
      cache_read: 845245
      cache_write: 55735
      cost: 0.3842
  strategic:
    - kind: planner
      seconds: 202
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 378
          output: 107
          cache_read: 3131614
          cache_write: 93189
          cost: 0.6938
        - model: claude-opus-5-5
          input: 251
          output: 22496
          cache_read: 6544260
          cache_write: 220332
          cost: 3.8994
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: flai
    at: 2026-10-04T06:29:48Z
  value: 5
  by: planner-S-0260
  at: 2026-10-05T03:13:22Z
forecast:
  duration: 13m
  delivery: 2026-10-05T03:32:00Z
  basis: "flai forecast: median 60 s per unit of size over 8 done remediation stories on claude-opus-5-5 in the large band, times size 13 (2 criteria, 11 touches), 5th in the pull order behind S-0253, S-0257, S-0244, S-0262 and S-0258."
  by: planner-S-0260
  at: 2026-10-05T03:13:22Z
finalized:
  by: alex
  at: 2026-10-04T23:14:24Z
---
# S-0260 item_edit refuses a touch that starts with a dot, though flai touches accepts it

## Goal

This story remediates [I-0071](../../../design/issues/I-0071-item-edit-refuses-a-touch-that-starts-with-a-dot-though-flai-touches-accepts-it.md), "item_edit refuses a touch that starts with a dot, though flai touches accepts it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0071 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0071 is closed with `flai issue close I-0071 --reason` saying what fixed it

## Tasks
- T-0862 item_edit and the dashboard's writes accept a touch that starts with a dot, by one touch rule in workitem
- T-0863 flai touches, item_new, and flai story and task new check touches by the same rule
- T-0864 work-hierarchy.md states the touch rule, and I-0071 is closed
- T-0865 One touch rule for every writer: a path starting with a dot is accepted, one that escapes the repository is refused
- T-0866 The design records the touch rule and I-0071 is closed

## Notes

Cost of delay inputs set by flai from I-0071. time_lost_per_cycle 2m: 2m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-04T04:45:59Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0260 on 2026-10-05.

The cause: `itemedit.cleanList` checks touches with `listValue`, `^[A-Za-z0-9][A-Za-z0-9 _./@+-]*$`, so a touch must start with a letter or a digit. The refusal itself says only "not starting with a dash". The dashboard's writes in `flai/internal/hostapi/writes.go` copy the rule. `flai touches` and `workitem.Create`, behind `item_new`, check nothing.

The proposed fix is one touch rule, `workitem.CleanTouches`, used by every way of setting touches. It accepts a leading dot. It refuses a leading dash, an absolute path, and a `..` segment. Tags keep their rule.

The tasks are duplicated. I drafted T-0862, T-0863, and T-0864 at 03:13Z. Setting the touches let agent-S-0260 pull the story at 03:14Z, and it wrote T-0865 and T-0866 for the same fix. T-0865 covers T-0862 and T-0863, and T-0866 covers T-0864. The plan's thread proposes dropping T-0862, T-0863, and T-0864.

Touches. The story declared none, so `flai touches suggest` had nothing to start from. I seeded it with the two files that hold the rule.

- Layout, read from the code:
  - `flai/internal/itemedit/itemedit.go` has `cleanList` and `listValue`.
  - `flai/internal/hostapi/writes.go` has the dashboard's copy.
  - `flai/internal/workitem/create.go` has `CleanTopics`, where `CleanTouches` goes, and `Create`, whose touches go unchecked.
  - `flai/cmd/touches.go` is `flai touches`.
- Co-change: `flai touches suggest` from those seeds gave `flai/internal/hostapi/writes_test.go` (67%). The other paths it listed are the dashboard's agent, serve actions, and generated docs, which this fix does not reach.
- Layout, the tests that pin the rule:
  - `flai/cmd/edit_test.go` asserts "not starting with a dash".
  - `flai/internal/mcpserver/items_write_test.go` tests `item_edit` and `item_new`.
  - `flai/internal/workitem/create_test.go` and `flai/cmd/touches_test.go` test the other two ways of setting touches.
- Design:
  - `design/system/work-hierarchy.md` describes `touches`, and should state the rule.
  - `design/issues` covers closing I-0071 and `summary.md`.
- `docs/users/flai-reference.md` is left out: no help text needs to change.

Forecast: 13m, as `flai forecast` gave it once the touches were set. That is the median of 60 s per unit of size over 8 done remediation stories on claude-opus-5-5, times size 13 (2 criteria and 11 touches). It stands because recent remediation stories of similar size took 8m (S-0268, 10 touches), 29m (S-0249, 24 touches), and 43m (S-0266, 14 touches), and 13m lies within that range. The delivery of 03:32Z was worked out with the story 5th in the pull order. It was pulled at 03:14Z, so it may land sooner.

Cost of delay: 5.00 USD a week, as `flai cod` worked it out from flai's input of 2m per 168h cycle at 150 USD an hour. It stands: I-0071 has one occurrence, and `flai touches` is a workaround that costs a retry.
