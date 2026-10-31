---
id: T-2384
type: task
nature: improvement
title: A verified batch whose landing failed on main's state is kept, not dropped, and lands again from its verified commits once main is clean, with no new verification
status: backlog
parent: S-0549
owner: alex
created: 2026-10-31T07:41:20Z
updated: 2026-10-31T07:41:20Z
transitions: []
stream: S-0549
tags: [flai]
touches: [flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/log.go]
after: [T-2383]
---
# T-2384 A verified batch whose landing failed on main's state is kept, not dropped, and lands again from its verified commits once main is clean, with no new verification

## Work

Today a failed `flai accept --batch` in `Land` (`land.go`) calls `dropBatch` with `landing <story> failed: ...`, so the next look builds the batch again and verifies it again: S-0541 went through 16 batches and 28 failed landings this way.

Tell a landing that failed on main's state (the preflight's findings, or the acceptance's own refusal of the main checkout's unfixed files) from one that failed on its stories. For the first, keep the batch with its `Landing` and its verified heads (`Heads[].Verified`, `Rebased`), mark it waiting on main in `state.go`, and land it again from those commits once the preflight passes, with no tier run, as long as main's head is still the base it was verified against and each story's branch is still at its verified head (`Landable`). When main moved meanwhile, drop it as today. Append an event to the event log (`log.go`) that says the batch landed again without verification, so `flai serve` can journal it.

Keep the edit to `land.go` as small as it can be: S-0552 changes `flai/internal/mergequeue` too (MS-0263).

Tests: a passing batch whose landing is refused for an uncommitted file's lint is kept, lands from its verified commit once the file is fixed, runs no tier, and logs the event; a batch whose main moved meanwhile is dropped as before.

It waits for the preflight task, which writes `state.go` and `log.go` first and gives the check this task runs again.

## Done when

A verified batch kept on main's state lands again with no new verification once main is clean, the event log says so, a moved main still drops it, and `go test -race ./internal/mergequeue` passes. Story criterion 3.

## Notes

Drafted by the planner of S-0549 on 2026-10-31.
