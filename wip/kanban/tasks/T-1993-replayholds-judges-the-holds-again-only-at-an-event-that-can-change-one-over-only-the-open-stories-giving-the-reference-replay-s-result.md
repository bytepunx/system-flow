---
id: T-1993
type: task
nature: remediation
title: replayHolds judges the holds again only at an event that can change one, over only the open stories, giving the reference replay's result
status: done
parent: S-0471
owner: alex
created: 2026-10-21T12:08:09Z
updated: 2026-10-21T12:40:01Z
transitions:
  - to: ready
    at: 2026-10-21T12:21:27Z
    by: agent-S-0471
  - to: in-progress
    at: 2026-10-21T12:21:27Z
    by: agent-S-0471
  - to: done
    at: 2026-10-21T12:40:01Z
    by: agent-S-0471
stream: S-0471
tags: [flai]
touches: [flai/internal/metrics/claims.go, flai/internal/metrics/claims_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go]
after: [T-1991]
usage:
  source: log
  seconds: 1113
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 188
      output: 873
      cache_read: 9122319
      cache_write: 255957
      cost: 3.1776
---
# T-1993 replayHolds judges the holds again only at an event that can change one, over only the open stories, giving the reference replay's result

## Work

Build the fix T-1991 chose and recorded in the narrative's `## Decisions`. It waits for T-1991 because it needs its generator, benchmark, and reference replay to show the result is unchanged and the time is gone.

- In `flai/internal/metrics/claims.go`, change `replayHolds` so that it does not build a `workitem.Holds` and judge every ready story at every distinct event time since the first item. Judge again only when an event can change a hold, and carry the previous instant's held set and counts by reason over the others.
- Reuse each story's claim while its touches and its tasks' states stay the same, rather than calling `Holds.Claim` again at each rebuild.
- Change `flai/internal/workitem/hold.go` only where the replay needs a way in, such as building `Holds` from claims already worked out. Keep `Holds.Of`, `Claim`, `Overlaps`, and `PathsOverlap` giving what they give now: the board, the pull hold, `flai check`'s `wip.overlap`, and `flai serve` call them.
- In `flai/internal/metrics/claims_test.go`, add a test that the fixed replay gives the same `holdReplay` as `replayHoldsReference` on the generated history and on the existing small cases, and a test that counts the holds built in one replay and fails at today's count, one per event time. This is the test that reproduces I-0155.
- Run the benchmark again and record the time per replay against T-1991's in `## Decisions`. Aim for the stats read to be back within tens of milliseconds warm on this repository.

## Done when

- The equivalence test and the counting test pass, and `TestHeldSecondsReplayTheHoldRules` and the other tests in `claims_test.go` and `flai/internal/workitem/hold_test.go` pass unchanged.
- `flai test flai/internal/metrics flai/internal/workitem` passes.
- The benchmark's time per replay, before and after, is in `## Decisions`.

## Notes
