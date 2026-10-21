---
id: T-2002
type: task
nature: remediation
title: flai serve's usage writer, flai edit and item_edit, touches, and widening write their items under the items lock from a fresh read
status: backlog
parent: S-0479
owner: alex
created: 2026-10-21T12:10:53Z
updated: 2026-10-21T12:10:53Z
transitions: []
stream: S-0479
tags: [flai, serve, itemedit]
touches: [flai/internal/serve/usage.go, flai/internal/serve/usage_test.go, flai/internal/itemedit/itemedit.go, flai/internal/itemedit/widen.go, flai/internal/itemedit/lock_test.go, flai/cmd/touches.go]
after: [T-1999]
---
# T-2002 flai serve's usage writer, flai edit and item_edit, touches, and widening write their items under the items lock from a fresh read

## Work

Move the edit and usage writers onto the lock T-1999 adds. They wait for T-1999 because they call its `Repo.Update`.

- `put` in `flai/internal/serve/usage.go` reads the item again before it saves, but with no lock. Run its read, merge, and save inside the items lock.
- `itemedit.Apply` in `flai/internal/itemedit/itemedit.go` reads the item at its start, runs `check.Run` over the project, and saves its early copy seconds later. This is the window the replanner's forecast write left open in I-0175. Do not hold the lock across `check.Run`. Validate as now, then under the lock read the item, and the parents it saves, again. Compare the caller's `hash` with that fresh read, so the hash check guards the write. Apply the change to the fresh copies and save them. The replanner's call in `flai/internal/serve/replan.go` then needs no change of its own.
- Do the same in `Widen` in `flai/internal/itemedit/widen.go` and in `flai touches` in `flai/cmd/touches.go`.
- In `flai/internal/itemedit/lock_test.go`, reproduce the forecast case. Start an `Apply` of a forecast, move the item to in-progress while `Apply` is between its read and its save (a hook in `Options` or a slow check runner), and assert that both the forecast and the transition stay. In `flai/internal/serve/usage_test.go`, assert that a usage write keeps a transition saved after the usage was measured.

## Done when

- Every `Save` in the files this task touches runs inside the items lock on an item read under it, and `Apply`'s hash check compares against that read.
- The new tests fail on the code before this story and pass after.
- `flai test` on the touched paths passes.

## Notes

Drafted by the planner. It runs beside T-2001, since the two touch no file in common.
