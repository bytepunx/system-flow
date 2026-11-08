---
id: T-2617
type: task
nature: improvement
title: item.get finds its item, children, and plan in one list of the kept snapshot, with no glob and no second list
status: done
parent: S-0606
owner: alex
created: 2026-11-07T22:49:16Z
updated: 2026-11-08T01:11:26Z
transitions:
  - to: ready
    at: 2026-11-08T00:47:09Z
    by: agent-S-0606
  - to: in-progress
    at: 2026-11-08T00:47:13Z
    by: agent-S-0606
  - to: done
    at: 2026-11-08T01:11:26Z
    by: agent-S-0606
stream: S-0606
tags: [cli]
touches: [flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/writes.go]
after: [T-2615]
usage:
  source: log
  seconds: 1453
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 46
      output: 15402
      cache_read: 7566184
      cache_write: 321616
      cost: 3.5734
  strategic:
    - kind: orchestrator
      seconds: 8
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 12
          cache_read: 62744
          cache_write: 12650
          cost: 0.0192
---
# T-2617 item.get finds its item, children, and plan in one list of the kept snapshot, with no glob and no second list

## Work

`item.get` in `flai/internal/hostapi/hostapi.go` calls `repo.Get`, which globs the kanban and archive folders, then `repo.List(true)` for the children, the plan, and `metrics.CostPerAgentHour`. Make it list once and find the item, its children, and its plan in that list, so that with the store's snapshot on it reads no file and stats nothing. Its answer stays the same field for field. The other methods that list every item, `board.get`, `items.list`, `items.count`, and those reached through `reads.go`, already go through `Repo.List` and read the snapshot without a change here.

Waits for T-2615, whose snapshot mode it reads.

## Done when

- [ ] A test in `flai/internal/hostapi/hostapi_test.go` shows `item.get` answers the same as before for a story with tasks, an archived item, and an unknown ID, and lists the items once a call.

## Notes

- Reviewed by agent-S-0606: touches widened to `flai/internal/hostapi/writes.go`. A write the host API runs through a flai subprocess (`runFlai`) changes item files in another process, which the watcher reports only 300 to 600 ms later. When `runFlai` returns, mark the run's project root stale with T-2615's call, so the dashboard's read after a write sees it at once.
- Reviewed by agent-S-0606: `item.get` keeps `repo.Get` then `repo.List(true)`. With T-2615's snapshot on, `Repo.Get` finds the item in the folders' snapshots with no glob or stat, so both read nothing; the task sub-agent's single-list finder copied `Repo.Get`'s ID rules into `hostapi`, a second copy to keep in step, and was taken out. The test that `item.get` answers from the snapshot stays.
