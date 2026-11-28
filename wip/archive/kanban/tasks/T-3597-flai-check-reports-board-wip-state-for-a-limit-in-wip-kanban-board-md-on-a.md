---
id: T-3597
type: task
nature: remediation
title: flai check reports board.wip-state for a limit in wip/kanban/board.md on a state that takes none
status: done
parent: S-0844
owner: alex
created: 2026-11-17T12:36:49Z
updated: 2026-11-28T00:22:07Z
transitions:
  - to: ready
    at: 2026-11-27T23:38:20Z
    by: agent-S-0844
  - to: in-progress
    at: 2026-11-27T23:38:20Z
    by: agent-S-0844
  - to: done
    at: 2026-11-28T00:22:07Z
    by: agent-S-0844
stream: S-0844
tags: [flai, check]
touches: [flai/internal/check/check_items.go, flai/internal/check/items_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/internal/check/itemwrite_test.go]
usage:
  source: log
  seconds: 2627
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 73
      output: 18348
      cache_read: 5559955
      cache_write: 152987
      cost: 2.3869
---
# T-3597 flai check reports board.wip-state for a limit in wip/kanban/board.md on a state that takes none

## Work

This task waits for no other task. It shares no file with the other tasks of the first layer.

ADR-0301 §6 names the finding `board.wip-state`: a limit in `wip/kanban/board.md` on a state that takes none. No code reports it (I-0380, instance 3). It matters because `Board.Save` (`flai/internal/workitem/board.go`) drops such a key without a word on the next write of the board. A limit written by hand, or one left behind when a story takes `limit: true` off a state or removes the state, then disappears with no warning.

- In `flai/internal/check/check_items.go`, `(*checker).board`, loop over the keys of the board's `WIPLimits`. Report `board.wip-state` as a warning on `wip/kanban/board.md`, at the key's line, for each key that names no state of the workflow with `limit: true`. Match names and aliases the way `wipLimit` does. Say in the message that the next write of the board drops the limit, and that the fix is to give the state `limit: true` or remove the line. Skip the loop while `c.flowRefused()` holds.
- Make it a warning, as `board.wip-limit` is. The ADR gives no severity, and the finding stops no move.
- In `flai/internal/check/scope.go`, decide whether a check scoped to a story leaves out a `board.wip-state` outside it, as it does `board.wip-limit`. The board file is main's, and no story's branch changes it. Record the choice in the narrative's `## Decisions` and test it in `scope_test.go`.
- Add tests to `flai/internal/check/items_test.go`: a limit on a limited state reports nothing, and a limit on a state without `limit: true` reports `board.wip-state` at its line. Do the same for a state the workflow names only as an alias, and for a name the workflow lacks.

## Done when

- `flai check` on a board whose `wip_limits` names a state without `limit: true` reports `board.wip-state` at that line, as a warning.
- A board whose limits are all on limited states reports nothing new.
- The scoped check treats the finding as the narrative records, and a test covers that.
- `flai test` on the changed paths passes.

## Notes
