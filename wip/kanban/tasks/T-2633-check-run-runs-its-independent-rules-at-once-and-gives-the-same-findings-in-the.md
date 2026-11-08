---
id: T-2633
type: task
nature: improvement
title: check.Run runs its independent rules at once and gives the same findings in the same order
status: in-progress
parent: S-0601
owner: alex
created: 2026-11-07T23:48:15Z
updated: 2026-11-08T01:11:52Z
transitions:
  - to: ready
    at: 2026-11-08T01:11:50Z
    by: agent-S-0601
  - to: in-progress
    at: 2026-11-08T01:11:52Z
    by: agent-S-0601
stream: S-0601
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/concurrent_test.go]
---
# T-2633 check.Run runs its independent rules at once and gives the same findings in the same order

## Work

`Run` in `flai/internal/check/check.go` calls about 30 rules one after another on one `checker`. They read only, but they share state the order hides:

- `workItems` fills `c.byID`, which `named` reads in `after`, `afterHub`, `afterReason`, and others. Run `workItems` first, alone, then the rest.
- `topicWords` sets `c.vocab` lazily. Compute it once before the rules start.
- `add`, `advise`, `adviseUnchecked`, and four rules that set `Stories` write the last finding by its index in `c.res.Findings`. Give each rule its own `Result`, so those writes stay in its own slice.

Run the rules after `workItems` on a bounded pool, each into its own `Result`. Then merge the results under one lock in the order `Run` lists the rules today, sum their counts, and sort as `sortFindings` does now. The sort is stable, so findings at one path and line keep their rule order.

`nodeRunner` is a package variable that `components_test.go` swaps. Leave it as it is: no rule writes it.

Add `flai/internal/check/concurrent_test.go`. It builds a fixture project in a temporary folder with at least one finding of every rule `Run` calls, the way the existing rule tests build theirs. It runs the rules one after another, the old order, and through `Run`, and checks that the findings, levels, counts, and order are the same. Run the check package's tests under `-race`.

This task waits for none.

## Done when

- `Run` runs every rule but `workItems` at once and merges their findings in the old order.
- `concurrent_test.go` passes, and so does `go test -race ./internal/check/...`.

## Notes
