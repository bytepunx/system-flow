---
id: T-3261
type: task
nature: feature
title: An integration test renders the template, adds a state with a rule and a timer to its workflow file, and drives the check, the board, the moves, serve, the acceptance, and stats through it
status: backlog
parent: S-0756
owner: alex
created: 2026-11-12T16:10:19Z
updated: 2026-11-12T16:10:19Z
transitions: []
stream: S-0756
tags: [cli]
touches: [flai/cmd/workflow_e2e_test.go, flai/internal/workflow/testdata/e2e-board.md]
---
# T-3261 An integration test renders the template, adds a state with a rule and a timer to its workflow file, and drives the check, the board, the moves, serve, the acceptance, and stats through it

## Work

Write `workflow_e2e_test.go`, skipped under `-short` as the integration tier's tests are: render the template into a temporary project, replace its workflow file with `testdata/e2e-board.md`, which adds a state between `in-progress` and `review` with a rule on the move into it and a timer from it to the acceptance. Then check that `flai check --strict` passes and fails, naming the state, when the state is left unreachable; that `flai board` shows the column; that `flai move` and `item_move` make and refuse the same moves; that a served story's agent keeps running in the new state; that `flai accept` merges and archives the story; and that `flai stats --json` reports the timer.

## Done when

The test passes in the integration tier, and fails if any one of those surfaces names a state.

## Notes
