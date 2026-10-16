---
id: T-1803
type: task
nature: experiment
title: Weigh what Postgres would add over SQLite and, where it might, run the same test in Docker
status: done
parent: S-0427
owner: alex
created: 2026-10-16T05:59:35Z
updated: 2026-10-16T07:27:08Z
transitions:
  - to: ready
    at: 2026-10-16T07:01:51Z
    by: agent-S-0427
  - to: in-progress
    at: 2026-10-16T07:01:52Z
    by: agent-S-0427
  - to: done
    at: 2026-10-16T07:27:08Z
    by: agent-S-0427
stream: S-0427
tags: [cli]
touches: [design/experiments/S-0427-sqlite-and-postgres-read-models.md]
after: [T-1802]
usage:
  source: log
  seconds: 1516
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 135
      output: 77508
      cache_read: 8539571
      cache_write: 281965
      cost: 4.8477
---
# T-1803 Weigh what Postgres would add over SQLite and, where it might, run the same test in Docker

## Work

Write under `## What was done` what a Postgres instance could give that SQLite cannot, for this host: one server that `flai serve`, each project's MCP server, the dashboard, and every story worktree share; concurrent writers; `LISTEN`/`NOTIFY` to push a change to the dashboard instead of polling; and what it costs, a daemon and a container to keep running.

If any of those might pay, run the same test:

- Start a throw-away `postgres:16-alpine` container with `docker run`, as `scripts/litellm.sh` runs its Postgres, on a free local port and a scratch volume.
- Load the same data with T-1801's loader pointed at Postgres, such as through `pgx`, and run the same queries.
- Time the full load, a cold connect and query, and a warm query, the median of three each, and check the figures against T-1800's reference.
- Measure the container's RAM and CPU with `docker stats`, and its volume's size.
- Remove the container and its volume.

If none would pay, say why, and leave Postgres out of the table with that reason.

Waits for T-1802: it reuses T-1801's loader and queries and compares with T-1802's overhead figures, and all three write the results document.

## Done when

- The results document says which advantages Postgres might have here, and why.
- Either its timings and its RAM, CPU, and disk figures sit in the same table as the baseline's and SQLite's, or the document says why it was not run.
- No container or volume of the test is left on the host.

## Notes
