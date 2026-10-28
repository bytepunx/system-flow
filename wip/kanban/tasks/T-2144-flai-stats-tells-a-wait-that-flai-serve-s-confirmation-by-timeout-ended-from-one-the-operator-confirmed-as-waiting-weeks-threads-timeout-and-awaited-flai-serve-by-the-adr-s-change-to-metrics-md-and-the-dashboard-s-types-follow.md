---
id: T-2144
type: task
nature: improvement
title: "flai stats tells a wait that flai serve's confirmation by timeout ended from one the operator confirmed, as waiting.weeks[].threads.timeout and awaited flai serve, by the ADR's change to metrics.md, and the dashboard's types follow"
status: done
parent: S-0502
owner: alex
created: 2026-10-26T04:45:09Z
updated: 2026-10-28T13:23:19Z
transitions:
  - to: ready
    at: 2026-10-28T13:05:46Z
    by: agent-S-0502
  - to: in-progress
    at: 2026-10-28T13:05:48Z
    by: agent-S-0502
  - to: done
    at: 2026-10-28T13:23:19Z
    by: agent-S-0502
stream: S-0502
tags: [flai, flaiover]
touches: [flai/internal/metrics/waiting.go, flai/internal/metrics/waiting_test.go, design/system/metrics.md, flaiover/src/lib/viz/charts.ts]
after: [T-2143]
usage:
  source: log
  seconds: 1051
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 35
      output: 11773
      cache_read: 2871583
      cache_write: 73073
      cost: 1.2575
---
# T-2144 flai stats tells a wait that flai serve's confirmation by timeout ended from one the operator confirmed, as waiting.weeks[].threads.timeout and awaited flai serve, by the ADR's change to metrics.md, and the dashboard's types follow

## Work

Point 6 of T-2141's ADR. `design/system/metrics.md` is the contract between `flai stats` and the dashboard and changes only with an ADR, which T-2141 is.

- `flai/internal/metrics/waiting.go`: a third ender beside `byOrchestrator` and `byConfirmation`, `byTimeout`, for an entry that confirms a recommendation (`Confirms()`) and is by `flai serve`; `byConfirmation` then covers only the rest, the operator's. `ThreadWaitTotal` gains `Timeout` with `count` and `total_seconds`, computed as `Confirmed` is in `waiting()`, a union of its own waits within the item's `in-progress` intervals, present in every week and 0 without such waits. `waiting.longest[].awaited` already reads the ending entry's author, so it reads `flai serve`; `items[].wait_threads_seconds` is unchanged. Tests in `waiting_test.go`: a thread the operator confirmed counts under `confirmed`, one `flai serve` confirmed by timeout counts under `timeout` and not under `confirmed`, and the week's totals and `longest` agree.
- `design/system/metrics.md` § Waiting: the row `Who ended a wait` gains the timeout, and the paragraph on `waiting.weeks[].threads` names `timeout` beside `orchestrator` and `confirmed`, with `confirmed` now the operator's alone; the `updated` date moves.
- `flaiover/src/lib/viz/charts.ts`: `WaitWeek.threads` gains `timeout: WaitShare`, so the dashboard's type matches `flai stats`; the agent-waiting chart draws the totals and is unchanged. `flai cmd/stats.go`'s text output prints no split of the enders today and is unchanged.

## Done when

- `flai stats --json` on a fixture with both confirmations reports `threads.confirmed.count` 1 and `threads.timeout.count` 1 for the week, and `longest[]` names `flai serve` as `awaited` for the timeout's wait.
- metrics.md states the field, and `flai test` on the changed paths passes, the flaiover lint and vitest for `charts.ts` among them.

## Notes

Layer 3, waits for T-2143 (the confirming entry's author and line). The docs task waits for it. It shares no path with the serve look or the prompt task.
