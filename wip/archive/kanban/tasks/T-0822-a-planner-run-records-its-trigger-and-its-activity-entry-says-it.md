---
id: T-0822
type: task
nature: improvement
title: A planner run records its trigger, and its activity entry says it
status: done
parent: S-0211
owner: alex
created: 2026-10-04T20:41:54Z
updated: 2026-10-04T20:49:32Z
transitions:
  - to: ready
    at: 2026-10-04T20:42:48Z
    by: agent-S-0211
  - to: in-progress
    at: 2026-10-04T20:42:49Z
    by: agent-S-0211
  - to: done
    at: 2026-10-04T20:49:32Z
    by: agent-S-0211
stream: S-0211
tags: [flai]
touches: [flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/agents.go, design/system/agent-narrative.md]
usage:
  source: log
  seconds: 403
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 52
      output: 279
      cache_read: 2011864
      cache_write: 86191
      cost: 0.845
---
# T-0822 A planner run records its trigger, and its activity entry says it

## Work

`ActivityEntry` gains `Trigger`, written as an optional `- Trigger:` line after the summary and read back; an entry without one stays valid (ADR-0084). `AgentRun` gains `trigger` in `serve/agents.json`. `serve.Plan` keeps its signature and records `asked`; an internal start takes a trigger, for the queue to come. `planEnded` passes the run's trigger to `LogRunEnd`, which writes it on the entry. `design/system/agent-narrative.md` shows the line in its example. It waits for nothing: the first layer; the queue task waits for it.

## Done when

- [ ] An activity document round-trips with and without the Trigger line
- [ ] A planner run the operator asked for is logged with `- Trigger: asked`
- [ ] Tests in `flai/internal/workitem` and `flai/internal/serve` for these pass

## Notes
