---
id: T-1806
type: task
nature: experiment
title: Run the same load and chart queries against a throw-away Postgres container and write the results and recommendation
status: backlog
parent: S-0427
owner: alex
created: 2026-10-16T06:00:31Z
updated: 2026-10-16T06:00:40Z
transitions: []
stream: S-0427
tags: []
touches: [design/experiments/S-0427-sqlite-and-postgres-read-models.md]
after: [T-1805]
---
# T-1806 Run the same load and chart queries against a throw-away Postgres container and write the results and recommendation

## Work

Start a throw-away Postgres container of its own, such as `docker run -d --rm --name s0427-pg -p 127.0.0.1:55432:5432 -e POSTGRES_PASSWORD=local postgres:16-alpine`, never touching another container on the host. Point T-1801's loader and queries at it with a pure Go driver such as `github.com/jackc/pgx/v5`, in the same scratch module under `bin/experiments/S-0427/`, and measure what T-1801 and T-1805 measured for SQLite: full load, rescan with nothing and one file changed, cold and warm queries, the median of three, and the container's storage (`pg_database_size`), RAM, and CPU (`docker stats --no-stream`, idle and under the query loop), on today's corpus and the ten-times corpus. Note what Postgres offers beyond SQLite for this use, such as concurrent writers or `LISTEN/NOTIFY`, and what it costs the operator: a container to run, a port, a credential. Stop and remove the container when done.

Then write `## Results` against the success measure and `## Recommendation` (adopt, adapt, or drop, and why), with the Postgres figures beside the baseline and SQLite in the tables, and set the document's `status`.

Waits for T-1805: it reuses the loader, the queries, and the ten-times corpus, and all three write the results document.

## Done when

- The results document's tables have a Postgres column beside the baseline and SQLite, for timings and for storage, RAM, and CPU.
- `## Results` and `## Recommendation` are written, the recommendation naming what would carry it on.
- The container is stopped and removed.

## Notes
