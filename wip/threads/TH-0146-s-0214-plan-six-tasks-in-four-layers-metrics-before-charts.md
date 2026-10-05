---
id: TH-0146
title: "S-0214 plan: six tasks in four layers, metrics before charts"
anchor:
  path: wip/kanban/stories/S-0214-charts-show-parallelism-holds-and-touches-drift.md
  item: S-0214
status: resolved
participants: [planner-S-0214, alex]
created: 2026-10-05T05:46:19Z
updated: 2026-10-05T05:46:33Z
---

# TH-0146 S-0214 plan: six tasks in four layers, metrics before charts

On wip/kanban/stories/S-0214-charts-show-parallelism-holds-and-touches-drift.md.

## Entries

### 2026-10-05T05:46:19Z planner-S-0214
Plan for S-0214, by planner-S-0214.

Tasks and layers:

- Layer 1: T-0924, which defines in `metrics.md` held stories per day, held hours by reason per week, and the weekly exact-touches share, with an ADR refining ADR-0081.
- Layer 2: T-0929, which computes them in `flai/internal/metrics`, and T-0936, which adds `claims` to the dashboard's `Report` and the three chart builders in `flaiover/src/lib/viz`. These two run together: both wait only for T-0924, and their paths are apart.
- Layer 3: T-0933, which adds them to `flai stats` printing and help and to `docs/users/flai.md` and `flai-reference.md` (after T-0929). T-0943 runs beside it: it puts the three charts under Planning on the Charts page and compares them by hand with `flai stats --json` (after T-0929 and T-0936).
- Layer 4: T-0946, which updates `flaiover-dashboard.md` and `docs/users/flaiover.md` (after T-0943).

Assumptions:

1. The criteria need figures `flai stats` lacks: held stories per day, held time by reason and by week, and the weekly exact share. ADR-0081 says the charts compute nothing of their own, so I put them in flai behind a new ADR rather than deriving them in the dashboard.
2. A hold with more than one reason counts under its first, so that the hold-time stack sums to the held time. The alternative is to count it under each reason, which double-counts.
3. "Empty claim" is the `no-touches` hold.
4. The Planning group is on the Charts page, not in `sitemenu.ts`. Whichever of S-0212, S-0213, and S-0214 lands first creates `PLANNING_KINDS`, and the others reuse it. S-0214 waits only for S-0205, so it may be first.
5. `flaiover/src/lib/charts` and `flaiover/src/lib/sitemenu.ts` stay declared, but no task changes them.

Figures: forecast 1h45m, kept over flai's 31m because the Go metrics and the ADR are extra work. Delivery 2026-10-05T15:47Z. Cost of delay 106.42 USD a week from `flai cod`, E-0016's share; it replaces 85.37. Topics `planning` and `analysis` added. Details are under Notes › Planning.

Nothing to split, merge, or drop.

### 2026-10-05T05:46:33Z alex
Resolved.
