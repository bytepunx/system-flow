---
id: T-3437
type: task
nature: feature
title: flai test --at ci, batch, or publish runs the tiers of the suites the workflow's verify_at lists, over all paths
status: backlog
parent: S-0812
owner: alex
created: 2026-11-16T23:47:16Z
updated: 2026-11-16T23:47:16Z
transitions: []
stream: S-0812
tags: [flai]
touches: [flai/cmd/test.go, flai/cmd/test_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/verify/verify.go, flai/internal/workflow/workflow.go, flai/internal/workflow/workflow_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
---
# T-3437 flai test --at ci, batch, or publish runs the tiers of the suites the workflow's verify_at lists, over all paths

## Work

ADR-0312 §6. Add `--at <point>` to `flai test` in `flai/cmd/test.go`. It takes `ci`, `batch`, or `publish`, reads the suites the project's workflow gives that point under `verify_at:`, as S-0806 parses them, and runs every tier of those suites for the whole checkout, as `--all` runs a tier: all-only tiers too, manual ones never. `all` means every tier but the manual ones, which is what the baseline's `ci: all` gives.

- A point the file leaves out, as the baseline leaves out `publish`, runs nothing and says so, exiting 0.
- An unknown point is refused, naming the three.
- `--at` given with paths, with a tier's name, or with `--all` is refused.
- With `--suite`, from S-0805, `--at` runs only the suites named that the point also lists. A suite the point does not list is refused. This lets one CI job run its share of the point.

Put the selection by suite over all paths in `flai/internal/verify/select.go`, beside the selection `--suite` uses. If the workflow needs a lookup by point, add it to `flai/internal/workflow/workflow.go`. Regenerate the reference and the flag index with `make flai-reference`.

This task waits for nothing in the story. It is the first layer.

## Done when

- Tests in `flai/cmd/test_test.go` and `flai/internal/verify/select_test.go` show that `--at ci` with `ci: all` runs every tier but the manual ones. They show that a point listing suites runs only those suites' tiers, all-only ones included, and that `--suite` narrows the point.
- The same tests show refusals for an unknown point, for paths, for a tier's name, for `--all`, and for a suite the point does not list.
- A point the file leaves out runs nothing and exits 0.
- `docs/users/flai-reference.md` and `docs/operators/settings.md` list `--at`, and `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent` pass.
- `flai test` on the changed paths passes.

## Notes
