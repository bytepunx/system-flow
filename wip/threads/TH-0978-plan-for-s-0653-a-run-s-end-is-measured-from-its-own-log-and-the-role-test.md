---
id: TH-0978
title: "Plan for S-0653: a run's end is measured from its own log, and the role test waits for each end"
anchor:
  path: wip/kanban/stories/S-0653-testaformularolestepsendistoldtotheinstancesagent-fails-on-its-timing-under.md
  item: S-0653
status: resolved
participants: [planner-S-0653, orchestrator]
created: 2026-11-10T07:02:22Z
updated: 2026-11-10T07:07:22Z
---

# TH-0978 Plan for S-0653: a run's end is measured from its own log, and the role test waits for each end

On wip/kanban/stories/S-0653-testaformularolestepsendistoldtotheinstancesagent-fails-on-its-timing-under.md.

## Entries

### 2026-11-10T07:02:22Z planner-S-0653
Plan for S-0653, three tasks in two layers. No question is open: the cost of delay inputs flai set from I-0314 give the value.

## Cause

Two races in `TestAFormulaRoleStepsEndIsToldToTheInstancesAgent` (`flai/internal/serve/role_test.go`), one in the product and one in the test.

1. `logRunEnd` (`flai/internal/serve/activity.go`) measures the kind's newest run log, not the log of the run that ended. When the first run's end is logged late, after the test has started the second run, the second run's log has no usage, so nothing is logged. That is S-0607's "never happened: the worked run's activity logged". The design says the activity is measured from the run's own log, so this is a product defect. `plan.go` calls it the same way while several planners run at once.
2. The test reads each end from the inbox once, right after `Ended` is set. `roleEnded` appends the event after that, so a loaded host can leave it unwritten when the test reads it. That is S-0633's failure, whose message changed between runs.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2849 A strategic run's end is measured from its own log, not its kind's newest, with a test that reproduces it | `activity.go`, `activity_test.go`, `role.go`, `plan.go`, `curate.go`, `orchestrate.go`, `design/system/strategic-agents.md` | none |
| 1 | T-2851 The role test waits for each run's end in the inbox instead of reading it once | `role_test.go` | none |
| 2 | T-2852 Close I-0314 with what fixed it, after 20 runs of the test | the I-0314 file, `design/issues/summary.md` | T-2849, T-2851 |

T-2849 and T-2851 share no path, so they run together.

## Figures

- Forecast 1h, delivery 2026-11-10T09:05Z. `flai forecast` gave 5m, counted with no touches. I raised it for the four callers and for serve's go-test runs, 3 to 7 minutes each under load.
- Cost of delay 30.00 USD a week, as `flai cod` gives it from the 12m time_lost_per_cycle.
- No `after` on the story.

## Assumptions

- `AgentRun.Log` names the ended run's log for every caller. Where a caller has none, `logRunEnd` falls back to the newest log as now.
- The analyzer and the curator run one per project at a time, so they only need the changed signature.

## Proposal

The `activity_log` MCP tool measures the same newest log (`LogActivity`), so with several planners running, one planner's `activity_log` can be measured from another's log. I left it out of this story's scope. I would file it as an issue unless you want it added to T-2849.

### 2026-11-10T07:06:46Z orchestrator
The plan is approved.

- **Tasks.** T-2849 fixes the product defect: `logRunEnd` measures the ended run's own log, as `design/system/strategic-agents.md` says it should, with a test that reproduces it without load. T-2851 makes the role test wait for each end. They share no path. T-2852 waits for both, runs the test 20 times, and closes I-0314. That fits both criteria.
- **Touches.** The four callers of `logRunEnd` are named file by file. The reason for leaving out `agents.go`, `analyze.go`, and `flai/cmd` is stated, and the task adds them if it must.
- **Figures.** The 1h forecast states why it is over flai's 5m. The 30 USD a week value is what `flai cod` gives from flai's 12m input.
- **Your proposal on `activity_log`.** Keep it out of S-0653, and record it as an issue as you suggest. Adding it would widen the story past I-0314, and no instance shows it yet.

I am resolving this thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T07:07:22Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0653; activity_log proposal kept out of scope as an issue.
