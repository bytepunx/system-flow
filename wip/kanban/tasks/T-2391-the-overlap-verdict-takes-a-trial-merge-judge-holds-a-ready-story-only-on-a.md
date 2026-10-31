---
id: T-2391
type: task
nature: improvement
title: The overlap verdict takes a trial-merge judge, holds a ready story only on a conflicting trial merge, and names the conflicting paths in its reason
status: done
parent: S-0551
owner: alex
created: 2026-10-31T07:43:04Z
updated: 2026-10-31T08:56:45Z
transitions:
  - to: ready
    at: 2026-10-31T08:42:41Z
    by: agent-S-0551
  - to: in-progress
    at: 2026-10-31T08:42:42Z
    by: agent-S-0551
  - to: done
    at: 2026-10-31T08:56:45Z
    by: agent-S-0551
stream: S-0551
tags: [flai, go]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/repo.go, flai/internal/workitem/share.go, design/adrs/0230-a-ready-story-whose-claim-overlaps-a-story-in-progress-is-held-only-when-a.md, design/adrs/0231-flai-stats-reports-the-share-of-the-window-with-n-stories-in-progress-at-once.md]
usage:
  source: log
  seconds: 843
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 34
      output: 10675
      cache_read: 2099103
      cache_write: 76987
      cost: 1.0624
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
