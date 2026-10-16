---
id: T-1805
type: task
nature: feature
title: Measure the storage, RAM, and CPU a SQLite read model adds, on today's corpus and a synthetic ten-times corpus
status: cancelled
parent: S-0427
owner: alex
created: 2026-10-16T06:00:08Z
updated: 2026-10-16T06:36:57Z
transitions:
  - to: cancelled
    at: 2026-10-16T06:36:57Z
    by: agent-S-0427
stream: S-0427
tags: []
touches: [design/experiments/S-0427-sqlite-and-postgres-read-models.md]
after: [T-1801]
---
# T-1805 Measure the storage, RAM, and CPU a SQLite read model adds, on today's corpus and a synthetic ten-times corpus

## Work

Using T-1801's loader and queries in the ignored scratch folder `bin/experiments/S-0427/`, measure what a SQLite read model costs beside flai's file reads:

- Storage: the database file's size after a full load, with its WAL and indexes, against the corpus's markdown size from T-1800.
- RAM: max RSS of a full load, a rescan, a cold one-shot query, and a warm query loop, with `/usr/bin/time -v`, beside the baseline's RSS for `flai stats --json` and `stats.get`.
- CPU: user and system CPU time of the same runs.
- Growth: generate a synthetic corpus ten times today's, copying the work items with fresh IDs into a scratch folder outside the repository, and repeat the baseline's `flai stats --json` against it (as a copy of the repository with the copied items) and the SQLite load, rescan, and queries, so the results show how each grows.

Record the commands and a table of the figures under `## What was done`. Commit no code and no data.

Waits for T-1801: it reuses its loader and queries, and both write the results document.

## Done when

- The results document has a table of storage, max RSS, and CPU time for the baseline and for SQLite, on today's corpus and the ten-times corpus, each with the command that measured it.
- It says how the baseline and SQLite timings grow from today's corpus to the ten-times corpus.

## Notes
- 2026-10-16T06:36:57Z: moved to cancelled: duplicates the planner's T-1802, written in parallel; its ten-times corpus is folded into T-1802
