---
id: T-3573
type: task
nature: improvement
title: flai serve opens, resolves, and adds to the starved board's thread only once the board has held its state for the settle window
status: backlog
parent: S-0847
owner: alex
created: 2026-11-17T12:29:41Z
updated: 2026-11-17T12:29:50Z
transitions: []
stream: S-0847
tags: [flai, serve]
touches: [flai/internal/serve/starved.go, flai/internal/serve/starved_test.go, flai/internal/serve/readyempty_test.go, docs/users/flai.md]
after: [T-3569]
---
# T-3573 flai serve opens, resolves, and adds to the starved board's thread only once the board has held its state for the settle window

## Work

This task waits for T-3569: it builds what the ADR that task writes decides, and the window, the opening, the resolving, and the coalescing may change on the operator's answer.

In `flai/internal/serve/starved.go`, as the ADR decides:

- Add to `starvedState` when the board was first found starved and, while the thread is open, when it was first found not starved, so that both survive a restart in `.flai-cache/starved.json`.
- `noteStarved` and `starve` open the thread only at a look that finds the board starved for the whole window. `unstarve` resolves it only at a look that finds the board not starved for the whole window, and a starving within the window keeps the same thread.
- `starve` writes a changed list at most once per window since flai serve's last entry, the newest list then.
- Keep the rest of ADR-0285 as it is: the anchor, the author, the title, the list's text, the reminder channel's line, and a thread the operator resolved not being opened again until the list changes.

Write a test in `flai/internal/serve/starved_test.go` that reproduces I-0389 on `newStarvedLab`'s clock. A board flips between starved and not at each look, as TH-1135 and TH-1166 did: no thread opens while no state lasts the window. A board starved past the window that flips back for less than it keeps one thread. A list that changes at each look, as TH-1134's did, gains at most one entry per window. Change `TestAStarvedBoardOfI0288TellsTheOperatorOnce`, which expects a thread at the first look and a new one at the next starving, to the ADR's rules. Change `TestReadyEmptyAndTheStarvedBoardCountTheActiveWorkRoom` in `flai/internal/serve/readyempty_test.go` where it expects the thread at once.

Rewrite the starved-board paragraph of `docs/users/flai.md` § flai serve to say when the thread opens, gains an entry, and is resolved, linking the ADR.

## Done when

- The test reproducing I-0389 passes. A board flipping faster than the window opens no thread. A board starved past the window that flips back briefly keeps its one thread. A churning list gains at most one entry per window.
- The changed starved-board tests pass, and the rest of ADR-0285's tests pass unchanged.
- `docs/users/flai.md` describes the window and links the ADR.
- `flai test` on the changed paths passes.

## Notes
