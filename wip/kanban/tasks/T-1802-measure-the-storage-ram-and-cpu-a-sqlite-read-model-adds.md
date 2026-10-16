---
id: T-1802
type: task
nature: experiment
title: Measure the storage, RAM, and CPU a SQLite read model adds
status: backlog
parent: S-0427
owner: alex
created: 2026-10-16T05:59:26Z
updated: 2026-10-16T05:59:26Z
transitions: []
stream: S-0427
tags: [cli]
touches: [design/experiments/S-0427-sqlite-and-postgres-read-models.md]
after: [T-1801]
---
# T-1802 Measure the storage, RAM, and CPU a SQLite read model adds

## Work

With T-1801's loader and queries, measure what a SQLite read model would add to the host, and record it under `## What was done` and in the results table:

- Storage: the database file and its WAL after a full load, per thousand items, and how they grow with the corpus.
- RAM: max RSS of the full load, of a cold open and query, and of a process holding the database open warm, against the baseline's RSS from T-1800, with `/usr/bin/time -v`.
- CPU: CPU time of a full load, of a rescan with nothing changed, and of a rescan with one file changed, and what a rescan every second or on each request would cost a long-running `flai serve` per hour.
- Where the database would live, such as under `.flai-cache/`, one per checkout or worktree, and what a story worktree or a fresh clone pays to build it.

Waits for T-1801: it measures T-1801's loader, and both write the results document.

## Done when

- The results document gives the storage, RAM, and CPU figures above, each with the command that measured it, beside the baseline's.
- It says where the database would live and what building it costs a fresh checkout.

## Notes
