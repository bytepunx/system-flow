---
id: S-0620
type: story
nature: remediation
title: A formula's instance declares no touches, so the story that started it holds it until that story leaves in-progress
status: in-progress
owner: alex
created: 2026-11-07T23:39:44Z
updated: 2026-11-09T17:10:27Z
transitions:
  - to: ready
    at: 2026-11-09T17:09:52Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T17:10:27Z
    by: agent-S-0620
tags: [flai, formulas, holds]
topics: [formulas, claims]
touches: [flai/cmd/formula.go, flai/cmd/formula_test.go, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0282-a-formula-s-instance-declares-no-touches-so-the-story-that-started-it-holds-it.md, design/issues/summary.md]
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
      seconds: 28
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 66
          cache_read: 515889
          cache_write: 24168
          cost: 0.1409
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-07T23:39:44Z
  value: 12.5
  by: planner-S-0620
  at: 2026-11-09T17:05:44Z
forecast:
  duration: 30m
  delivery: 2026-11-09T17:40:00Z
  basis: flai forecast's 18m (median 120 s per unit over 44 done remediation stories on claude-opus-5-5, size 9) raised to 30m for the proposal the goal asks for before the build and a test that reproduces the hold
  by: planner-S-0620
  at: 2026-11-09T17:05:44Z
finalized:
  by: orchestrator
  at: 2026-11-09T17:09:28Z
---
# S-0620 A formula's instance declares no touches, so the story that started it holds it until that story leaves in-progress

## Goal

This story remediates [I-0282](../../../design/issues/I-0282-a-formula-s-instance-declares-no-touches-so-the-story-that-started-it-holds-it.md), "A formula's instance declares no touches, so the story that started it holds it until that story leaves in-progress". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0282 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0282 is closed with `flai issue close I-0282 --reason` saying what fixed it

## Tasks
- T-2740 flai formula run gives the instance story its own story file as its touch, with a test that the instance is not held by a story in progress
- T-2741 Say in the design and the user guide that a formula's instance claims its own story file
- T-2742 Close I-0282 saying that flai formula run now gives the instance its own story file as its touch

## Notes

Cost of delay inputs set by flai from I-0282. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T08:04:18Z, 0.6 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0282's instance and D-0509: `flai formula run` gives the instance story its own story file as its touch before `--ready` moves it, so its claim is not empty. A ready story with an empty claim is held `no-touches` by every story in progress, the one that started it too (`Holds.overlap`, `flai/internal/workitem/hold.go`). An instance in progress with an empty claim also holds every ready story. A claim on its own story file overlaps no other story, so neither hold happens. The other fix the issue names, a hold that passes a story with `formula:`, was not chosen: it changes the hold rule that `metrics.md` and ADR-0046 state, and an instance whose steps do change files would then claim nothing. `formulaRun` is the one path that makes an instance: `flai formula run`, `flai import` (`flai/cmd/import_commit.go`), and the orchestrator's publish all go through it, so one change covers them all.

Touches, all files, none a folder:

- `flai/cmd/formula.go`, `flai/cmd/formula_test.go`: layout. `formulaRun` makes the story through `itemnew.Create` with no `Touches`; `TestFormulaRun` and `TestFormulaRunReady` assert nothing about touches or holds.
- `design/system/strategic-agents.md` (§ Formulas, "What run makes"), `design/system/flai-cli.md` (the `flai formula run` row), `docs/users/flai.md` (§ A story from a formula): design. Each states what run makes and gains the instance's touch. The last two are in `claims.shared`.
- `design/issues/I-0282-…md`, `design/issues/summary.md`: the second criterion. `flai issue close` writes both. Both are in `claims.shared`.

`flai touches suggest S-0620` listed nothing, because the story declared no touches.

Forecast: `flai forecast` gave 5m from size 2 (2 criteria, 0 touches) before the touches were set, and 18m with the 7 touches (median 120 s per unit of size over 44 done remediation stories, size 9). It was raised to 30m because the goal asks for the fix to be proposed before it is built, and the test must reproduce the hold with a story in progress. Delivery 2026-11-09T17:40:00Z keeps flai's 3rd place in the pull order.

Cost of delay: `flai cod` gives 12.50 USD a week, from the 5m `time_lost_per_cycle` input flai set from I-0282 at 150 USD an hour. That value stands. The instance names two stories held the same way, S-0584 and S-0585, and every `publish` instance the orchestrator starts with `--ready` meets it. The input is the operator's, so it was not changed.

No `after`: the story reads nothing that an open story writes.
