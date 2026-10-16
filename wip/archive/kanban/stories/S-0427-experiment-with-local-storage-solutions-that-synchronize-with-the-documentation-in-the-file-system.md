---
id: S-0427
type: story
nature: experiment
title: Experiment with local storage solutions that synchronize with the documentation in the file system
status: done
owner: alex
created: 2026-10-16T05:55:07Z
updated: 2026-10-16T07:30:22Z
transitions:
  - to: ready
    at: 2026-10-16T05:55:09Z
    by: alex
  - to: in-progress
    at: 2026-10-16T05:59:11Z
    by: agent-S-0427
  - to: review
    at: 2026-10-16T07:29:32Z
    by: agent-S-0427
  - to: done
    at: 2026-10-16T07:30:22Z
    by: orchestrator
tags: []
topics: [server-side]
touches: [design/experiments/S-0427-sqlite-and-postgres-read-models.md, design/issues/I-0155-the-chart-data-takes-about-4-s-warm-against-10-ms-measured-on-2026-09-29-almost-all-of-it-in-metrics-replayholds.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5447
  priced_by: harness
  turns:
    - day: 2026-10-16
      ceremony: 1
      hand_edits: 2
      work: 59
  models:
    - model: claude-opus-5-5
      input: 556
      output: 318777
      cache_read: 35121967
      cache_write: 1159681
      cost: 19.9379
  strategic:
    - kind: planner
      seconds: 3287
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 160
          output: 21546
          cache_read: 6592156
          cache_write: 362068
          cost: 0.4427
        - model: claude-opus-5-5
          input: 2942
          output: 78975
          cache_read: 523318036
          cache_write: 841375
          cost: 238.9423
    - kind: orchestrator
      seconds: 4846
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 158
          output: 2643
          cache_read: 37617568
          cache_write: 57586
          cost: 9.4452
        - model: claude-sonnet-5-5
          input: 8
          output: 64
          cache_read: 89116
          cache_write: 59918
          cost: 0.122
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h
    by: planner-S-0427
    at: 2026-10-16T06:49:35Z
  value: 300
  by: planner-S-0427
  at: 2026-10-16T06:49:43Z
forecast:
  duration: 3h
  delivery: 2026-10-16T09:00:00Z
  basis: "flai forecast gives 7m (94 s a unit times size 4), raised to 3h for a scratch SQLite loader checked against metrics.Compute, a Postgres run in Docker, and a five-layer chain from 06:00Z."
  by: planner-S-0427
  at: 2026-10-16T06:00:42Z
---
# S-0427 Experiment with local storage solutions that synchronize with the documentation in the file system

## Goal

As of now, flai relies heavily on reading front-matter and documents from across the repository in order to present critical information for both the operator and agents.

As projects grow larger, this overhead starts to show up in both agents and tools like the board and chart.

Lets look at what it might take to leverage a database instance that gets deterministically populated from the file system and kept up-to-date with the changes on the local system such that the agents and flai can perform more demanding read operations over a SQLite connection vs. reading from hundreds or thousands of files.

Don't worry about the overhead of introducing a pass-through storage adapter in addition to the file operations, focus primarily on the potential time savings and look at how this might reduce both operator wait time and agent overhead.

## Acceptance criteria
- [x] Side-by-side comparisons of chart data (load a sql lite db then make comparative calls from flai's api vs. a throw-away script that reads the same broad information from Sqlite)
- [x] Exploration of storage, RAM, and CPU overhead that might be added with a SQLite instance
- [x] If there might be additional advantages to a postgres instance, spin one up via docker and perform the same kind of test and present the data alongside the other comparisons

## Tasks
- T-1800 Measure the file-read baseline for chart data on today's corpus and start the results document
- T-1801 Load the work items into SQLite deterministically and time the chart queries side by side with flai's API
- T-1802 Measure the storage, RAM, and CPU a SQLite read model adds
- T-1803 Weigh what Postgres would add over SQLite and, where it might, run the same test in Docker
- T-1804 Write the side-by-side results and the recommendation for the read model
- T-1805 Measure the storage, RAM, and CPU a SQLite read model adds, on today's corpus and a synthetic ten-times corpus
- T-1806 Run the same load and chart queries against a throw-away Postgres container and write the results and recommendation

## Notes

### Planning

Written by the planner (planner-S-0427) on 2026-10-16. The cost of delay input was settled on TH-0474; the plan's summary is TH-0478.

Touches:

- `design/experiments/S-0427-sqlite-and-postgres-read-models.md`, from design: an `experiment` story records its results in one document under `design/experiments/` ([ADR-0066](../../../design/adrs/0066-an-experiment-story-is-accepted-like-any-other-and-records-its-results-in-a.md)), with the sections of `design/experiments/template.md`. The story declared no touches, and `flai touches suggest` found no co-change for the path.
- The two `design/issues/` paths were added by the story's agent while it worked; they are kept.
- No folder touch is kept. The throw-away loader, queries, and data live in the git-ignored scratch folder `bin/experiments/S-0427/` and are not committed, so flai gains no dependency in `flai/go.mod` and no new folder needs an ADR. `design/experiments/README.md` has no index, so it needs no edit.
- What is measured, not changed: the chart data's path, `stats.get` in `flai/internal/hostapi/reads.go`, `statsread.Read` in `flai/internal/statsread/statsread.go`, `metrics.Compute` in `flai/internal/metrics/metrics.go`, and the file reads of `flai/internal/workitem/store.go`; the method of `design/system/server-performance.md`.

Figures:

- Forecast: `flai forecast` gives 7m, the median of 94 s a unit of size over 235 done stories times size 4. Raised to 3h: the loader must parse every work item, transition, and usage block, its queries must be checked against `metrics.Compute`'s figures, and Postgres runs in Docker, as a chain of five layers. Delivery counts from the start at 06:00Z.
- Cost of delay: 300 USD a week, as `flai cod` works it out: 2h of time lost per 168h cycle at 150 USD an hour. The input, `time_lost_per_cycle: 2h`, is the operator's, confirmed on TH-0474 on 2026-10-16 and set by the planner on that word. The value stands unadjusted: the experiment has no revenue or penalty to add.
- Topics: `server-side` added, so the pack carries `design/system/server-performance.md`, the baseline method and the 2026-09-29 figures.

Tasks: the planner drafted T-1800 to T-1804, one layer each. The story's agent wrote T-1805 and T-1806 after it started the story, then cancelled them and kept the planner's chain.

### Accepted by the orchestrator

- Verified: ce03beae4fd125112fbc828f9a735e5c9536b444
- At: 2026-10-16T07:30:22Z

Verdict: all three criteria met at ce03bea; flai verify passed at that commit against main 6e0f297; Postgres was run in Docker; no convention breach.
- 1: design/experiments/S-0427-sqlite-and-postgres-read-models.md
- 2: design/experiments/S-0427-sqlite-and-postgres-read-models.md
- 3: design/experiments/S-0427-sqlite-and-postgres-read-models.md
