---
id: S-0427
type: story
nature: experiment
title: Experiment with local storage solutions that synchronize with the documentation in the file system
status: in-progress
owner: alex
created: 2026-10-16T05:55:07Z
updated: 2026-10-16T06:01:28Z
transitions:
  - to: ready
    at: 2026-10-16T05:55:09Z
    by: alex
  - to: in-progress
    at: 2026-10-16T05:59:11Z
    by: agent-S-0427
tags: []
topics: [server-side]
touches: [design/experiments/S-0427-sqlite-and-postgres-read-models.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 313
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 66
          output: 1188
          cache_read: 31380205
          cache_write: 16267
          cost: 7.8708
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
- [ ] Side-by-side comparisons of chart data (load a sql lite db then make comparative calls from flai's api vs. a throw-away script that reads the same broad information from Sqlite)
- [ ] Exploration of storage, RAM, and CPU overhead that might be added with a SQLite instance
- [ ] If there might be additional advantages to a postgres instance, spin one up via docker and perform the same kind of test and present the data alongside the other comparisons

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

Written by the planner (planner-S-0427) on 2026-10-16. The cost of delay input is asked on TH-0474; the plan's summary is in its own thread on this story.

Touches:

- `design/experiments/S-0427-sqlite-and-postgres-read-models.md`, from design: an `experiment` story records its results in one document under `design/experiments/` ([ADR-0066](../../../design/adrs/0066-an-experiment-story-is-accepted-like-any-other-and-records-its-results-in-a.md)), with the sections of `design/experiments/template.md`. The story declared no touches, and `flai touches suggest` found no co-change for the path.
- No folder touch is kept. The throw-away loader, queries, and data live in the git-ignored scratch folder `bin/experiments/S-0427/` and are not committed, so flai gains no dependency in `flai/go.mod` and no new folder needs an ADR. `design/experiments/README.md` has no index, so it needs no edit.
- What is measured, not changed: the chart data's path, `stats.get` in `flai/internal/hostapi/reads.go`, `statsread.Read` in `flai/internal/statsread/statsread.go`, `metrics.Compute` in `flai/internal/metrics/metrics.go`, and the file reads of `flai/internal/workitem/store.go`; the method of `design/system/server-performance.md`.

Figures:

- Forecast: `flai forecast` gives 7m, the median of 94 s a unit of size over 235 done stories times size 4. Raised to 3h: the loader must parse every work item, transition, and usage block, its queries must be checked against `metrics.Compute`'s figures, and Postgres runs in Docker, as a chain of five layers. Delivery counts from the start at 06:00Z.
- Cost of delay: no value yet. The story has no inputs and no epic, so `flai cod` cannot work one out; TH-0474 asks the operator, recommending `time_lost_per_cycle: 2h`.
- Topics: `server-side` added, so the pack carries `design/system/server-performance.md`, the baseline method and the 2026-09-29 figures.

Tasks: the planner drafted T-1800 to T-1804. T-1805 and T-1806 are the story's agent's, written after it started the story; they cover what T-1802, T-1803, and T-1804 cover, adding a ten-times corpus.
