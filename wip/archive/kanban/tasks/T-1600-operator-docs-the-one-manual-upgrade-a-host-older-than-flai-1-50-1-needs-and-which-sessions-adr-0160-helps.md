---
id: T-1600
type: task
nature: feature
title: "Operator docs: the one manual upgrade a host older than flai 1.50.1 needs, and which sessions ADR-0160 helps"
status: done
parent: S-0391
owner: alex
created: 2026-10-15T05:09:36Z
updated: 2026-10-15T05:21:58Z
transitions:
  - to: ready
    at: 2026-10-15T05:09:58Z
    by: agent-S-0391
  - to: in-progress
    at: 2026-10-15T05:09:58Z
    by: agent-S-0391
  - to: done
    at: 2026-10-15T05:21:58Z
    by: agent-S-0391
stream: S-0391
tags: []
touches: [docs/operators/index.md, docs/operators/runbooks/update.md, docs/operators/runbooks/migrate.md]
usage:
  source: log
  seconds: 720
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 5299
      cache_read: 1304778
      cache_write: 83477
      cost: 0.9642
---

# T-1600 Operator docs: the one manual upgrade a host older than flai 1.50.1 needs, and which sessions ADR-0160 helps

## Work

I-0149: the host runs flai 1.48.0, older than ADR-0158 (1.50.1) and ADR-0160 (1.51.1), so nothing on it upgrades or asks, and sessions whose MCP server predates 1.51.1 keep refusing `task_done` even after an upgrade. Document the one manual `flai host upgrade` a host older than 1.50.1 needs in `docs/operators/runbooks/update.md` and as a row in `docs/operators/runbooks/migrate.md`. Correct `docs/operators/index.md`'s claim that sessions started before an upgrade go on without a restart: only those whose MCP server is 1.51.1 or newer. Say there that the thread waits for the release to be published (T-1599). Waits for nothing: it shares no path with T-1599.

## Done when

- `update.md` and `migrate.md` name the one manual upgrade for hosts older than 1.50.1.
- `index.md` says which sessions ADR-0160 helps, and that the thread waits for a published release.
- `flai test` on the changed paths passes (markdown lint).

## Notes
