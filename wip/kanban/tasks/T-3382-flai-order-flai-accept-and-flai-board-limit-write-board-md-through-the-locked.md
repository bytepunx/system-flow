---
id: T-3382
type: task
nature: improvement
title: flai order, flai accept, and flai board limit write board.md through the locked update
status: ready
parent: S-0790
owner: alex
created: 2026-11-15T21:11:50Z
updated: 2026-11-15T21:18:46Z
transitions:
  - to: ready
    at: 2026-11-15T21:18:46Z
    by: system-flow
stream: S-0790
tags: [cli]
touches: [flai/cmd/order.go, flai/cmd/accept.go, flai/cmd/board.go, flai/cmd/order_test.go]
after: [T-3381]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 20828
          cache_write: 530
          cost: 0.0065
---
# T-3382 flai order, flai accept, and flai board limit write board.md through the locked update

## Work

Move the other writers of `board.md` onto the function T-3381 adds, so that no write loads the board outside the lock and saves it later.

- `flai/cmd/order.go`: the placement (`LoadBoard` then `board.Save`, near line 92 to 101) and `--apply` (near line 152 to 160) each make their change inside the locked update. `ApplyReadyOrder` and `ReadyOrderKeeping` (`flai/internal/workitem/policy.go`) read the board and the items before the write: under `--apply`, work the order out from the board and items read inside the lock.
- `flai/cmd/accept.go`: the move to done (`LoadBoard` near line 607, `Move` near 614, `board.Save` near 644) applies its board change inside the locked update.
- `flai/cmd/board.go`: `flai board limit` (near line 127 to 134) sets the limit inside it.

In `flai/cmd/order_test.go`, add a test that `flai order --by <policy> --apply`, run while a story in the order has moved to in-progress since the board was last written, leaves that story out of the order. Run the existing tests of the three commands.

It waits for T-3381, whose function these call.

## Done when

- No code outside the locked update calls `Board.Save`, as a search of `flai/` for `.Save(` on a board shows.
- The new test and the existing tests of `flai order`, `flai accept`, and `flai board` pass.
- `flai test` on the changed paths passes.

## Notes
