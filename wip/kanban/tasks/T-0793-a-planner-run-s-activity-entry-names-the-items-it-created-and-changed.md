---
id: T-0793
type: task
nature: feature
title: A planner run's activity entry names the items it created and changed
status: done
parent: S-0209
owner: alex
created: 2026-10-04T04:02:45Z
updated: 2026-10-04T04:08:07Z
transitions:
  - to: ready
    at: 2026-10-04T04:03:32Z
    by: agent-S-0209
  - to: in-progress
    at: 2026-10-04T04:03:33Z
    by: agent-S-0209
  - to: done
    at: 2026-10-04T04:08:07Z
    by: agent-S-0209
stream: S-0209
tags: []
touches: [flai/internal/serve]
usage:
  source: log
  seconds: 274
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 47
      output: 222
      cache_read: 1657109
      cache_write: 75109
      cost: 0.7066
---
# T-0793 A planner run's activity entry names the items it created and changed

## Work

- When a planner run ends (`planEnded`, `flai/internal/serve/plan.go`), its activity entry's items are the planned item, then the items under it created since the run started, then those under it changed since then (their `updated`), so the entry names the stories created and revisited beside the run's cost. `LogRunEnd` takes the items.
- Waits for nothing: the first layer.

## Done when

- [ ] A test in `flai/internal/serve` shows a planner run's entry naming the epic, a story created during the run, and a story edited during it, and not one left alone.
- [ ] `go test ./internal/serve` passes.

## Notes

`LogRunEnd` logs no items today, so the entry names nothing unless the planner's summary does.
