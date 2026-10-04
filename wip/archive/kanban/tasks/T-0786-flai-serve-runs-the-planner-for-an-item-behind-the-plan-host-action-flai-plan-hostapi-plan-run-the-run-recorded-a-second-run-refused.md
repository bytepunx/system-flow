---
id: T-0786
type: task
nature: feature
title: "flai serve runs the planner for an item behind the plan host action: flai plan, hostapi plan.run, the run recorded, a second run refused"
status: done
parent: S-0208
owner: alex
created: 2026-10-04T00:43:57Z
updated: 2026-10-04T03:14:29Z
transitions:
  - to: ready
    at: 2026-10-04T00:44:32Z
    by: agent-S-0208
  - to: in-progress
    at: 2026-10-04T00:52:31Z
    by: agent-S-0208
  - to: done
    at: 2026-10-04T03:14:29Z
    by: agent-S-0208
stream: S-0208
tags: []
touches: [flai/internal/serve, flai/internal/hostapi, flai/cmd]
after: [T-0784]
usage:
  source: log
  seconds: 812
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 223
      output: 90189
      cache_read: 13989103
      cache_write: 351565
      cost: 6.655
---
# T-0786 flai serve runs the planner for an item behind the plan host action: flai plan, hostapi plan.run, the run recorded, a second run refused

## Work

- A `plan` host action, off by default, in `hostapi.Actions`, seen by the dashboard, enabled and disabled with `flai serve enable|disable plan`.
- `flai serve` starts a planner run for an epic or a story: in the main checkout, on the planning agent (T-0784), with the harness's planner request; `FLAI_AGENT` names the planner and its item (`planner-E-nnnn`), `FLAI_ROLE=plan`, `FLAI_ITEM`, `FLAI_SESSION`, `FLAI_STARTED_BY=flai-serve`, and no `FLAI_STORY`.
- The run is recorded as a story's run is: its log as `<key>-planner-<start>.log` beside the story runs' logs (what `Dir.ActivityLogs` reads, ADR-0079), its session, its exit, and its outcome (`asked` when it left a thread on the item awaiting the operator, `failed` on a failed exit, `stopped` when stopped, `worked` otherwise), kept in the serve state per project by item. When it ends, `serve.LogRunEnd` logs its activity in `wip/agents/planner.md`.
- A second run on the same item while one runs is refused, naming the running one; a task, an archived item, or an item that is neither an epic nor a story is refused.
- `flai plan <E-nnnn|S-nnnn>` on the host asks the running flai serve to start it, as `flai serve agent start` does, and the hostapi method `plan.run` with `{id}` does the same for a dashboard, gated on the `plan` action (refused with `Disabled` while it is off) and journalled.

Waits for T-0784: it starts the run with the planning agent and the harness's planner request.

## Done when

- [x] Tests cover the start (log name, environment, recorded run, activity logged on end), the refusal of a second run on the same item, and `plan.run` refused while `plan` is off and started while on
- [x] `flai plan` works against a running flai serve, and says how to enable the action when it is off

## Notes
