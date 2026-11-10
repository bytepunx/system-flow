---
id: S-0653
type: story
nature: remediation
title: TestAFormulaRoleStepsEndIsToldToTheInstancesAgent fails on its timing under load, in the go-test tier of any story that changes a package serve imports
status: backlog
owner: alex
created: 2026-11-10T05:45:21Z
updated: 2026-11-10T07:02:19Z
transitions: []
tags: []
topics: [testing]
touches: [flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/serve/role.go, flai/internal/serve/plan.go, flai/internal/serve/curate.go, flai/internal/serve/orchestrate.go, flai/internal/serve/role_test.go, design/system/strategic-agents.md, design/issues/I-0314-testaformularolestepsendistoldtotheinstancesagent-fails-on-its-timing-under.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-11-10T05:45:21Z
  value: 30
  by: planner-S-0653
  at: 2026-11-10T07:01:00Z
forecast:
  duration: 1h
  delivery: 2026-11-10T09:37:00Z
  basis: "Its own forecast of 1h; 8th in the pull order with an in-progress limit of 5, behind S-0628, S-0636, S-0638, S-0647, S-0656, S-0241, S-0637, S-0648, S-0649, S-0650, S-0651 and S-0652."
  by: flai
  at: 2026-11-10T07:02:19Z
---
# S-0653 TestAFormulaRoleStepsEndIsToldToTheInstancesAgent fails on its timing under load, in the go-test tier of any story that changes a package serve imports

## Goal

This story remediates [I-0314](../../../design/issues/I-0314-testaformularolestepsendistoldtotheinstancesagent-fails-on-its-timing-under.md), "TestAFormulaRoleStepsEndIsToldToTheInstancesAgent fails on its timing under load, in the go-test tier of any story that changes a package serve imports". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0314 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0314 is closed with `flai issue close I-0314 --reason` saying what fixed it

## Tasks
- T-2849 A strategic run's end is measured from its own log, not its kind's newest, with a test that ends an older run after a newer one started
- T-2851 TestAFormulaRoleStepsEndIsToldToTheInstancesAgent waits for each run's end in the inbox instead of reading it once
- T-2852 Close I-0314 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0314. time_lost_per_cycle 12m: 6m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-08T02:50:07Z, 2.1 days before this story; under one cycle counts as one).

### Planning

Cause, from the code. The test is in `flai/internal/serve/role_test.go`. It starts three runs of the role `reviewer` one after another, each once the one before has ended, and the first run alone has usage in its log. Two failures follow from what `roleEnded` in `flai/internal/serve/role.go` does after it records the run as ended.

1. The activity. `roleEnded` logs the activity last, with `logRunEnd` in `flai/internal/serve/activity.go`. `logRunEnd` measures the kind's newest run log (`activityMeasure.newest`), not the log of the run that ended. Under load, the first run's end is logged only after the test has started the second run. The second run's log has no usage, so nothing is logged, and `waitFor` gives up after 30 seconds: "never happened: the worked run's activity logged", as in S-0607's instance. This is a product defect, not only the test's: the design says the activity is measured from the run's log (`design/system/strategic-agents.md` § The planner › The run and how it ended). `plan.go` makes the same call, and several planners of one project run at once, so a planner's end that is logged late can be measured from another planner's log.
2. The inbox. The test reads each run's end from the inbox once, right after `Ended` is set. `roleEnded` appends the agent event after it sets `Ended`, so a loaded host can leave the event unwritten when the test reads it, and the test fails with "agent-S-nnnn is told of []". The message changes with which run lost the race, as in S-0633's instance.

Proposed fix: `logRunEnd` measures the log of the run that ended, `AgentRun.Log`, which each caller passes. The test polls the inbox for each end with `waitFor`, as it does for the activity. A test in `activity_test.go` that ends an older run after a newer run of the same kind has started reproduces the first cause without load.

Touches:

| Path | Source | Why |
|------|--------|-----|
| `flai/internal/serve/activity.go` | design | `logRunEnd` and `activityMeasure.newest` pick the log measured |
| `flai/internal/serve/activity_test.go` | co-change | `flai touches suggest`: 8 of 24 commits; the reproducing test goes here |
| `flai/internal/serve/role.go` | design | `roleEnded` calls `logRunEnd` |
| `flai/internal/serve/plan.go` | co-change | 6 of 24 commits; the planner's end calls `logRunEnd` the same way |
| `flai/internal/serve/curate.go` | co-change | 3 of 24 commits; calls `logRunEnd`, so a changed signature reaches it |
| `flai/internal/serve/orchestrate.go` | co-change | 3 of 24 commits; calls `logRunEnd`, as `curate.go` does |
| `flai/internal/serve/role_test.go` | design | the failing test, named in both instances |
| `design/system/strategic-agents.md` | design | says which log a run's end is measured from |
| `design/issues/I-0314-…md` | design | criterion 2 closes it |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

Left out from `flai touches suggest`: `agents.go`, `serve.go`, `analyze.go`, and the `flai/cmd` files. None calls `logRunEnd`, and the analyzer's end is logged by `logRunEndSaying` with one run per project at a time. The task that changes `activity.go` adds them with `flai touches` if the change reaches them.

Forecast: `flai forecast` gave 5m from the median of 17 small remediation stories at size 2, counted with no touches. It is raised to 1h. The fix changes how serve measures a run's activity for four callers and adds a test that reproduces the cause. Each go-test run of serve takes 3 to 7 minutes under load, as I-0314 says, and the close-out runs the integration tier. The delivery moves 55m past `flai forecast`'s 08:10 for the same reason.

Cost of delay: `flai cod` gives 30.00 USD a week from the 12m time_lost_per_cycle flai set from I-0314, at 150 USD an hour. It is written as it stands.

No `after`: the story reads no other story's output.
