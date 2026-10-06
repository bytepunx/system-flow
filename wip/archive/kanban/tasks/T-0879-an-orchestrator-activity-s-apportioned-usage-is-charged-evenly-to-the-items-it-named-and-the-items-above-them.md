---
id: T-0879
type: task
nature: improvement
title: An orchestrator activity's apportioned usage is charged evenly to the items it named and the items above them
status: done
parent: S-0226
owner: alex
created: 2026-10-05T04:44:41Z
updated: 2026-10-06T18:11:55Z
transitions:
  - to: ready
    at: 2026-10-06T18:02:21Z
    by: agent-S-0226
  - to: in-progress
    at: 2026-10-06T18:02:22Z
    by: agent-S-0226
  - to: done
    at: 2026-10-06T18:11:55Z
    by: agent-S-0226
stream: S-0226
tags: [flai]
touches: [flai/internal/usage, flai/internal/serve, flai/internal/workitem, flai/internal/mcpserver, flai/cmd/activity.go, design/system/flai-cli.md]
after: [T-0878]
usage:
  source: log
  seconds: 573
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 88
      output: 32100
      cache_read: 4341168
      cache_write: 142533
      cost: 2.4361
---
# T-0879 An orchestrator activity's apportioned usage is charged evenly to the items it named and the items above them

## Work

Build the charge T-0878's ADR decides. It waits for T-0878 so that the split, the edge cases, and which kinds charge are settled before the code; it shares no path with T-0880 and runs beside it.

- `flai/internal/usage`: an exported way to split a `Usage` into n even shares whose tokens, cost, and seconds add up to the whole, built on `Model.scaled`, with tests.
- `flai/internal/serve/activity.go`: `charge` today returns at once for any kind but the planner. For the orchestrator, take the items the logged entry names, drop duplicates and IDs that do not exist, split the apportioned usage between the rest with the entry's seconds, and call `ChargeStrategic` once per item, so each share is summed up to its epic. Both paths reach it: `LogActivity` (`activity_log`) and `LogRunEnd` (the run's end).
- `Logged` reports every item charged, not one `Planned`; `flai activity log` (`flai/cmd/activity.go`) and the MCP `activity_log` result and description (`flai/internal/mcpserver`) say what was charged to which items, and that an activity naming none goes to the project total.
- Change `flai/internal/workitem` only if `ChargeStrategic` needs it.
- Tests in `flai/internal/serve` pin: one item gets the whole; two items get half each; a story and its epic named together give the epic both shares; a duplicate and an unknown ID are left out; an activity naming none charges no item; the planner's charge is unchanged.
- Update `design/system/flai-cli.md` where it describes `activity_log` and `flai activity log`.

## Done when

- An orchestrator activity logged through `activity_log` or at its run's end leaves its apportioned usage under the `orchestrator` entry of each item it named, split evenly, and summed up to each one's epic.
- The tests above pass with `scripts/flai-test.sh`, and the planner's existing charge tests still pass.
- `design/system/flai-cli.md` says what the orchestrator's activity charges.

## Notes

- Drafted by planner-S-0226.
