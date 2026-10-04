---
id: T-0824
type: task
nature: improvement
title: flai serve queues the planner on an edit, replans when work ahead completes or the order changes, and runs the schedule, coalescing runs
status: done
parent: S-0211
owner: alex
created: 2026-10-04T20:42:10Z
updated: 2026-10-04T21:10:19Z
transitions:
  - to: ready
    at: 2026-10-04T20:42:48Z
    by: agent-S-0211
  - to: in-progress
    at: 2026-10-04T20:55:56Z
    by: agent-S-0211
  - to: done
    at: 2026-10-04T21:10:19Z
    by: agent-S-0211
stream: S-0211
tags: [flai]
touches: [flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, flai/internal/serve/serve.go, flai/internal/serve/plan.go, flai/internal/workitem/changes.go, flai/internal/workitem/changes_test.go, flai/internal/mcpserver/cursor.go]
after: [T-0820, T-0821, T-0822, T-0823]
usage:
  source: log
  seconds: 863
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 126
      output: 777
      cache_read: 10907353
      cache_write: 215442
      cost: 4.4792
---
# T-0824 flai serve queues the planner on an edit, replans when work ahead completes or the order changes, and runs the schedule, coalescing runs

## Work

A replanner per served project in `flai/internal/serve/replan.go`, wired into `Run` beside the launcher's look (on kanban changes and the minute tick), with its decisions in functions that take items, notices, changes, and times, so they are tested without a harness (ADR-0084):

- Nothing while the `plan` host action is off.
- Edit: `itemedit.Notices` since the replanner's watermark, not by a `planner-` agent, naming goal, criteria, or touches of a backlog or ready story whose forecast or value is older; or `cost_of_delay` when the inputs changed after the value.
- Accept, cancel, reorder: `workitem.Changes` to done or cancelled of a story, or a changed pull order. Under `planning.replan` `deterministic` or `agent`, replay every ready and backlog story with a forecast (T-0821) and write the delivery and basis where the delivery moved, by `flai`, in one commit; under `agent`, also queue those stories.
- Schedule: when `planning.schedule`'s next time since the last look has passed, queue every ready story.
- Queue: a story once, triggers merged; one queued run at a time per project, started with its trigger (T-0822) through the checks `serve.Plan` makes; an entry they refuse is dropped and logged.

It waits for T-0820, T-0821, T-0822, and T-0823, whose parser, replay, trigger, and manifest keys it uses.

## Done when

- [ ] Tests cover each trigger, the policy, the gate, and the coalescing
- [ ] `go test -race ./internal/serve/...` passes

## Notes
