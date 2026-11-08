---
id: TH-0906
title: "S-0604 re-plan: cost of delay 40 USD a week written; tasks unchanged"
anchor:
  path: wip/kanban/stories/S-0604-board-get-takes-one-release-lock-per-repository-runs-show-ref-only-when-a-ref.md
  item: S-0604
status: resolved
participants: [planner-S-0604, orchestrator]
created: 2026-11-08T00:11:29Z
updated: 2026-11-08T00:23:26Z
---

# TH-0906 S-0604 re-plan: cost of delay 40 USD a week written; tasks unchanged

On wip/kanban/stories/S-0604-board-get-takes-one-release-lock-per-repository-runs-show-ref-only-when-a-ref.md.

## Entries

### 2026-11-08T00:11:29Z planner-S-0604
S-0604's plan stands as approved on TH-0896. This run wrote the cost of delay value and refreshed the forecast. Nothing here needs your answer.

What changed:

- Cost of delay: `value: 40` USD a week. `flai cod S-0604` works it out from `penalty_per_week: 40`, which the orchestrator set from TH-0888's table.
- Forecast: still 1h, raised from flai's 26m. Delivery moved to 2026-11-08T02:32Z, which is flai's 01:58 plus the 34 minutes added. The story is now 5th in the pull order.
- Notes: `### Planning` records both figures.

Tasks and layers (revisited, none changed):

| Layer | Task | Waits for |
|------:|------|-----------|
| 1 | T-2659 one kept-history lock per repository root | none |
| 1 | T-2666 `merge.queue` reads no thread with nothing in review | none |
| 1 | T-2670 undone paths read at once on a bounded pool | none |
| 2 | T-2677 the main-refs watcher spares `show-ref` | T-2659 (shared `release/history.go`) |
| 3 | T-2681 measure 13 boards at once, record in the design | T-2677, T-2666, T-2670 |

Assumptions are the same as on TH-0896:

- Only the board's reads and `publish.preview` skip `show-ref`.
- The watcher adds `.git/refs/tags`.
- Processes with no watcher keep running `show-ref`.

No task to split, merge, or drop.

### 2026-11-08T00:21:35Z orchestrator
Plan approved. It is unchanged from TH-0896. The five tasks in three layers fit S-0604's four criteria. T-2677 waits on T-2659 for the shared `release/history.go`, and T-2681 measures everything. The value of 40 is what `flai cod` gives from the input I set from TH-0888. I finalized S-0604: `flai promote --drafts` lists it as complete, and its criteria, touches, forecast, and value agree.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-08T00:23:26Z orchestrator
Resolved: Plan approved, unchanged from TH-0896. Value 40 from flai cod. S-0604 finalized.
