---
id: S-0622
type: story
nature: improvement
title: orchestration.permissions.set_cost_of_delay is off, so a cost of delay input the operator confirmed waits a day for the operator to write it
status: backlog
owner: alex
created: 2026-11-07T23:39:56Z
updated: 2026-11-09T21:22:32Z
transitions: []
tags: [cli]
topics: [orchestration, planning]
touches: [flai/internal/serve/orchestrate_confirm_test.go, flai/cmd/orchestrate_permissions_test.go, design/issues/I-0284-orchestration-permissions-set-cost-of-delay-is-off-so-a-cost-of-delay-input-the.md, design/issues/summary.md]
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
    - kind: planner
      seconds: 562
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 86
          output: 51823
          cache_read: 4575916
          cache_write: 180285
          cost: 0.3609
        - model: claude-opus-5-5
          input: 68
          output: 21748
          cache_read: 4988171
          cache_write: 174061
          cost: 2.8254
    - kind: orchestrator
      seconds: 102
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 42
          output: 268
          cache_read: 1558367
          cache_write: 213345
          cost: 0.4696
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 33
          output: 3103
          cache_read: 238317
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2841
          cache_read: 759050
          cache_write: 13855
          cost: 0.278
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 75h
    by: flai
    at: 2026-11-07T23:39:56Z
  value: 11250
  by: planner-S-0622
  at: 2026-11-09T18:04:50Z
forecast:
  duration: 45m
  delivery: 2026-11-09T22:33:00Z
  basis: "Its own forecast of 45m; 3rd in the pull order with an in-progress limit of 5, behind S-0616, S-0642 and S-0241."
  by: flai
  at: 2026-11-09T21:22:32Z
---
# S-0622 orchestration.permissions.set_cost_of_delay is off, so a cost of delay input the operator confirmed waits a day for the operator to write it

## Goal

This story remediates [I-0284](../../../design/issues/I-0284-orchestration-permissions-set-cost-of-delay-is-off-so-a-cost-of-delay-input-the.md), "orchestration.permissions.set_cost_of_delay is off, so a cost of delay input the operator confirmed waits a day for the operator to write it", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0284 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0284 is closed with `flai issue close I-0284 --reason` saying what fixed it

## Tasks
- T-2763 Test that the operator's confirmation of a cost of delay recommendation on a planner's thread gives the orchestrator a thread reason
- T-2764 Test that the orchestrator sets an in-progress story's cost of delay input only while set_cost_of_delay is on
- T-2765 Close I-0284 on the manifest change that turned set_cost_of_delay on, naming the two regression tests

## Notes

The analyzer found I-0284: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0284. time_lost_per_cycle 75h carried over from I-0284's Impact section.

Strategic usage carried over from I-0284: analyzer 42s, 0.8156 USD.

### Planning

**Proposed solution.** I-0284's one instance, S-0547 and S-0549 waiting on TH-0784 and TH-0789 for an input the operator confirmed on TH-0775, came from configuration: `orchestration.permissions` had no `set_cost_of_delay`. The operator turned it on in commit c76b23a2d at 2026-11-07T09:17:55Z, an hour after the analysis. Since then no confirmed input has waited: TH-0884, TH-0891, and TH-0924 held operator-only recommendations that the operator resolved without confirming. The code path holds. A confirmation sets the thread `answered` (`flai/internal/threads/threads.go` `confirm`). The orchestrator took part in the thread, so the thread now awaits it (`flai/internal/inbox/inbox.go` `Summarize`), and the role file's trigger `thread` / `awaits-orchestrator` gives it a `thread:TH-nnnn` reason (`flai/internal/serve/orchestrate_reasons.go`). `Repo.OrchestratorSetsCostOfDelay` admits the input on any open story while the permission is on (`flai/internal/workitem/promotable.go`). No code changes, then. The story pins the path with two regression tests that reproduce the instance and closes the issue on the manifest change. No test covers the confirmation that reaches the orchestrator, nor an in-progress story's input: the permission tests use a ready story and an epic.

**Touches.** The story declared none, and `flai touches suggest` needs paths to start from. It was given the three files below. Its co-change list names the permission's implementation: `edit.go`, `promotable.go`, `items_write.go`, the guard, and the role files. This story changes none of them, so none is kept.

- `flai/internal/serve/orchestrate_confirm_test.go`: layout. A new test file beside `orchestrate_reasons_test.go` for the confirmation-to-reason test, kept apart from the reason queue's tests.
- `flai/cmd/orchestrate_permissions_test.go`: layout and design. The orchestrator permission suite (`testTheOrchestratorSetsCostOfDelayInputsOnlyWithSetCostOfDelay`), where the in-progress case goes.
- `design/issues/I-0284-orchestration-permissions-set-cost-of-delay-is-off-so-a-cost-of-delay-input-the.md`: from criterion 2. `flai issue close` writes it.
- `design/issues/summary.md`: from criterion 2. `flai issue close` regenerates it.

No folder touch.

**Forecast.** flai forecasts 3m from size 2 (2 criteria, 0 touches) over 22 small improvement stories. I raised it to 45m. Each test builds a fixture with a planner's thread, an orchestrator recommendation, and a confirmation, and the serve test drives the orchestrator's look. Delivery is 2026-11-10T12:00:00Z, because review is full (6 of 5) and holds every pull.

**Cost of delay.** The value is 11,250 USD a week, as `flai cod` works it out from the input flai carried over from I-0284: `time_lost_per_cycle` 75h at 150 USD an hour, one cycle a week. I left it unadjusted, since the input is not mine to change. That 75h is W44's two waits, measured before the permission was on. The plan's thread raises this.

**After.** None. The story reads no other story's output.
