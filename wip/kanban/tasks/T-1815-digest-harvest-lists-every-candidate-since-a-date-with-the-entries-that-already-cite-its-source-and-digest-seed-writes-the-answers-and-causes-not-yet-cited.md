---
id: T-1815
type: task
nature: feature
title: digest.Harvest lists every candidate since a date with the entries that already cite its source, and digest.Seed writes the answers and causes not yet cited
status: in-progress
parent: S-0429
owner: alex
created: 2026-10-16T06:04:04Z
updated: 2026-10-23T17:47:16Z
transitions:
  - to: ready
    at: 2026-10-23T17:37:51Z
    by: agent-S-0429
  - to: in-progress
    at: 2026-10-23T17:47:16Z
    by: agent-S-0429
stream: S-0429
tags: [cli]
touches: [flai/internal/digest/harvest.go, flai/internal/digest/harvest_test.go, flai/internal/digest/seed.go, flai/internal/digest/seed_test.go]
after: [T-1813, T-1814]
usage:
  source: log
  seconds: 111
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 18
      output: 99
      cache_read: 718624
      cache_write: 85825
      cost: 0.2714
---
# T-1815 digest.Harvest lists every candidate since a date with the entries that already cite its source, and digest.Seed writes the answers and causes not yet cited

## Work

- `harvest.go`: `Harvest(root, since, kinds)` joins the candidates of `Answers` over every thread, `Decisions`, and `Causes`, sorted by date, each with the entries that already cite its source (from `digest.List`). With a zero `since` it lists everything.
- `seed.go`: `Seed(root, since, dryRun)` writes, through `digest.New`, one entry per `answer`, `recommendation`, and `cause` candidate whose source no entry cites, and returns what it wrote and what it skipped and why. Decisions are never seeded.

Waits for the answer rule and the sources tasks, which it joins.

## Done when

`harvest_test.go` covers a mixed fixture sorted by date, `since`, `kinds`, and a cited source shown with its entry; `seed_test.go` covers a first run that writes, a second that writes nothing, a dry run that writes nothing, and decisions left out; `flai test flai/internal/digest` passes.

## Notes

The curator (S-0431) reads `Harvest` since its last run on every run, so its output must be stable for the same files.

Planner-S-0429, 2026-10-23: `threads.List` returns every thread, resolved ones included, and there is no threads archive, so `Harvest` reads `wip/threads` alone (610 threads, 591 of them resolved). Sort by date, then by source ID, so two candidates of one second keep one order. `digest.New` and `digest.List` are the names the draft assumed before S-0428 landed its package: take the exported names `flai/internal/digest` has once S-0428 is accepted, and the entry format from ADR-0191, which S-0428 writes.
