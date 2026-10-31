---
id: T-2377
type: task
nature: improvement
title: The go-test tier runs the changed packages and every package that imports them, and a test shows a change that breaks a caller's tests fails verify
status: backlog
parent: S-0547
owner: alex
created: 2026-10-31T07:39:42Z
updated: 2026-10-31T07:39:42Z
transitions: []
stream: S-0547
tags: [flai]
touches: [system-flow.yaml, flai/internal/verify/story_test.go, flai/internal/verify/importers_test.go]
---
# T-2377 The go-test tier runs the changed packages and every package that imports them, and a test shows a change that breaks a caller's tests fails verify

## Work

Criterion 1. flai already fills `{importers}` with the changed packages and every package of the module that imports them, directly or through others, in their own or their test files (S-0364, `flai/internal/verify/importers.go`); the `vet` tier uses it. Change this repository's `go-test` tier command in `system-flow.yaml` from `{packages}` to `{importers}`, leaving its other arguments as they are on main when the task starts (S-0548 changes the same line for `-timeout` and `-short`).

Add a test in `flai/internal/verify/story_test.go` that builds a two-package Go module in `t.TempDir()`, a callee and a caller whose test calls it, commits it on a branch that changes only the callee in a way that breaks the caller's test, and runs the story verification with a `go-test` tier on `{importers}`: the run fails at that tier with a finding in the caller's package. Show the same run with `{packages}` passes, which is I-0221's defect. Extend `importers_test.go` only if the fixture finds a case `importers` misses.

Waits for nothing: no other first-layer task touches these files.

## Done when

- `system-flow.yaml`'s `go-test` tier passes `{importers}`.
- The new test fails the verification on the caller's broken test with `{importers}` and passes it with `{packages}`, and is not skipped under `-short`.
- `flai test` on the changed paths passes.

## Notes

Drafted by planner-S-0547. Reading criterion 1's "`{packages}` expands to" as this repository's go-test tier taking `{importers}`, rather than changing what `{packages}` means for every project and tier, is an assumption on the plan's thread.
