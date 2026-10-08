---
id: S-0348
type: story
nature: improvement
title: "flai check finds `board.wip-limit` outside the story at close-out"
status: in-progress
owner: alex
created: 2026-10-08T08:37:06Z
updated: 2026-10-08T09:18:00Z
transitions:
  - to: ready
    at: 2026-10-08T09:09:03Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-08T09:18:00Z
    by: agent-S-0348
tags: [flai, check]
topics: [cli, conventions, template]
touches: [design/adrs, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0123-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md, design/issues/summary.md]
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
      seconds: 250
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 210
          output: 7700
          cache_read: 1170700
          cache_write: 64849
          cost: 0.2368
        - model: claude-opus-5-5
          input: 87
          output: 34492
          cache_read: 7010398
          cache_write: 159138
          cost: 3.3654
    - kind: orchestrator
      seconds: 360
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 106
          output: 1800
          cache_read: 34819572
          cache_write: 39129
          cost: 8.5887
        - model: claude-sonnet-5-5
          input: 4
          output: 35
          cache_read: 29103
          cache_write: 17518
          cost: 0.0402
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-10-08T08:58:23Z
  value: 12.5
  by: planner-S-0348
  at: 2026-10-08T08:59:43Z
forecast:
  duration: 21m
  delivery: 2026-10-08T09:33:00Z
  basis: "Its own forecast of 21m; 1st in the pull order with an in-progress limit of 5, behind S-0232, S-0322, S-0341, S-0344 and S-0345."
  by: flai
  at: 2026-10-08T09:09:50Z
finalized:
  by: orchestrator
  at: 2026-10-08T09:00:47Z
---
# S-0348 flai check finds `board.wip-limit` outside the story at close-out

## Goal

This story remediates [I-0123](../../../design/issues/I-0123-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md), "flai check finds `board.wip-limit` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0123 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0123 is closed with `flai issue close I-0123 --reason` saying what fixed it

## Tasks
- T-1422 Record in an ADR that a close-out records no board.wip-limit and a check scoped to a story leaves it out
- T-1423 Leave board.wip-limit out of a check scoped to a story, with a test that reproduces I-0123
- T-1424 Say in the design, the user guide, and the work-management convention that a close-out leaves out board.wip-limit
- T-1425 Close I-0123 saying what fixed it

## Notes

### Planning

Proposed fix, from I-0123's one instance (S-0339's close-out, `wip/kanban/board.md: 6 stories in in-progress, limit 5`): a check scoped to a story leaves `board.wip-limit` out, as it leaves out `item.archive` (ADR-0122) and the findings on another open story's narrative (ADR-0123, ADR-0125). The board's counts are the main checkout's state, which no story branch changes, and the breach names no story. `flai check --strict` in the main checkout still fails on it (ADR-0073 unchanged). The finding comes from `checker.board` in `flai/internal/check/check.go`. The scoped check is `ScopeToStory` in `flai/internal/check/scope.go`. The issue is recorded by `recordOutside` in `flai/cmd/check.go`.

Tasks, in layers:

1. T-1422, the ADR.
2. T-1423, the code and its tests, and T-1424, the design, user guide, and convention. Both wait for T-1422 for its number, and they share no path.
3. T-1425 closes I-0123, after T-1423 and T-1424.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/check/scope.go`, `flai/cmd/check.go` | layout: where the finding is scoped and recorded |
| `flai/internal/check/scope_test.go`, `flai/cmd/check_test.go`, `docs/users/flai-reference.md`, `docs/users/flai.md`, `design/system/flai-cli.md` | co-change with the two above (`flai touches suggest`, 29-41%); S-0280, S-0318, and S-0323 changed each |
| `design/system/continuous-improvement.md`, `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md` | design: they list what a close-out leaves out; S-0280, S-0318, and S-0323 changed each |
| `design/adrs` | design: the new ADR, as in each precedent |
| `design/issues/I-0123-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md`, `design/issues/summary.md` | the issue the story closes |

`flai/internal/check/check.go` (35% co-change) is left out: the fix scopes the finding and leaves its production alone.

Folder touch kept: `design/adrs`, because T-1422 adds an ADR whose number is not known until it is written, and other open stories add ADRs too. While this story is in progress it holds S-0334, S-0337, and S-0297, which touch it as well.

Forecast: 21m, delivery 2026-10-08T14:34:00Z, as `flai forecast` gives it (78 s per unit of size over 51 large-band improvement stories, size 16). Kept: the three precedent stories took 12m, 22m, and 49m of agent time.

Cost of delay: 12.50 USD a week, as `flai cod` gives it from `time_lost_per_cycle: 5m`, the input the orchestrator set on TH-0383 as recommended. Kept: it matches S-0346, a close-out finding of the same kind.
