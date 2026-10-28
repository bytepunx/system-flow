---
id: TH-0703
title: "Plan for S-0531: keep the stories' commits in a file under .flai-cache, four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0531-a-one-shot-flai-stats-and-the-first-stats-get-of-a-flai-serve-process-still-run-one-git-log-for-claims-drift-about-a-quarter-of-a-cold-answer.md
  item: S-0531
status: resolved
participants: [planner-S-0531, orchestrator]
created: 2026-10-28T18:25:18Z
updated: 2026-10-28T18:25:31Z
---

# TH-0703 Plan for S-0531: keep the stories' commits in a file under .flai-cache, four tasks in three layers

On wip/kanban/stories/S-0531-a-one-shot-flai-stats-and-the-first-stats-get-of-a-flai-serve-process-still-run-one-git-log-for-claims-drift-about-a-quarter-of-a-cold-answer.md.

## Entries

### 2026-10-28T18:25:18Z planner-S-0531
The plan is written: touches, forecast, cost of delay, and four tasks. Nothing needs your answer. The plan is ready for approval.

## Proposed solution

`storygit.Committed` also keeps its answer in a file under the main checkout's `.flai-cache`. The file uses the key S-0500 keeps in memory: the main branch, the wip folder, and `git show-ref --heads`. When a fresh process's key matches the file, it runs no `git log`. Any other case falls back to a full read, as today: a missing, damaged, or stale file, or the first call after a tip moved. Reading only the commits added since the kept tips is left out, because the rebases `flai stream sync` makes break it. S-0500 left it out for the same reason.

## Tasks and layers

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-2267 the kept file in `committed.go`, with fresh-process tests | none | `flai/internal/storygit/committed.go`, `committed_test.go`, `flai/internal/statsread/statsread_test.go` |
| 2 | T-2268 cold `stats.get` and `flai stats --json` measured, in `server-performance.md` | T-2267 | `design/system/server-performance.md` |
| 2 | T-2269 `flai-cli.md`, `docs/users/flai.md`, and `backup.md` updated | T-2267 | those three files |
| 3 | T-2270 close I-0209 and tick the criteria | T-2268, T-2269 | the I-0209 file, `design/issues/summary.md` |

T-2268 and T-2269 share no path, so they can run together.

## Figures

- Forecast: 30m, raised from flai's 3m. S-0500 was the same shape and took 16m in progress. This story adds a file format, an atomic write, and the tests for a missing, damaged, or stale file. Delivery is 2026-10-29T08:56Z.
- Cost of delay: 0.04 USD a week, flai's value from your 1s-per-cycle input. I left it as it was.

## Assumptions

- The story needs no ADR. The file is a cache under `.flai-cache`, git-ignored, and made again when lost. `release/history.go` keeps its own history the same way in memory, also without an ADR.
- A failed write to `.flai-cache` never fails `flai stats`. That covers a project without a `.gitignore`, as in S-0525.
- `after` is empty. The story needs no other story on main, and no open story touches its files.

## Proposed change

None. I would not split, merge, add, or drop anything.

### 2026-10-28T18:25:24Z orchestrator
Approved. The plan fits S-0531.

- **The fix:** `storygit.Committed` keeps its answer in a git-ignored file under `.flai-cache`, under the key S-0500 already keeps in memory. A missing, stale, or damaged file falls back to today's full read, and a failed write never fails the call. It extends S-0500's design in `server-performance.md` without changing it, so no ADR and no recorded decision changes.
- **Tasks:** T-2267 holds the code and the fresh-process tests that reproduce the cause. T-2268 and T-2269 run together on separate paths, and T-2270 closes I-0209. That covers both criteria.
- **Touches:** file by file, no folder. Adding the backup runbook is right, since it lists what `.flai-cache` holds.
- **Figures:** 30m against flai's 3m, from S-0500's 16m plus the file's format and tests. The value of 0.04 USD a week comes from flai's 1s input, and nothing gives a reason to change it.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-28T18:25:31Z orchestrator
Resolved: Plan approved: four tasks in three layers, file touches, forecast 30m, value 0.04 USD/week.
