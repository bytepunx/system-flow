---
id: T-1801
type: task
nature: experiment
title: Load the work items into SQLite deterministically and time the chart queries side by side with flai's API
status: done
parent: S-0427
owner: alex
created: 2026-10-16T05:59:18Z
updated: 2026-10-16T06:36:37Z
transitions:
  - to: ready
    at: 2026-10-16T06:09:55Z
    by: agent-S-0427
  - to: in-progress
    at: 2026-10-16T06:09:56Z
    by: agent-S-0427
  - to: done
    at: 2026-10-16T06:36:37Z
    by: agent-S-0427
stream: S-0427
tags: [cli]
touches: [design/experiments/S-0427-sqlite-and-postgres-read-models.md]
after: [T-1800]
usage:
  source: log
  seconds: 1601
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 148
      output: 85134
      cache_read: 9379784
      cache_write: 309708
      cost: 5.3247
---
# T-1801 Load the work items into SQLite deterministically and time the chart queries side by side with flai's API

## Work

Write a throw-away loader and query script in a scratch Go module under the ignored scratch folder (`bin/experiments/S-0427/`), outside `flai/go.mod`, so that flai gains no dependency. Use a pure Go SQLite driver such as `modernc.org/sqlite`, so the build needs no cgo; say which one and why.

- Load: walk `wip/kanban`, `wip/archive`, `wip/threads`, and `wip/messages`, parse each file's front matter, and write items, transitions, usage, and threads into tables. Load in sorted path order in one transaction, so that two loads from empty give the same rows; check it by hashing a sorted dump of both.
- Keep up to date: rescan by mtime and size, as `flai/internal/workitem/store.go` does, and reload only the changed files. Time a rescan with nothing changed and with one file changed.
- Query: the broad reads behind the charts, at least items per column, throughput per day, cycle time per story, time in state, and usage cost per model over time, over the same window the baseline used. Check each against the reference `flai stats --json` saved in T-1800, and say where a figure differs and why.
- Time the full load, a cold open and query (a one-shot agent call), and a warm query (a long-running `flai serve`), the median of three each.

Under `## What was done`, record the driver, the schema, the queries, and the commands, so the run can be repeated without the code. Put the timings in a table beside T-1800's baseline. Commit no code and no data.

Waits for T-1800: it needs the baseline and the reference output, and both write the results document.

## Done when

- Two loads from empty give the same dump hash.
- The results document has a side-by-side table of the baseline and the SQLite timings, cold and warm, for each query, and the time of a full load and of each rescan.
- Each query's figures are checked against the reference, and any difference is explained.

## Notes
