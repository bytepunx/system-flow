---
id: T-3141
type: task
nature: remediation
title: A move to done or cancelled closes the item's open blocked interval at the move
status: done
parent: S-0723
owner: alex
created: 2026-11-12T15:53:22Z
updated: 2026-11-12T16:54:29Z
transitions:
  - to: ready
    at: 2026-11-12T16:46:50Z
    by: agent-S-0723
  - to: in-progress
    at: 2026-11-12T16:46:50Z
    by: agent-S-0723
  - to: done
    at: 2026-11-12T16:54:29Z
    by: agent-S-0723
stream: S-0723
tags: [cli]
touches: [flai/internal/workitem/rules.go, flai/internal/workitem/rules_test.go]
usage:
  source: log
  seconds: 459
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 2627
      cache_read: 1068776
      cache_write: 25250
      cost: 0.4683
---
# T-3141 A move to done or cancelled closes the item's open blocked interval at the move

## Work

This is the cause I-0347 records. S-0678 was blocked at 2026-11-11T05:26:16Z ("pause for urgent fix"), went to review at 06:13:48Z, and the orchestrator accepted it at 06:15:07Z with the interval still open. `Repo.Move` in `flai/internal/workitem/rules.go` refuses nothing and closes nothing for a blocked item on its way to `done` or `cancelled`. Every close-out since has reported the archived S-0678 as `item.blocked-closed` outside its story (32 instances), and `flai stats` counts its blocked time to now, since `metrics.md` ends an open interval at now.

In `Repo.Move`, when `to` is `done` or `cancelled` and the item has an open blocked interval, close it at the move's own timestamp, the `now` the transition is stamped with, and return a warning that names the item and the interval's reason. Do not refuse the move: acceptance is the operator's or the orchestrator's, and a block left open must not stop it. Every path that closes an item goes through `Repo.Move`: `flai move`, `flai accept` (`flai/cmd/accept.go`), the merge queue, the MCP `item_move`, the dashboard's writes, and a cancel's cascade to the open items under it. So this one change covers them all.

Add tests to `flai/internal/workitem/rules_test.go` that reproduce the cause: a blocked story moved from review to done, and a blocked task cancelled by its story's cascade, each ends with no open interval, `Until` equal to its last transition's `at`, and the warning given.

Waits for nothing: it is the first layer.

## Done when

- A move to `done` or `cancelled` leaves no open blocked interval, closed at the move's timestamp, with a warning naming it.
- The new tests in `flai/internal/workitem/rules_test.go` fail without the change and pass with it.
- `flai test` on the two paths passes.

## Notes
