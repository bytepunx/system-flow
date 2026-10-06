---
id: T-1030
type: task
nature: improvement
title: An overlap wholly inside a shared path neither holds a ready story nor counts as a wip.overlap or a grown claim
status: backlog
parent: S-0295
owner: alex
created: 2026-10-06T12:16:01Z
updated: 2026-10-06T12:16:01Z
transitions: []
stream: S-0295
tags: [flai]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/itemedit/claim.go, flai/internal/itemedit/claim_test.go, flai/internal/workitem/boardview.go]
after: [T-1025, T-1027]
---
# T-1030 An overlap wholly inside a shared path neither holds a ready story nor counts as a wip.overlap or a grown claim

## Work

`Holds` receives the manifest's shared patterns, from `NewHolds` and `Repo.Holds`. In `holdBy`, an overlapping pair of entries does not hold when the narrower of the two lies wholly inside a shared pattern, by T-1025's matcher. A folder entry that only partly lies in a shared path still holds.

Apply the same exception in two more places:

- `flai check`'s `wip.overlap` (`flai/internal/check/check.go`).
- The notice that a write grew a claim onto another open story's (`flai/internal/itemedit/claim.go`).

Leave `flai stream sync`'s trial merge and `flai accept`'s overlap notice as they are, so that a real conflict on a shared file is still reported. Keep the board view's held reason (`boardview.go`) free of shared paths.

This task waits for T-1025 for the matcher, and for T-1027, which changes `hold.go` and the claim before it.

## Done when

- A test reproduces I-0087's S-0283 and S-0285 holds through `docs/users/flai.md` and `design/adrs`. With those paths shared, no ready story is held by them. With the list empty, they hold as before.
- A test shows `flai check` gives no `wip.overlap` and the grown-claim notice stays silent on an overlap inside a shared path, and that both still report one outside it.
- `scripts/flai-test.sh` passes, the race detector included.

## Notes
