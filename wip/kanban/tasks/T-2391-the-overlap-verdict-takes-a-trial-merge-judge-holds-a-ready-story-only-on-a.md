---
id: T-2391
type: task
nature: improvement
title: The overlap verdict takes a trial-merge judge, holds a ready story only on a conflicting trial merge, and names the conflicting paths in its reason
status: backlog
parent: S-0551
owner: alex
created: 2026-10-31T07:43:04Z
updated: 2026-10-31T07:43:04Z
transitions: []
stream: S-0551
tags: [flai, go]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go]
---
# T-2391 The overlap verdict takes a trial-merge judge, holds a ready story only on a conflicting trial merge, and names the conflicting paths in its reason

## Work

In `flai/internal/workitem/hold.go`, let the holds take an optional trial-merge judge, a function given the ready story and the in-progress story it overlaps that answers clean, the conflicting paths, or that it could not run. `workitem` cannot import `storygit`, which imports it, so the judge is a function type declared here. In `overlap` and `holdBy`, an overlap with an in-progress story whose trial merge is clean does not hold; one that conflicts holds with the code `overlap` and a reason naming the conflicting paths; one the judge could not run holds as today. With no judge, the verdict is today's, so that callers not yet wired keep their behaviour. The `after` verdict is unchanged.

It waits for nothing: the judge is a function, and the git side is a task of its own.

## Done when

- Tests in `hold_test.go` show a clean judge leaving the ready story unheld, a conflicting judge holding it with the conflicting path in the reason, a judge that could not run holding it as today, and no judge keeping today's verdict.
- `flai test flai/internal/workitem` passes.

## Notes

The explorer found the verdict in `overlap` (about line 649), `holdBy` (706), and `firstHolding` (741), and the constructors `NewHolds` (132) and `Repo.Holds` (160).
