---
id: T-2644
type: task
nature: feature
title: verify.Run starts the tiers the manifest marks independent at once and stops the later ones when one fails
status: done
parent: S-0610
owner: alex
created: 2026-11-07T23:48:50Z
updated: 2026-11-09T08:08:29Z
transitions:
  - to: ready
    at: 2026-11-09T07:52:33Z
    by: agent-S-0610
  - to: in-progress
    at: 2026-11-09T07:52:34Z
    by: agent-S-0610
  - to: done
    at: 2026-11-09T08:08:29Z
    by: agent-S-0610
stream: S-0610
tags: [cli]
touches: [flai/internal/verify/run.go, flai/internal/verify/run_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go]
after: [T-2634]
usage:
  source: log
  seconds: 952
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 38
      output: 239
      cache_read: 5557979
      cache_write: 363120
      cost: 2.1511
---
# T-2644 verify.Run starts the tiers the manifest marks independent at once and stops the later ones when one fails

## Work

Make `verify.Run` (`flai/internal/verify/run.go`) run at once the selected tiers that T-2634's key lets run beside each other, for criterion 1. It waits for T-2634, whose key and `verify.Tier` field it reads.

- Keep the list order: results stay in the order `Select` gave, cheapest first, and a tier starts no earlier than its key allows.
- When a tier fails, stop through its context every tier after it in list order that is running, and start none after it. Let a running tier before it finish, so that the tier `stopped_at` names is the first in list order that fails, as a serial run's is. A stopped tier is listed as not reached, never failed, and its findings are not counted against `opts.Max`; the failing tiers' findings are capped in list order.
- Each tier keeps its own output buffers and its own log file, as now. The tiers' `Proc.Run` already stops a command's process group when its context ends.
- `Run` returns `ctx`'s error when the caller's context ends, as now, stopping every running tier.
- The close-out's call in `storyRun.run` (`flai/internal/verify/story.go`) needs no change of shape; check that reused tiers and `RanFor` still line up with the results. `Tiers`, the merge queue's batch, runs its tiers one at a time through `Run` with its `Started` and `Ended` callbacks: leave it serial unless running it through one `Run` keeps those callbacks true, and say which in the narrative's `## Decisions`.
- Test with a fake `Proc`: two independent tiers overlap in time; a dependent tier waits; a later tier's failure stops nothing before it and the earlier failure is the one named; a failure stops the running later tier and it is listed not reached; the verdict and `stopped_at` match a serial run's over the same results; and the race detector is clean.

## Done when

- The tests above pass under `-race`.
- A run with no tier marked independent behaves as before: the existing tests of `run_test.go` and `story_test.go` pass unchanged.

## Notes
