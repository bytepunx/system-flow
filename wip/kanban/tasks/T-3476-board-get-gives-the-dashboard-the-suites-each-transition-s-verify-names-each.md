---
id: T-3476
type: task
nature: feature
title: board.get gives the dashboard the suites each transition's verify names, each state's runs, and verify_at's batch, publish, and ci
status: backlog
parent: S-0814
owner: alex
created: 2026-11-16T23:50:53Z
updated: 2026-11-16T23:50:53Z
transitions: []
stream: S-0814
tags: [flai]
touches: [flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go]
---
# T-3476 board.get gives the dashboard the suites each transition's verify names, each state's runs, and verify_at's batch, publish, and ci

## Work

Criterion 3, the host API's part. Give `board.get` what the workflow file says each step needs and runs, as S-0806 parses it: for each transition with `verify:`, its from, to, suites (or `all`), and severity; for each state with `runs:`, its suites and its `kluster:` steps; and `verify_at:`'s `batch`, `publish`, and `ci`. Build it in `NewBoardView` beside the columns S-0753 gives, from the project's workflow, so the baseline answers `verify: all` on the move to review and the acceptance, and `runs: {kluster: [deploy, checks]}` on `review`.

Waits for no task: it reads only the parsed workflow, which is on main once S-0806 is accepted.

## Done when

- `board.get` answers the verify, runs, and verify_at blocks for the baseline and for a fixture workflow with an integration state that runs suites.
- A workflow with none of the keys answers empty blocks, and an older dashboard that ignores them reads the board as before.
- `boardview_test.go` and `hostapi_test.go` cover both, and `flai test` on the changed paths passes.

## Notes

- Drafted by planner-S-0814. `board.get` is in `hostapi.go`, and `NewBoardView` in `boardview.go` builds its answer.
