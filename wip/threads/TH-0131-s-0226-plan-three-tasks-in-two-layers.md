---
id: TH-0131
title: "S-0226 plan: three tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0226-the-orchestrator-s-cost-is-recorded-on-the-story-or-epic-each-decision-concerned.md
  item: S-0226
status: resolved
participants: [planner-S-0226, alex]
created: 2026-10-05T04:45:14Z
updated: 2026-10-05T04:57:23Z
---

# TH-0131 S-0226 plan: three tasks in two layers

On wip/kanban/stories/S-0226-the-orchestrator-s-cost-is-recorded-on-the-story-or-epic-each-decision-concerned.md.

## Entries

### 2026-10-05T04:45:14Z planner-S-0226
The plan for S-0226, by planner-S-0226.

Tasks and layers:

- Layer 1: T-0878. An ADR that extends ADR-0083 to the orchestrator and a project strategic total, plus `strategic-agents.md` and `work-hierarchy.md`. It decides the split, the edge cases, and where the project total comes from.
- Layer 2, run together: T-0879 and T-0880. Each waits for T-0878, and they share no path.
  - T-0879 charges an orchestrator activity evenly to the items it named, through `activity_log` and at the run's end. It touches usage, serve, mcpserver, `flai/cmd/activity.go`, and `flai-cli.md`.
  - T-0880 adds the project total per kind to `flai stats`. It touches metrics, `stats.go`, `metrics.md`, and the user docs.

Enrichment:

- Touches went from 8 to 15.
- Forecast is 35m, up from flai's 16m, because S-0225 took 52m.
- Cost of delay is 33.55 USD a week, from `flai cod` after the forecast change. The story has no inputs of its own.
- Topics `orchestration` and `planning` added.

The reasons are under `### Planning` in the Notes.

Assumptions:

1. Recommended: the project total is worked out when `flai stats` runs, from the activity entries that charged no item. Nothing new is stored, so the per-kind totals equal the activity document's totals by construction. If T-0878 decides to store the total instead, T-0880 gets an `after` on T-0879.
2. A duplicate or unknown item ID is left out of the split. An activity that names only unknown IDs goes to the project total. An activity that names a story and its epic gives the epic both shares, so each epic's sum still equals what was spent under it.
3. The charge is written in a way S-0227 (the analyzer) can reuse. T-0878 decides whether it covers every kind except the planner, or the orchestrator alone.
4. The dashboard is out of scope. The criteria name only `flai stats`.

I propose no split, merge, or drop.

### 2026-10-05T04:57:23Z alex
Resolved.
