---
id: S-0599
type: story
nature: improvement
title: flai serve's replanner writes every forecast it moves in one pass with one check, so an acceptance's replan lands in seconds, not 80 s
status: review
parent: E-0033
owner: alex
created: 2026-11-07T22:39:17Z
updated: 2026-11-08T02:06:30Z
transitions:
  - to: ready
    at: 2026-11-07T23:49:24Z
    by: alex
  - to: in-progress
    at: 2026-11-07T23:49:39Z
    by: agent-S-0599
  - to: review
    at: 2026-11-08T02:06:30Z
    by: agent-S-0599
tags: [cli]
touches: [flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, flai/internal/itemedit/forecasts.go, flai/internal/itemedit/forecasts_test.go, flai/internal/workitem/itemlock.go, flai/internal/workitem/itemlock_test.go, flai/internal/check/items_test.go, design/system/server-performance.md, design/system/strategic-agents.md, flai/internal/check/items.go, flai/internal/workitem/repo.go, design/issues/I-0311-teststartdevserverstopsvitewhenitscontextends-fails-under-load-the-fake-vite.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 8245
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 1
      hand_edits: 1
      work: 18
    - day: 2026-11-08
      ceremony: 6
      hand_edits: 1
      work: 178
  models:
    - model: claude-opus-5-5
      input: 568
      output: 99004
      cache_read: 103766644
      cache_write: 1683866
      cost: 32.5425
  strategic:
    - kind: orchestrator
      seconds: 232
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 25
          output: 204
          cache_read: 1414233
          cache_write: 98335
          cost: 0.387
        - model: claude-sonnet-5-5
          input: 4
          output: 41
          cache_read: 12661
          cache_write: 30294
          cost: 0.0314
forecast:
  duration: 2h30m
  delivery: 2026-11-08T02:30:00Z
  basis: "Five tasks in four layers across workitem, check, itemedit, and serve, a measurement in a scratch copy, and a close-out whose integration tier means about 4 min; flai forecast's 6m counts only 3 criteria and 1 touch."
  by: planner-S-0599
  at: 2026-11-07T23:47:15Z
finalized:
  by: alex
  at: 2026-11-07T23:49:05Z
---
# S-0599 flai serve's replanner writes every forecast it moves in one pass with one check, so an acceptance's replan lands in seconds, not 80 s

## Goal

Cause 9 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): the replanner writes each moved forecast through `itemedit.Apply`, one after another (`flai/internal/serve/replan.go`), so each pays a checked edit's two whole-repository checks, 5 to 6.4 s a forecast. Over the week to 2026-11-07, 98 replans moved 1,200 forecasts, about 1.7 to 2 h of serve CPU, and an acceptance's replan commit landed p50 80 s and p90 127 s after the acceptance.

## Acceptance criteria

- [x] A replan writes every forecast it moves in one pass under the items lock, with one check over the written items and one write of `wip/agents/index.md`.
- [x] A replan that moves 20 forecasts in a scratch copy of this repository takes under 5 s, and server-performance.md records the before and after.
- [x] A test shows a replan writes the same forecasts as before, and that a forecast another writer changed meanwhile is merged or refused as ADR-0184 says.

## Tasks
- T-2628 Save many work items under one hold of the items lock
- T-2632 Check the work-item rules and the markdown lint over a given set of items
- T-2641 Apply many forecast changes in one pass with one check and one index write
- T-2647 The replanner writes its moved forecasts through the batch entry
- T-2655 Measure a 20-forecast replan in a scratch copy and record cause 9's before and after

## Notes

- Proposed fix: an `itemedit` entry that applies many forecast changes at once. Forecasts are values flai owns, so the check can be the work-item rules over the written items.
- Independent of the story for cause 8, but it gains from it: with a scoped checked edit alone, 20 forecasts would still cost about 40 s.

### Planning

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/serve/replan.go` | declared | `replan` calls `itemedit.Apply` once per moved forecast |
| `flai/internal/serve/replan_test.go` | co-change | 4 of the 7 commits that changed `replan.go` changed it too; criterion 3's test |
| `flai/internal/itemedit/forecasts.go` | layout | New file: the batch forecast entry the Notes propose, beside `Apply` |
| `flai/internal/itemedit/forecasts_test.go` | layout | Its tests: same forecasts as `Apply`, merge and refusal under ADR-0184 |
| `flai/internal/workitem/itemlock.go` | design | ADR-0184's items lock; `Save` takes it per file, so one pass needs a save of many items under one hold |
| `flai/internal/workitem/itemlock_test.go` | layout | Tests of that save |
| `flai/internal/check/check.go` | design | The Notes' "work-item rules over the written items": `oneItem` and the lint, exported for a set of items |
| `flai/internal/check/items_test.go` | layout | New file: tests of that entry |
| `design/system/server-performance.md` | design | Criterion 2: the before and after of cause 9 |
| `design/system/strategic-agents.md` | design | § The replan policy says the replay writes through `itemedit.Apply` and names its per-story cost |

Left out: the other co-change files (`analysis_schedule.go`, `curate_schedule.go`, `plan.go`, `role_schedule.go`, `serve.go`), which changed with `replan.go` for the schedules, not for the write path. No folder touch.

No `after`. S-0598 (cause 8) scopes the checked edit, which this story does not call: the batch runs its own item-level check. Both may change `check.go` and `itemedit`; the trial merge judges that.

Forecast: 2h30m, not the 6m `flai forecast` gives. Its 79 s a unit of size times size 4 counts 3 criteria and 1 touch. The plan has 5 tasks over 4 packages, a measurement in a scratch copy, and a close-out whose `integration` tier means about 238 s (server-performance.md, cause 20). Delivery 2026-11-08T02:30Z assumes the draft is finalized and pulled tonight; the replanner moves it otherwise.

Cost of delay: no value yet. Neither S-0599 nor E-0033 has inputs, so `flai cod` refuses. The input is asked on TH-0888 (recommended `penalty_per_week: 60`; the orchestrator suggests 100 may fit) and on this story's plan thread.
