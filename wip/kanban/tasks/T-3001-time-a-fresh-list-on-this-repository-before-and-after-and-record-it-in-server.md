---
id: T-3001
type: task
nature: improvement
title: Time a fresh list on this repository before and after, and record it in server-performance.md and the backup runbook
status: backlog
parent: S-0676
owner: alex
created: 2026-11-10T10:31:09Z
updated: 2026-11-10T10:31:09Z
transitions: []
stream: S-0676
tags: [flai, performance, docs]
touches: [design/system/server-performance.md, docs/operators/runbooks/backup.md]
after: [T-2998]
---
# T-3001 Time a fresh list on this repository before and after, and record it in server-performance.md and the backup runbook

## Work

Time a fresh `Repo.List(true)` with the kept files present, as S-0600 did (`design/system/server-performance.md` § Parses and lint findings kept between processes (S-0600) › Measured on this repository):

- a small program built with `go build -overlay` against `main` and against the story branch;
- a fresh process each round, ten rounds;
- in a scratch copy of the main checkout under `/tmp`;
- the median of wall and CPU time, with the one-minute load average.

Time a fresh `flai move` and `flai check --strict` the same way, five rounds each.

This task waits for T-2998, whose change it measures.

Record the figures and the fix in a section of their own in `design/system/server-performance.md`, and rewrite the S-0600 section's first "What is left" bullet to point at it. If T-2998 renamed the kept items file or changed its format, update its row in `docs/operators/runbooks/backup.md`.

## Done when

- `design/system/server-performance.md` gives the before and after figures, the host, the load, and the commits measured, with `updated` bumped.
- `docs/operators/runbooks/backup.md` names the kept items file as flai now writes it.
- `flai test` passes on both files.

## Notes
