---
id: TH-1029
title: "Plan for S-0676: measure the fresh list, replace the kept items JSON, record and close I-0310, four tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0676-a-fresh-flai-list-decodes-10-5-mb-of-kept-items-json-so-it-takes-0-3-to-1-4-s.md
  item: S-0676
status: open
participants: [planner-S-0676]
created: 2026-11-10T10:31:59Z
updated: 2026-11-10T10:31:59Z
---

# TH-1029 Plan for S-0676: measure the fresh list, replace the kept items JSON, record and close I-0310, four tasks in four layers

On wip/kanban/stories/S-0676-a-fresh-flai-list-decodes-10-5-mb-of-kept-items-json-so-it-takes-0-3-to-1-4-s.md.

## Entries

### 2026-11-10T10:31:59Z planner-S-0676
S-0676 is planned: four tasks in four layers, each waiting for the one before.

## Tasks

| Layer | Task | Nature | Touches | Waits for |
|-------|------|--------|---------|-----------|
| 1 | T-2995 Measure where a fresh `Repo.List(true)` with kept files spends its time, and choose the remedy | research | `flai/internal/workitem/kept_bench_test.go` | none |
| 2 | T-2998 Restore the kept items files without decoding 10.5 MB of JSON, as T-2995 chose | improvement | `kept.go`, `kept_test.go`, `store.go`, `repo.go` under `flai/internal/workitem/` | T-2995: its figures choose the remedy |
| 3 | T-3001 Time a fresh list on this repository before and after, and record it | improvement | `design/system/server-performance.md`, `docs/operators/runbooks/backup.md` | T-2998: it measures that change |
| 4 | T-3006 Close I-0310 with the fix and the figures | improvement | the I-0310 file, `design/issues/summary.md` | T-3001: the reason quotes its figures |

## Figures

- Forecast 1h30m, against flai's 20m. S-0600 did the same kind of work and took 1.9 h of agent time.
- Delivery 2026-11-10T18:00:00Z.
- Cost of delay value 0.04 USD a week, from flai's input. It is unchanged.
- No folder touch and no `after`.

## Assumptions

- **Proposed remedy, for T-2995 to confirm or overturn.** Encode the kept items file with `encoding/gob` instead of JSON, and raise `keptVersion` so that today's JSON files are a miss. If T-2995's profile shows those phases matter, also restore the six item folders at once in `Repo.List`, and make `readAs` lazily. The story asks its agent to propose a remedy before building it, so T-2995 records the choice in the narrative's Decisions. Answer here if you want a different remedy.
- **"No longer occurs" means** decoding the kept files is no longer most of a fresh list's time. It is shown by T-2995's benchmark before and after, and by T-3001's in-process timing on this repository, measured the way S-0600 measured. A wall-time pass/fail test would be flaky under this host's load (I-0309), so the test that reproduces the cause is the benchmark plus format tests in `kept_test.go`.
- Not reaching 0.3 s median under a load average above 30 does not fail the story. T-3001 records CPU time beside wall time.

## Proposed apart, not planned here

`flai move` lists the archive only so that `Follow` sees an epic's archived stories (`flai/internal/workitem/transition.go`). The issue's third way forward, listing less there, would cut the cost of a move, not of a list. If you want it, I would make it a story of its own.

The prompt said you had written on a thread named "asked" since this story's last plan. No such thread exists, and S-0676 had no plan, tasks, or threads. So this is the story's first plan.
