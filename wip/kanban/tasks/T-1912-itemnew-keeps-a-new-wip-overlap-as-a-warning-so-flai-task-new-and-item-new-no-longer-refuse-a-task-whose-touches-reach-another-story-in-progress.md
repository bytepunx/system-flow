---
id: T-1912
type: task
nature: feature
title: itemnew keeps a new wip.overlap as a warning, so flai task new and item_new no longer refuse a task whose touches reach another story in progress
status: done
parent: S-0451
owner: alex
created: 2026-10-20T19:04:09Z
updated: 2026-10-20T19:05:55Z
transitions:
  - to: ready
    at: 2026-10-20T19:04:18Z
    by: agent-S-0451
  - to: in-progress
    at: 2026-10-20T19:04:18Z
    by: agent-S-0451
  - to: done
    at: 2026-10-20T19:05:55Z
    by: agent-S-0451
stream: S-0451
tags: []
touches: [flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go]
usage:
  source: log
  seconds: 97
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 18
      output: 3688
      cache_read: 1114957
      cache_write: 44164
      cost: 0.6379
---
# T-1912 itemnew keeps a new wip.overlap as a warning, so flai task new and item_new no longer refuse a task whose touches reach another story in progress

## Work

In `flai/internal/itemnew/itemnew.go`, `Create` keeps a new `wip.overlap` finding as a warning in `Result.Warnings` instead of refusing the item, as `itemedit` does for an edit (ADR-0019, I-0059): an overlap between two stories in progress is advisory, and `flai task new` and `item_new` tell both stories through the claim watch. Add a test in `itemnew_test.go` that reproduces I-0169: two stories in progress, a new task of one whose touches overlap the other's claim, kept with the overlap as a warning. Waits for nothing.

## Done when

- The new test fails on main and passes here.
- `flai test` on `flai/internal/itemnew` passes.

## Notes
