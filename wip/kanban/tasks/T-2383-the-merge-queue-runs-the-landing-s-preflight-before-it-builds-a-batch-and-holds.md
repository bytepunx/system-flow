---
id: T-2383
type: task
nature: improvement
title: The merge queue runs the landing's preflight before it builds a batch and holds itself on a finding, with the files and findings in its state and event log
status: backlog
parent: S-0549
owner: alex
created: 2026-10-31T07:41:02Z
updated: 2026-10-31T07:42:31Z
transitions: []
stream: S-0549
tags: [flai]
touches: [flai/internal/mergequeue/preflight.go, flai/internal/mergequeue/preflight_test.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/log.go, flai/internal/mergequeue/log_test.go]
after: [T-2381, T-2382]
---
# T-2383 The merge queue runs the landing's preflight before it builds a batch and holds itself on a finding, with the files and findings in its state and event log

## Work

In a new `flai/internal/mergequeue/preflight.go`, a function that runs, on the main checkout, the fixing tiers' check over its uncommitted files (`verify.CheckoutUnfixed`, with the findings of the findings task) and `check.Run` on main, keeping the findings the ADR says hold the queue, and returns them with each file, line, and rule. It changes no file: the operator's uncommitted files are theirs.

Call it from `RunBatch` (`batch.go`) before a batch is built, with as small an edit there as it needs: S-0552 changes `RunBatch` and `pick` for its coalescing window, and whichever lands second rebases (MS-0263). On a finding, record the hold in the queue's state (`state.go`) in a field of its own beside the operator's `Held`, with the findings, when it began, and room for the ID of the thread `flai serve` opens for it; build nothing, and append an event to the event log (`log.go`). On a clean run, clear it and append the event that it cleared. Keep `flai queue resume` from clearing it.

Tests: a main checkout with an uncommitted note that breaks a markdown rule holds the queue with that file and rule, builds no batch, and runs no tier; fixing the note and looking again clears the hold and builds the batch. Story criterion 1's test.

It waits for the ADR task, which settles what holds the queue, and the findings task, whose findings it reports.

## Done when

`RunBatch` builds no batch while the preflight finds anything, the state and the event log name each file and finding, a clean look clears the hold, and `go test -race ./internal/mergequeue` passes.

## Notes

Drafted by the planner of S-0549 on 2026-10-31.
