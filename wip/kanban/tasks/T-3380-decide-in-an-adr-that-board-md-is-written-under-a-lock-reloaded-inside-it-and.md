---
id: T-3380
type: task
nature: improvement
title: Decide in an ADR that board.md is written under a lock, reloaded inside it, and normalised on every write
status: done
parent: S-0790
owner: alex
created: 2026-11-15T21:11:30Z
updated: 2026-11-15T21:20:04Z
transitions:
  - to: ready
    at: 2026-11-15T21:18:46Z
    by: system-flow
  - to: in-progress
    at: 2026-11-15T21:18:47Z
    by: system-flow
  - to: done
    at: 2026-11-15T21:20:04Z
    by: agent-S-0790
stream: S-0790
tags: [design]
touches: [design/adrs/README.md, design/system/workflow.md, design/system/flai-cli.md]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 7
          cache_read: 20828
          cache_write: 531
          cost: 0.0065
---
# T-3380 Decide in an ADR that board.md is written under a lock, reloaded inside it, and normalised on every write

## Work

Write the decision I-0375's instances call for, with `flai adr new`, and bring the design in line with it.

The cause, from the history of `wip/kanban/board.md` and the archived items: S-0738, S-0746, and S-0739 moved to in-progress in the same second, 2026-11-14T23:49:15Z, when `flai serve` started their agents after S-0737's acceptance. Each start ran `TransitionAll` (`flai/internal/workitem/transition.go`), which loads the board, removes its own story from `order` with `RemoveFromOrder`, and saves it. `Board.Save` (`flai/internal/workitem/board.go`) takes no lock, unlike an item's `Save`, which takes the items lock (ADR-0184). So one start's save wrote back a board read before another start removed its story: S-0738 stayed in the order while in progress, and the close-outs of S-0738 and S-0746 found `board.order` outside the story. The commit `449ee65da`, S-0789's creation 15 seconds later, still names S-0738.

The ADR decides:

- Every write of `board.md` takes a lock, reloads the board inside it, applies its change, and saves before it lets go. The writes are `TransitionAll`, `flai order` (both its placement and `--apply`), `flai accept`, and `flai board limit`. Say whether it is the items lock or a lock of its own, and the lock order with the items lock, so that a write holding one never waits on a writer holding the other.
- Every such write normalises `order` against the items read inside the lock: it drops a name that is no longer a ready or backlog story. A stale name left by any earlier write then goes at the next write.
- `board.order` stays a finding a close-out records: it names a real defect, not a state only the operator clears, so it is not left out of a check scoped to a story as ADR-0133 leaves out `board.wip-limit`.

In `design/system/workflow.md` § The pull order (S-0057), the sentence "On every write the list is normalised" is true only of `flai order` today: say that every write of `board.md` normalises it, under the lock, and cite the ADR. Update `design/system/flai-cli.md` where it describes how `board.md` is written, if it does. `flai adr new` regenerates `design/adrs/README.md`.

It waits for no task: T-3381 and T-3382 build what it decides.

## Done when

- A proposed ADR records the lock, the reload, the normalisation, and why `board.order` stays recorded, citing I-0375 and ADR-0184.
- `design/system/workflow.md` says every write of `board.md` normalises the order under the lock and links the ADR.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes
