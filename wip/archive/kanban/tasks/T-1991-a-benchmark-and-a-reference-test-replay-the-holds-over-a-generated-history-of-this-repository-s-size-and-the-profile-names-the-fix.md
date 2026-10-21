---
id: T-1991
type: task
nature: remediation
title: A benchmark and a reference test replay the holds over a generated history of this repository's size, and the profile names the fix
status: done
parent: S-0471
owner: alex
created: 2026-10-21T12:07:50Z
updated: 2026-10-21T12:21:16Z
transitions:
  - to: ready
    at: 2026-10-21T12:13:52Z
    by: agent-S-0471
  - to: in-progress
    at: 2026-10-21T12:13:52Z
    by: agent-S-0471
  - to: done
    at: 2026-10-21T12:21:16Z
    by: agent-S-0471
stream: S-0471
tags: [flai]
touches: [flai/internal/metrics/claims_bench_test.go]
usage:
  source: log
  seconds: 444
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 40
      output: 16861
      cache_read: 1861803
      cache_write: 66291
      cost: 1.1136
---
# T-1991 A benchmark and a reference test replay the holds over a generated history of this repository's size, and the profile names the fix

## Work

Reproduce I-0155 before changing anything, as the story's goal asks.

- In `flai/internal/metrics/claims_bench_test.go`, add a generator of a history at this repository's scale: about 500 stories and 2,000 tasks with creations and transitions spread over months, touches that overlap, folder touches narrowed by tasks, `after`, and a few shared paths in `Options.Shared`.
- Add `BenchmarkReplayHolds` over it, calling `replayHolds` directly, and `BenchmarkComputeClaims` through `Compute`.
- Copy today's `replayHolds` into the test file as `replayHoldsReference`, unchanged, so that the next task can prove the fixed replay gives the same `holdReplay`.
- Run the benchmark with `-cpuprofile` and read it. Record the time per replay, and the share of samples under `replayHolds`, `Holds.Claim`, `Holds.Overlaps`, `PathsOverlap`, and the garbage collector, in the narrative's `## Decisions`, with the fix chosen.
- The fix to weigh first: rebuild the holds only at an event that can change one, such as a story entering or leaving ready or in progress, a task of a ready or in-progress story created, closed, or cancelled, or a story named in `after` reaching done. Carry the previous instant's held set and counts by reason otherwise, and reuse each story's claim until its tasks' states change. `replayHolds` came with S-0205 on 2026-10-03, after the 10 ms measurement, and rebuilds every claim at every event time since the first item.

This task waits for none.

## Done when

- `flai test flai/internal/metrics` passes, with the generator, the two benchmarks, and the reference replay in place.
- The narrative's `## Decisions` names the profile's figures and the fix the next task builds.

## Notes
