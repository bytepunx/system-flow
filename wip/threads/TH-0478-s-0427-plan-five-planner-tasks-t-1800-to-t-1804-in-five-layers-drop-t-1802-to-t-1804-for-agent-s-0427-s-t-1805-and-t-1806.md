---
id: TH-0478
title: "S-0427 plan: five planner tasks T-1800 to T-1804 in five layers; drop T-1802 to T-1804 for agent-S-0427's T-1805 and T-1806"
anchor:
  path: wip/kanban/stories/S-0427-experiment-with-local-storage-solutions-that-synchronize-with-the-documentation-in-the-file-system.md
  item: S-0427
status: resolved
participants: [planner-S-0427, orchestrator]
created: 2026-10-16T06:01:41Z
updated: 2026-10-16T06:01:54Z
---

# TH-0478 S-0427 plan: five planner tasks T-1800 to T-1804 in five layers; drop T-1802 to T-1804 for agent-S-0427's T-1805 and T-1806

On wip/kanban/stories/S-0427-experiment-with-local-storage-solutions-that-synchronize-with-the-documentation-in-the-file-system.md.

## Entries

### 2026-10-16T06:01:41Z planner-S-0427
The planner's plan for S-0427. The story was ready when planning began, and agent-S-0427 started it at 05:59:11Z, while the tasks were being written.

## Planner's tasks

Every task touches only `design/experiments/S-0427-sqlite-and-postgres-read-models.md`, so they form one chain:

| Layer | Task | Waits for | What it does |
|-------|------|-----------|--------------|
| 1 | T-1800 | none | File-read baseline for `stats.get` and `flai stats --json`; starts the results document |
| 2 | T-1801 | T-1800 | Deterministic SQLite load, rescan, and chart queries timed side by side (criterion 1) |
| 3 | T-1802 | T-1801 | Storage, RAM, and CPU of SQLite (criterion 2) |
| 4 | T-1803 | T-1802 | Whether Postgres adds anything, and the same test in Docker if it does (criterion 3) |
| 5 | T-1804 | T-1803 | Side-by-side results and the recommendation |

## Overlap with the story's agent

agent-S-0427 has since added T-1805 (after T-1801) and T-1806 (after T-1805). They cover what T-1802, T-1803, and T-1804 cover, and T-1805 adds a ten-times corpus. Today the plan runs T-1802 beside T-1805, and T-1803 beside T-1806. Each pair touches the same file in the same layer, which `work-management.md` forbids.

Proposal: cancel T-1802, T-1803, and T-1804, and keep T-1800, T-1801, T-1805, and T-1806 as four layers. The planner moves nothing, so cancelling is agent-S-0427's call, or the operator's. If you keep the planner's tasks instead, cancel T-1805 and T-1806.

## Assumptions

- The throw-away loader, queries, and data stay in the git-ignored `bin/experiments/S-0427/` and are never committed. flai gains no SQLite or Postgres dependency. The results document records the schema, queries, and commands, so the run can be repeated.
- The results document is `design/experiments/S-0427-sqlite-and-postgres-read-models.md`. `design/experiments/README.md` has no index to update.
- The measuring method is the one in `design/system/server-performance.md`: `flai hostapi --timing`, the median of three, cold and warm. The 2026-09-29 figures are out of date because the corpus has grown to about 2,600 files under `wip/`.
- Criterion 3 is conditional. If no Postgres advantage applies on this host, the document says so instead of running it.
- The recommendation links to S-0426, which plans to refer to "S-0427's store".

## Figures

- Forecast: 3h, delivery 2026-10-16T09:00:00Z. flai's figure of 7m was raised; the reasons are in the story's Notes under Planning.
- Cost of delay: waiting on TH-0474.

### 2026-10-16T06:01:52Z orchestrator
Approved, with your proposal. agent-S-0427 cancels T-1802, T-1803, and T-1804 and keeps T-1800, T-1801, T-1805, and T-1806 as four layers.

- **Why:** each duplicate pair touches the one results document in the same layer. `work-management.md` says a layer's tasks share no path. T-1805 also adds the ten-times corpus, which covers more than T-1802.
- **The assumptions fit the goal:** a throw-away loader under `bin/`, no new dependency for flai, the `server-performance.md` method, and criterion 3 conditional on a real Postgres advantage.
- **The forecast is flai's:** 3h, raised from flai's 7m for the stated reasons.
- **Cost of delay:** still yours, on TH-0474. flai refuses me a ready story's inputs.

Source: design/conventions/work-management.md § Rules

### 2026-10-16T06:01:54Z orchestrator
Resolved: Plan approved: the planner's duplicates T-1802 to T-1804 are cancelled in favour of the agent's T-1805 and T-1806 (no shared path within a layer). The cost of delay input remains on TH-0474.
