---
id: S-0851
type: story
nature: improvement
title: flai serve still reads three state questions by the baseline workflow directly, not through a wrapper, so a project's own workflow file will not reach them
status: backlog
owner: alex
created: 2026-11-27T22:34:11Z
updated: 2026-11-28T00:43:35Z
transitions: []
tags: [cli]
topics: [workflow]
touches: [flai/internal/serve/curate_schedule.go, flai/internal/serve/curate_schedule_test.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/components.go, flai/internal/serve/components_test.go, design/system/component-library.md, design/issues/I-0406-flai-serve-still-reads-three-state-questions-by-the-baseline-workflow-directly.md, design/issues/I-0402-flai-serve-stops-the-dev-server-of-a-story-in-a-state-a-project-s-workflow-adds.md, design/issues/summary.md]
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
      seconds: 87
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 39
          output: 6365
          cache_read: 3149186
          cache_write: 72919
          cost: 1.3464
    - kind: orchestrator
      seconds: 8
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 8
          cache_read: 75694
          cache_write: 406
          cost: 0.0243
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 7m
    by: flai
    at: 2026-11-27T22:34:11Z
  value: 17.5
  by: planner-S-0851
  at: 2026-11-28T00:03:14Z
forecast:
  duration: 1h
  delivery: 2026-11-28T08:29:00Z
  basis: "Its own forecast of 1h; 12th in the pull order with an in-progress limit of 3, behind S-0814, S-0835, S-0837, S-0852, S-0813, S-0822, S-0848, S-0853, S-0854, S-0824, S-0845, S-0241, S-0815 and S-0849."
  by: flai
  at: 2026-11-28T00:43:35Z
---
# S-0851 flai serve still reads three state questions by the baseline workflow directly, not through a wrapper, so a project's own workflow file will not reach them

## Goal

This story remediates [I-0406](../../../design/issues/I-0406-flai-serve-still-reads-three-state-questions-by-the-baseline-workflow-directly.md), "flai serve still reads three state questions by the baseline workflow directly, not through a wrapper, so a project's own workflow file will not reach them". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0406 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0406 is closed with `flai issue close I-0406 --reason` saying what fixed it

## Tasks
- T-3652 The curate scheduler matches an acceptance by the project's workflow, not the baseline's
- T-3653 StoryLeft reads the kinds of the story's project's workflow, so a dev server stays up in a state the workflow adds
- T-3654 Close I-0406 and I-0402, naming the two reads fixed and why kluster.go's legacy send-back check keeps the baseline

## Notes

Cost of delay inputs set by flai from I-0406. time_lost_per_cycle 7m: 10m per occurrence × 1 occurrence ÷ 1.4 cycles of 168h (first reported 2026-11-18T10:51:47Z, 9.5 days before this story).

### Planning

The proposed solution, from I-0406's three instances, read against main:

| Read | Where | Repository at hand | Proposal |
|------|-------|--------------------|----------|
| `acceptance()` | `flai/internal/serve/curate_schedule.go` | yes, the scheduler's look | read `repo.Flow()` (T-3652); S-0748, which its comment waits for, is done |
| `StoryLeft` | `flai/internal/serve/components.go` | yes, `readStory` opens the project | read the story's project's workflow (T-3653); this is also defect I-0402 |
| `klusterKept` | `flai/internal/serve/kluster.go` | yes | keep the baseline: it judges only a legacy send-back record without `To`, written by a flai that ran only the baseline; T-3654 checks this before closing |

`builtinTriggersIn` in `orchestrate_reasons.go` also reads `baselineFlow` (D-0676). It translates the baseline's state names in `builtinTriggers` to the project's states, so it is deliberate and outside this story.

Touches, and where each came from:

- `flai/internal/serve/curate_schedule.go`, `flai/internal/serve/components.go`: design, named by I-0406.
- `flai/internal/serve/curate_schedule_test.go`: co-change (`flai touches suggest`, 3 of 25 commits), and it calls `accepted()`, which T-3652 removes.
- `flai/internal/serve/components_test.go`: layout, the test file of `StoryLeft` (`TestStoryLeft`).
- `flai/internal/serve/orchestrate_reasons_test.go`: layout, it holds the test of `acceptance()`.
- `design/system/component-library.md`: design, § Leaving review describes the stop `StoryLeft` decides.
- `design/issues/I-0406-…md`, `design/issues/I-0402-…md`, `design/issues/summary.md`: design, the issues the story closes and their generated summary.
- `flai/internal/serve/kluster.go` and `kluster_test.go` are left out: the proposal changes neither. `flai/cmd/serve_actions.go` is left out: T-3653 picks the form that leaves it unchanged. The `flai touches suggest` paths under `flai/internal/kluster/` and `docs/` come from kluster.go's other changes, not this story's.

No folder touch. No `after`: the story reads nothing another open story writes.

Figures:

- Forecast: `flai forecast` gave 23m (125 s per unit of size, size 11, medium band). Raised to 1h. The story has two Go changes in `flai/internal/serve`, each with a new test on a renamed workflow, and its close-out runs the integration and smoke tiers that `flai/` selects. Delivery moved 37m later to match.
- Cost of delay: `flai cod` gave 17.50 USD a week from the operator's 7m time lost per 168h cycle at 150 USD an hour. It stands: the occurrence count is one. Closing I-0402, which T-3653 also fixes, would add its 10m, but its cost is not among this story's inputs.
