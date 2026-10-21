---
id: T-2001
type: task
nature: remediation
title: Moves, story start, blocks, acceptance, creation, and usage roll-ups write their items under the items lock from a fresh read
status: backlog
parent: S-0479
owner: alex
created: 2026-10-21T12:10:33Z
updated: 2026-10-21T12:11:16Z
transitions: []
stream: S-0479
tags: [flai, workitem]
touches: [flai/internal/workitem/transition.go, flai/internal/workitem/create.go, flai/internal/workitem/usage.go, flai/internal/storystart/start.go, flai/internal/storystart/start_test.go, flai/cmd/move.go, flai/cmd/accept.go]
after: [T-1999]
---
# T-2001 Moves, story start, blocks, acceptance, creation, and usage roll-ups write their items under the items lock from a fresh read

## Work

Move the state writers onto the lock T-1999 adds. They wait for T-1999 because they call its `Repo.Update`.

- `Repo.TransitionAll` in `flai/internal/workitem/transition.go` saves the `it` its caller loaded earlier (`flai/internal/storystart/start.go` loads it before `Begin`). Under the lock, read the item, the items list, and the board again. Apply the move, cascade, and epic follow to the fresh copies, and save them and the board before releasing the lock. Copy the fresh state back into the caller's `it`, or return it, so callers see what was written.
- In `flai/internal/storystart/start.go`, `refuse` must judge the fresh copy under the lock, so that two starts of one story cannot both pass.
- Do the same for block and unblock in `flai/cmd/move.go`, the saves in `flai/cmd/accept.go`, the parent save in `flai/internal/workitem/create.go`, and `RollUp` and `ChargeStrategic` in `flai/internal/workitem/usage.go`. Re-read each item under the lock and change only the fields the writer owns.
- In `flai/internal/storystart/start_test.go`, reproduce I-0175's instance. Load the ready story, let another writer save a usage or forecast change to it, then start it, and the reverse order. Assert that the story is in progress with its transition and keeps the other writer's field.

## Done when

- Every `Save` in the files this task touches runs inside the items lock on an item read under it.
- The `start_test.go` test fails on the code before this story and passes after.
- `flai test` on the touched paths passes.

## Notes

Drafted by the planner. It runs beside T-2002, since the two touch no file in common.
