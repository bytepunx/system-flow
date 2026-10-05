---
id: TH-0142
title: "S-0215 plan: five tasks in three layers, with a waiting.longest[] extension to flai stats"
anchor:
  path: wip/kanban/stories/S-0215-a-chart-shows-how-long-agents-spend-waiting-on-threads-and-on-review.md
  item: S-0215
status: resolved
participants: [planner-S-0215, alex]
created: 2026-10-05T05:45:49Z
updated: 2026-10-05T05:50:03Z
---

# TH-0142 S-0215 plan: five tasks in three layers, with a waiting.longest[] extension to flai stats

On wip/kanban/stories/S-0215-a-chart-shows-how-long-agents-spend-waiting-on-threads-and-on-review.md.

## Entries

### 2026-10-05T05:45:49Z planner-S-0215
The plan for S-0215: five tasks in three layers.

- Layer 1: T-0919. `metrics.md` defines `waiting.longest[]` (the longest waits that overlap the window, at most 10, each with `item`, `kind`, `thread`, `started`, `ended`, `seconds`, and `awaited`), and an ADR records it.
- Layer 2, run together: T-0923 builds it in `flai/internal/metrics`, and T-0927 maps the waiting weeks to the stacked chart in `flaiover/src/lib/viz/charts.ts`, adding `agent-waiting` to `FLOW_KINDS`. Both wait only for T-0919's contract and share no path.
- Layer 3, run together: T-0932 adds `WaitTable.svelte` under the chart and checks the page against `flai stats --json`, after T-0923 and T-0927. T-0934 describes the chart in the dashboard design and the user guide, after T-0927.

Assumptions:

1. The table needs data `flai stats` does not produce yet. Today's `waiting` holds totals per story and per week only, with no thread ID and no one awaited. So the story extends the metrics contract, which needs an ADR (T-0919). The criterion "matches `flai stats --json`" points this way rather than to working the table out in the dashboard.
2. "Who was awaited" is whoever ended the wait: the author of the first later entry by someone else on a thread, and the `by` of the move out of `review` for a review wait. It is empty while the wait is open, since a thread names no addressee.
3. "Longest waits in the window" means waits overlapping the window, counted only for the part inside it, with open ones included and marked open, at most 10. A thread wait counts only while the story was in progress, as `wait_threads_seconds` already counts it.
4. The mean wait per story is (thread total + review total) / the stories completed that week, the week's `items`. It is absent in a week with none.
5. "Listed under Flow" is met by `FLOW_KINDS` on the charts page. `sitemenu.ts` stays a declared touch but should not change.

Figures: forecast raised from flai's 31m to 1h15m (four layers of work, against flai's count of touches), delivery 2026-10-05T15:44Z. Cost of delay 76.01 USD a week, the story's share of E-0016's value, from `flai cod`. I propose no split, merge, or drop.

### 2026-10-05T05:50:03Z alex
Resolved.
