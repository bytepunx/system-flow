---
id: T-1349
type: task
nature: improvement
title: A selected tier that another selected tier covers is skipped, not run, with state skipped and the covering tier named
status: backlog
parent: S-0342
owner: alex
created: 2026-10-08T08:05:37Z
updated: 2026-10-08T08:05:37Z
transitions: []
stream: S-0342
tags: [cli]
touches: [flai/internal/verify/verify.go, flai/internal/verify/manifest.go, flai/internal/verify/manifest_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/verify/run.go, flai/internal/verify/run_test.go]
after: [T-1345]
---
# T-1349 A selected tier that another selected tier covers is skipped, not run, with state skipped and the covering tier named

## Work

Carry the key into a run and skip what it covers.

- `verify.Tier` gains `Covers`, and `FromManifest` copies it; a case in `manifest_test.go` shows it copied.
- A new state, `Skipped State = "skipped"`, and a field on `TierResult`, `CoveredBy string` with `json:"covered_by,omitempty"`, naming the tier that covers it.
- After `Select` and `SelectStory` choose their tiers, one step in `select.go` marks each selected tier that a selected tier covers, directly or through a chain, as covered by the nearest selected one, on `Selected` (`CoveredBy`). A covering tier that is itself covered still covers what it names, so the chain resolves to the tier that runs. A tier covered by one that was not selected is left as it is: alone it runs as before. `select_test.go` covers `Select` with `all`, `SelectStory` with and without the covering tier, and a chain.
- `Run` gives each covered tier a `TierResult` in its place in the list, state `skipped` and `covered_by` set, without running it, and goes on; a skipped tier never fails the run and is never the tier it stopped at. `run_test.go` shows a covered tier's command is not started and the result keeps the order.

It waits for T-1345, whose `TestTier.Covers` it reads.

## Done when

- With `integration` covering `go-test`, `Select(..., all=true)` and `SelectStory` for a path under `flai/` mark `go-test` covered by `integration`, and `Run` reports it `skipped` without starting it.
- `Select` for paths alone, which never selects an `all_only` tier, leaves `go-test` to run.
- `flai test ./internal/verify/` passes.

## Notes

Layer 2 of S-0342's plan. `run.go` and `run_test.go` are also in S-0341's touches.
