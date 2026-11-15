---
id: T-3381
type: task
nature: improvement
title: workitem writes board.md under a lock, reloaded and normalised inside it, with a test that reproduces I-0375's lost update
status: backlog
parent: S-0790
owner: alex
created: 2026-11-15T21:11:42Z
updated: 2026-11-15T21:11:42Z
transitions: []
stream: S-0790
tags: [cli]
touches: [flai/internal/workitem/board.go, flai/internal/workitem/transition.go, flai/internal/workitem/board_update_test.go]
after: [T-3380]
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
          cache_write: 531
          cost: 0.0065
---
# T-3381 workitem writes board.md under a lock, reloaded and normalised inside it, with a test that reproduces I-0375's lost update

## Work

Build what T-3380's ADR decides, in `flai/internal/workitem`.

In `board.go`, add one way to change `board.md`, such as `func (r *Repo) UpdateBoard(today string, change func(*Board) error) error`. It takes the lock the ADR names (the items lock from `itemlock.go`, or a lock of its own beside it), loads the board inside the lock, runs `change`, drops from `order` every name that is not a ready or backlog story as the items read inside the lock have it, and saves before it lets go. `Board.Save` stays for the board a caller already holds under that lock; say so in its doc comment.

In `transition.go`, make `TransitionAll` change the board through it: the `RemoveFromOrder` and `PlaceReadyLast` its moves make (`rules.go`) are applied to the board loaded inside the lock, not to the one loaded before the items were saved. Keep the lock order the ADR gives with the items lock that `SaveAll` takes.

In the new `board_update_test.go`, reproduce I-0375: put three ready stories in the order, as S-0738, S-0746, and S-0739 were, and move each to in-progress with `TransitionAll` from its own goroutine at once. Afterwards the order names none of them. It fails before the change, run with `-count` high enough to show the lost update, and passes under `-race` after it. Add a test that a write drops a stale name, an in-progress story left in the order, and keeps the ready and backlog names in their order.

It waits for T-3380 so that the code and its comments cite the ADR's number and build what it decides.

## Done when

- Concurrent `TransitionAll` moves of stories in the order leave none of them in it, as the new test shows under `-race`.
- A write of `board.md` through the new function drops a name that is not a ready or backlog story.
- `flai test` on the changed paths passes.

## Notes
