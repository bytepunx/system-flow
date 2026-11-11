---
id: S-0676
type: story
nature: improvement
title: A fresh flai list decodes 10.5 MB of kept items JSON, so it takes 0.3 to 1.4 s, not the 0.3 s S-0600's title names
status: in-progress
owner: alex
created: 2026-11-10T08:34:53Z
updated: 2026-11-11T01:41:18Z
transitions:
  - to: ready
    at: 2026-11-10T10:57:41Z
    by: alex
  - to: in-progress
    at: 2026-11-11T01:41:18Z
    by: agent-S-0676
tags: [flai, performance]
topics: [back-end, cli]
touches: [flai/internal/workitem/kept.go, flai/internal/workitem/kept_test.go, flai/internal/workitem/kept_bench_test.go, flai/internal/workitem/store.go, flai/internal/workitem/repo.go, design/system/server-performance.md, docs/operators/runbooks/backup.md, design/issues/I-0310-a-fresh-flai-list-decodes-10-5-mb-of-kept-items-json-so-it-takes-0-3-to-1-4-s.md, design/issues/summary.md]
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
      seconds: 16
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 67
          cache_read: 270231
          cache_write: 10640
          cost: 0.0768
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1s
    by: flai
    at: 2026-11-10T08:34:53Z
  value: 0.04
  by: planner-S-0676
  at: 2026-11-10T10:31:42Z
forecast:
  duration: 1h30m
  delivery: 2026-11-11T03:35:00Z
  basis: "Its own forecast of 1h30m; 1st in the pull order with an in-progress limit of 3, with nothing ahead of it."
  by: flai
  at: 2026-11-11T01:40:45Z
finalized:
  by: alex
  at: 2026-11-10T10:57:37Z
---
# S-0676 A fresh flai list decodes 10.5 MB of kept items JSON, so it takes 0.3 to 1.4 s, not the 0.3 s S-0600's title names

## Goal

This story remediates [I-0310](../../../design/issues/I-0310-a-fresh-flai-list-decodes-10-5-mb-of-kept-items-json-so-it-takes-0-3-to-1-4-s.md), "A fresh flai list decodes 10.5 MB of kept items JSON, so it takes 0.3 to 1.4 s, not the 0.3 s S-0600's title names". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0310 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0310 is closed with `flai issue close I-0310 --reason` saying what fixed it

## Tasks
- T-2995 Measure where a fresh Repo.List(true) with kept items files spends its time, and choose the remedy
- T-2998 Restore the kept items files without decoding 10.5 MB of JSON, as T-2995 chose
- T-3001 Time a fresh list on this repository before and after, and record it in server-performance.md and the backup runbook
- T-3006 Close I-0310 with the fix and the figures T-3001 recorded

## Notes

Cost of delay inputs set by flai from I-0310. time_lost_per_cycle 1s: 1s per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-08T01:17:08Z, 2.3 days before this story; under one cycle counts as one).

### Planning

Touches, file by file. The story declared none, so `flai touches suggest` had no path to start from. Run again with these paths, it listed only files most stories change, such as `design/system/flai-cli.md` at 7 percent and `docs/users/flai.md` at 6 percent. No command or flag changes here, so none of them is kept.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/workitem/kept.go` | layout | Writes and reads `.flai-cache/items-<folder>.json` (`readKeptItems`, `restoreItem`, `writeKeptItems`): the decode I-0310 names |
| `flai/internal/workitem/kept_test.go` | layout | The kept file's tests: the new format, and the old one as a miss |
| `flai/internal/workitem/kept_bench_test.go` | layout | New: the benchmark that reproduces the fresh list's cost (T-2995) |
| `flai/internal/workitem/store.go` | layout | `lookAll`, `lookup`, `parseAll`, and `keep` restore kept entries and compare `readAs` |
| `flai/internal/workitem/repo.go` | layout | `Repo.List` walks the six item folders one after another |
| `design/system/server-performance.md` | design | § Parses and lint findings kept between processes (S-0600) records the 0.94 s median this story remedies |
| `docs/operators/runbooks/backup.md` | design | Its row names `.flai-cache/items-<folder>.json`, which a new format may rename |
| `design/issues/I-0310-...md`, `design/issues/summary.md` | declared by criterion 2 | `flai issue close` writes both |

No folder touch. No `after`: the story reads nothing an open story writes.

Forecast 1h30m, against flai's 20m. flai's figure is 107 s per unit of size times size 11, taken from done improvement stories in the medium band. S-0600 built the kept files and measured them the same way, and took 6,691 s of agent time, 1.9 h. This story has fewer criteria than S-0600 but the same in-process and fresh-process timing on a loaded host. 1h30m is the nearer figure. Delivery 2026-11-10T18:00:00Z: `flai forecast` puts it 14th in the pull order behind 18 stories with an in-progress limit of 5, which is about four rounds ahead of it.

Cost of delay value 0.04 USD a week, as `flai cod` works it out from the input flai set from I-0310: 1 s lost per 168 h cycle at 150 USD an hour. It stands, since I-0310 has one instance. The value counts only one second per fresh list; every one-shot `flai` and every `flai move` lists the archive.
