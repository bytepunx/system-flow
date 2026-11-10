---
id: T-2847
type: task
nature: improvement
title: itemedit.ApplyForecasts holds the checked-write lock from its save until it keeps or restores the forecasts
status: done
parent: S-0654
owner: alex
created: 2026-11-10T07:01:15Z
updated: 2026-11-10T07:44:20Z
transitions:
  - to: ready
    at: 2026-11-10T07:35:02Z
    by: agent-S-0654
  - to: in-progress
    at: 2026-11-10T07:35:04Z
    by: agent-S-0654
  - to: done
    at: 2026-11-10T07:44:20Z
    by: agent-S-0654
stream: S-0654
tags: [flai, itemedit, locking]
touches: [flai/internal/itemedit/forecasts.go, flai/internal/itemedit/forecasts_test.go]
usage:
  source: log
  seconds: 539
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 27
      output: 192
      cache_read: 3023572
      cache_write: 32922
      cost: 1.1245
---
# T-2847 itemedit.ApplyForecasts holds the checked-write lock from its save until it keeps or restores the forecasts

## Work

In `flai/internal/itemedit/forecasts.go`, `ApplyForecasts` take `repo.LockCheckedWrites()` just before `repo.SaveAll(edited)` and hold it through the index write, the check (`runItems`), and each `repo.Restore` of a refused item. Let go once every refused item is put back, or once the check keeps them all, on every return path, the failed index write's included. A lock it cannot take refuses the batch with nothing written, as a failed read does, so that flai serve's replan warns and tries again on its next round (`flai/internal/serve/replan.go`).

The replan is the only caller and runs outside the look's `readyChecked`, so it never holds the lock already; confirm it. Keep the lock order: the checked-write lock first, the items lock that `SaveAll` and `Restore` take inside it.

Add a test to `flai/internal/itemedit/forecasts_test.go`: replace `runItems` so that it tries `filelock.Lock` on `.flai-cache/items-checked.lock` with no wait and records that it could not take it, and refuses one item with an error finding on its file. Assert the probe found the lock held, the refused item was put back, and the lock is free once `ApplyForecasts` returns.

It waits for no task: no other task of S-0654 changes these files.

## Done when

- `ApplyForecasts` holds the checked-write lock from its save until it keeps or restores the forecasts, and lets go on every return.
- The new test passes, and `go test ./internal/itemedit/... ./internal/serve/...` passes from `flai/`.

## Notes
