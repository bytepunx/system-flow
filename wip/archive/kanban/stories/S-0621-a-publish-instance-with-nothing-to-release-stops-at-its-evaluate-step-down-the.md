---
id: S-0621
type: story
nature: improvement
title: A publish instance with nothing to release stops at its evaluate step down the failure path, blocking and asking the operator every time
status: done
owner: alex
created: 2026-11-07T23:39:50Z
updated: 2026-11-09T21:52:17Z
transitions:
  - to: ready
    at: 2026-11-09T17:09:53Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T17:41:54Z
    by: agent-S-0621
  - to: review
    at: 2026-11-09T18:31:02Z
    by: agent-S-0621
  - to: done
    at: 2026-11-09T21:52:17Z
    by: alex
tags: [flai, template]
topics: [formulas, release, acceptance]
touches: [design/adrs/README.md, design/system/strategic-agents.md, design/system/flai-cli.md, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/check/formulas_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/cmd/formula_stop.go, flai/cmd/formula_stop_test.go, flai/internal/guard/guard.go, flai/internal/guard/formula_test.go, flai/internal/preview/orchestrator.go, flai/internal/preview/accept_test.go, flai/cmd/accept.go, flai/cmd/accept_formula_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, design/workflows/formulas/README.md, template/root/design/workflows/formulas/README.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, template/CHANGELOG.md, template/template.yaml, design/issues/I-0283-a-publish-instance-with-nothing-to-release-stops-at-its-evaluate-step-down-the.md, design/issues/summary.md, design/adrs/0271-a-formula-step-may-name-the-exits-that-stop-its-instance-cleanly-flai-formula.md, docs/operators/settings.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2954
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 3
      hand_edits: 1
      work: 54
  models:
    - model: claude-opus-5-5
      input: 498
      output: 132788
      cache_read: 28878077
      cache_write: 618734
      cost: 12.2889
  strategic:
    - kind: orchestrator
      seconds: 54
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 12
          output: 94
          cache_read: 697297
          cache_write: 52561
          cost: 0.1981
cost_of_delay:
  inputs:
    time_lost_per_cycle: 24m
    by: flai
    at: 2026-11-07T23:39:50Z
  value: 60
  by: planner-S-0621
  at: 2026-11-09T17:06:05Z
forecast:
  duration: 1h30m
  delivery: 2026-11-09T19:50:00Z
  basis: "Its own forecast of 1h30m; 4th in the pull order with an in-progress limit of 5, behind S-0611, S-0616, S-0619, S-0620, S-0624, S-0629, S-0626 and S-0627."
  by: flai
  at: 2026-11-09T17:18:47Z
finalized:
  by: orchestrator
  at: 2026-11-09T17:09:29Z
---
# S-0621 A publish instance with nothing to release stops at its evaluate step down the failure path, blocking and asking the operator every time

## Goal

This story remediates [I-0283](../../../design/issues/I-0283-a-publish-instance-with-nothing-to-release-stops-at-its-evaluate-step-down-the.md), "A publish instance with nothing to release stops at its evaluate step down the failure path, blocking and asking the operator every time". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0283 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0283 is closed with `flai issue close I-0283 --reason` saying what fixed it

## Tasks
- T-2743 An ADR and the design say a formula step may name the exits that stop its instance cleanly, and how the instance then reaches done
- T-2744 A formula step's stop lists the exits that stop its instance cleanly, flai check validates it, and run and show write it
- T-2745 flai formula stop closes the stopping step's task, cancels the steps after it, and ticks their criteria as skipped, so the instance can go to review
- T-2746 The orchestrator accepts an instance stopped cleanly, taking a step task cancelled by the stop as its criterion's evidence
- T-2747 The instance prompt and the conventions in both copies say a stop exit ends the run with flai formula stop, asking nothing
- T-2748 The publish formula's evaluate step stops on exit 2 in both copies, the formula and user docs say how a stop works, and I-0283 is closed

## Notes

Cost of delay inputs set by flai from I-0283. time_lost_per_cycle 24m: 8m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-11-07T08:16:08Z, 0.6 days before this story; under one cycle counts as one).

### Planning

The proposed fix, from I-0283's three instances and the operator's answer on TH-0854 (D-0498). A formula step may name the exits that stop its instance cleanly, `stop: [2]` on the publish formula's evaluate step. When such a step exits with one of them, the story's agent runs `flai formula stop <S-nnnn> --step <id> --exit <n>` with the step's output. That closes the step's task as done, with the exit and the output in its Notes. It cancels each step task after it with the reason `step <id> stopped the run with exit <n>`. It ticks their criteria, each marked skipped with that reason. The agent then closes the instance out and moves it to review, asking nothing. The orchestrator's acceptance takes a step task cancelled by a stop as its criterion's evidence. The planner proposes the field's name and the command; the story's agent settles them in the ADR its first task writes.

Touches, by where they came from:

| Source | Paths |
|--------|-------|
| Design | `design/system/strategic-agents.md` (Formulas, Working an instance), `design/system/flai-cli.md` (Commands), `design/adrs/README.md` (the new ADR's index row; the ADR's own file is widened by `task_done`), the two copies each of `design/conventions/work-management.md`, `design/conventions/strategic-agents.md`, `design/workflows/formulas/publish.md`, and `design/workflows/formulas/README.md` |
| Layout (the explorer's map of the code) | `flai/internal/formulas/formulas.go` and its test (the `Step` struct and its checks), `flai/internal/check/formulas_test.go`, `flai/cmd/formula.go` and its test (`formulaTaskSections`, `printFormula`, the subcommand), `flai/cmd/formula_stop.go` and its test (new), `flai/internal/guard/guard.go` and `formula_test.go` (`cliAuthors`), `flai/internal/preview/orchestrator.go` (`stepsUndone`) and `accept_test.go`, `flai/cmd/accept.go` and `accept_formula_test.go`, `flai/internal/harness/harness.go` and its test (the `instance` paragraph) |
| Docs | `docs/users/flai.md`, `docs/users/flai-reference.md` (generated by `make flai-reference`), `docs/operators/index.md` (Writing and running a formula) |
| Co-change | `template/template.yaml` (the template's version bump beside `template/CHANGELOG.md`); `flai touches suggest` listed nothing else over 9% |
| The criteria | the I-0283 issue file and `design/issues/summary.md`, for its close |

No folder touch. No `after`: S-0625, the judgement policy's half of TH-0811, changes `flai release --evaluate`, which this story only reads as it stands on main.

Figures:

- Forecast 1h30m, delivery 2026-11-10T17:00Z. `flai forecast` gave 57m from size 34. The three formula stories, S-0421 to S-0423, took 46 to 59 minutes each, and this one does more: an ADR, a new command, an acceptance rule, and the docs in two copies, in six tasks. The delivery leaves a day for the operator to finalize the draft.
- Cost of delay 60 USD a week, as `flai cod` worked it from the inputs flai set from I-0283: 24m lost per 168h cycle at 150 USD an hour. It stands: three operator interventions per stopped instance (TH-0854, the criterion tick, the hand acceptance) match the issue's 8m an occurrence.
