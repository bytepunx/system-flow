---
id: T-2939
type: task
nature: remediation
title: The touches drift leaves out the digest, analysis, issues, and wip folders on both sides, with a test that reproduces I-0286
status: ready
parent: S-0674
owner: alex
created: 2026-11-10T08:50:01Z
updated: 2026-11-11T01:43:21Z
transitions:
  - to: ready
    at: 2026-11-11T01:43:21Z
    by: system-flow
stream: S-0674
tags: [flai, metrics]
touches: [flai/internal/storygit/committed.go, flai/internal/storygit/committed_test.go, flai/internal/metrics/claims.go, flai/internal/metrics/claims_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, design/system/metrics.md]
after: [T-2936]
---
# T-2939 The touches drift leaves out the digest, analysis, issues, and wip folders on both sides, with a test that reproduces I-0286

## Work

Build decision 1 of T-2936's ADR. It waits for T-2936 because it follows that ADR.

- In `flai/internal/storygit/committed.go`, `readCommitted` also leaves out the files under the design folder's `digest`, `analysis`, and `issues` folders. It reads them from `repo.Manifest.Layout` as it reads `wip`. Put the left-out prefixes into the read's key (`CommittedKeyed`). Otherwise a `stories-commits.json` kept under the old rules answers the old files.
- In `flai/internal/metrics/claims.go`, `touchDrift` leaves the same prefixes and `wip/` out of the touches it compares, so that a touch there is not counted `unchanged`. Pass the prefixes through `metrics.Options` (`metrics.go`), set by `statsread.go`.
- Reproduce I-0286 in `committed_test.go` and `claims_test.go`. Make a story whose own commit changes `flai/x.go` and whose touches are `flai/x.go` and its narrative under `wip/`. Add a `chore: [S-nnnn] accept and archive` commit adding `design/digest/D-0001-x.md`, and a `docs: [S-nnnn] record I-0001` commit adding an issue file. Before the fix the story is not exact. After it, `outside` and `unchanged` are both empty.
- Rewrite the paragraph under the table in `design/system/metrics.md` § Claims and touches to say which folders are left out on each side, and link the ADR.

## Done when

- Both tests fail without the change and pass with it.
- `flai test` on the changed paths passes.
- `flai stats --json` in the main checkout's tree shows `claims.drift[]` with no `outside` file under the four folders.
- `design/system/metrics.md` § Claims and touches matches the code.

## Notes
