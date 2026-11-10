---
id: S-0677
type: story
nature: improvement
title: "A whole flai check takes 3.3 to 4.9 s on a loaded host after S-0601, not under 2 s: the item list, workItems alone, and the digest rule are its critical path"
status: ready
owner: alex
created: 2026-11-10T08:34:54Z
updated: 2026-11-10T11:05:34Z
transitions:
  - to: ready
    at: 2026-11-10T10:57:48Z
    by: alex
tags: [flai, performance, check]
topics: [cli, server-side]
touches: [flai/internal/check/check.go, flai/internal/check/check_items.go, flai/internal/check/concurrent_test.go, flai/internal/check/items_test.go, flai/internal/check/check_records.go, flai/internal/check/digest_test.go, flai/internal/check/check_test.go, flai/internal/digest/digest.go, flai/internal/digest/digest_test.go, flai/internal/digest/kept.go, flai/internal/digest/kept_test.go, flai/internal/conflictmark/kept.go, flai/internal/conflictmark/kept_test.go, design/system/server-performance.md, design/system/flai-cli.md, docs/operators/runbooks/backup.md, design/issues/I-0312-a-whole-flai-check-takes-3-3-to-4-9-s-on-a-loaded-host-after-s-0601-not-under-2.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2s
    by: flai
    at: 2026-11-10T08:34:54Z
  value: 0.08
  by: planner-S-0677
  at: 2026-11-10T10:32:31Z
forecast:
  duration: 1h
  delivery: 2026-11-10T13:44:00Z
  basis: "Its own forecast of 1h; 3rd in the pull order with an in-progress limit of 3, behind S-0665, S-0667, S-0708, S-0666 and S-0688."
  by: flai
  at: 2026-11-10T11:05:34Z
finalized:
  by: alex
  at: 2026-11-10T10:57:45Z
---
# S-0677 A whole flai check takes 3.3 to 4.9 s on a loaded host after S-0601, not under 2 s: the item list, workItems alone, and the digest rule are its critical path

## Goal

This story remediates [I-0312](../../../design/issues/I-0312-a-whole-flai-check-takes-3-3-to-4-9-s-on-a-loaded-host-after-s-0601-not-under-2.md), "A whole flai check takes 3.3 to 4.9 s on a loaded host after S-0601, not under 2 s: the item list, workItems alone, and the digest rule are its critical path". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0312 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0312 is closed with `flai issue close I-0312 --reason` saying what fixed it

## Tasks
- T-3010 Run oneItem over the items on check.Run's pool, with only the byID fill first and alone
- T-3012 Keep the digest rule's per-entry results under .flai-cache, keyed by each entry's size and modification time
- T-3016 Keep the conflict-marker scan's per-file results under .flai-cache, keyed by each file's size and modification time
- T-3020 Time a warm flai check against main, record it in server-performance.md and the docs, and close I-0312

## Notes

Cost of delay inputs set by flai from I-0312. time_lost_per_cycle 2s: 2s per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-08T01:58:32Z, 2.3 days before this story; under one cycle counts as one).

### Planning

The proposed fix, from I-0312's instance and `design/system/server-performance.md` § "Cause 11":

- `workItems` keeps only the `byID` fill and a children index alone. Its `oneItem` calls go on the pool (T-3010). `oneItem` also scans every item for each item's children, which the index removes.
- The digest rule keeps each entry's results under `.flai-cache`, keyed by size and modification time, as S-0600 keeps the lint findings (T-3012). The `Known` and `Exists` lookups an entry needs run every check, since they depend on other files.
- `conflictMarkers`, the next slowest rule at 0.45 to 1.10 s, keeps its per-file results the same way (T-3016).
- T-3020 times the result and closes I-0312.

`Repo.List`'s decode of the kept items, 0.66 to 1.10 s of the path, is left to S-0676 (I-0310). This story reads nothing S-0676 writes, so it has no `after`.

Touches, by where each came from:

| Touch | From |
|-------|------|
| `flai/internal/check/check.go`, `check_items.go` | design: `server-performance.md` § "Cause 11" names `check.Run` and `workItems`; the explorer found `oneItem`'s children scan |
| `flai/internal/check/concurrent_test.go`, `items_test.go` | layout: the tests of the run's order and of a run on given items |
| `flai/internal/check/check_records.go`, `digest_test.go`, `check_test.go` | layout: the digest and conflict-marker rules and their tests |
| `flai/internal/digest/digest.go` | layout: `Validate` and `EnvFor` |
| `flai/internal/digest/digest_test.go` | co-change: `flai touches suggest` gives it 3 of 4 commits with the seeds |
| `flai/internal/digest/kept.go`, `kept_test.go`, `flai/internal/conflictmark/kept.go`, `kept_test.go` | layout: new files, beside `flai/internal/mdlint/kept.go`'s pattern |
| `design/system/server-performance.md`, `design/system/flai-cli.md` | design: where the S-0600 and S-0601 caches and timings are recorded |
| `docs/operators/runbooks/backup.md` | design: the one table of `.flai-cache` files |
| `design/issues/I-0312-…`, `design/issues/summary.md` | declared by the criteria: closing I-0312 |

The story declared no touches before this plan. All its touches name files; none is a folder.

Figures:

- Forecast 1h, not flai's 33m. flai sized the story by its 17 touches at 103 s per unit. The story adds two kept-cache files, each with tests, as S-0600 did for one, plus timed runs on a loaded host.
- Cost of delay value 0.08 USD a week, as `flai cod` works it out from the inputs flai set from I-0312: 2 s per 168 h cycle at 150 USD an hour. The value stands. The input counts I-0312's one report, not every check. See the plan's thread.
