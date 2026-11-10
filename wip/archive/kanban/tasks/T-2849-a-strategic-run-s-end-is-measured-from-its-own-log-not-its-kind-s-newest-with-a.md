---
id: T-2849
type: task
nature: remediation
title: A strategic run's end is measured from its own log, not its kind's newest, with a test that ends an older run after a newer one started
status: done
parent: S-0653
owner: alex
created: 2026-11-10T07:01:26Z
updated: 2026-11-10T07:34:42Z
transitions:
  - to: ready
    at: 2026-11-10T07:34:24Z
    by: agent-S-0653
  - to: in-progress
    at: 2026-11-10T07:34:26Z
    by: agent-S-0653
  - to: done
    at: 2026-11-10T07:34:42Z
    by: agent-S-0653
stream: S-0653
tags: [flai, serve]
touches: [flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/serve/role.go, flai/internal/serve/plan.go, flai/internal/serve/curate.go, flai/internal/serve/orchestrate.go, design/system/strategic-agents.md, flai/internal/serve/analyze.go]
usage:
  source: log
  seconds: 16
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 617
      cache_read: 562839
      cache_write: 7334
      cost: 0.1836
---
# T-2849 A strategic run's end is measured from its own log, not its kind's newest, with a test that ends an older run after a newer one started

## Work

First layer: it waits for no task, and shares no path with the test task beside it.

- Reproduce first. In `flai/internal/serve/activity_test.go`, write two run logs of one kind, the older with usage and the newer with none, the newer starting after the older ended. Log the older run's end with `logRunEnd`. Today nothing is logged, because `activityMeasure.newest` picks the newer log. Record that failing run in the task's Notes.
- Fix. Have `logRunEnd` measure the log of the run that ended. Each caller has it as `AgentRun.Log`: `roleEnded` in `role.go`, the planner's end in `plan.go`, and the ends in `curate.go` and `orchestrate.go`. Fall back to the newest log only where the caller has no log, such as an activity logged outside a run. Check that the span still starts at the later of the run's start and the document's last entry, so two runs of one kind that overlap are not charged twice.
- Keep `logRunEndSaying`, the analyzer's, in step with it.
- Say in `design/system/strategic-agents.md`, under § The planner › The run and how it ended and the defined role's run, that the activity is measured from the ended run's own log.
- If the change reaches a file outside this task's touches, such as `agents.go` or `analyze.go`, add it with `flai touches` and say why in the narrative's Decisions.

## Done when

- Before the fix, the new test fails with no activity logged; the run is in the task's Notes.
- After the fix, the new test passes, and `go test ./internal/serve -run 'Activity|RunEnd|Role|Plan' -count 5` passes in `flai/`.
- `flai test` on the task's Go files passes.

## Notes
