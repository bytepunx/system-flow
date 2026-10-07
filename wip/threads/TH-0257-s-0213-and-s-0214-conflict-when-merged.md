---
id: TH-0257
title: S-0213 and S-0214 conflict when merged
anchor:
  path: wip/kanban/stories/S-0214-charts-show-parallelism-holds-and-touches-drift.md
  item: S-0214
status: resolved
participants: [flai, agent-S-0214, alex]
created: 2026-10-07T08:23:03Z
updated: 2026-10-07T08:28:41Z
---

# TH-0257 S-0213 and S-0214 conflict when merged

On wip/kanban/stories/S-0214-charts-show-parallelism-holds-and-touches-drift.md.

## Entries

### 2026-10-07T08:23:03Z flai
A trial merge of story/S-0213 with story/S-0214 at flai stream sync conflicts in:

- `design/adrs/README.md`
- `design/system/metrics.md`

Whichever of S-0213 and S-0214 is accepted second will stop on these paths when it rebases. Settle between the two stories who changes what: one narrows its change, or names the other in `after:` and waits for it. Ask the designer when it is not clear. The next sync that finds the two merging cleanly resolves this thread.

### 2026-10-07T08:23:12Z agent-S-0214
S-0214 keeps its change and does not wait for S-0213. Both stories only add to these files, so the conflicts are positional:

- `design/adrs/README.md`: each adds one row. S-0214's ADR is ADR-0113, so the numbers no longer collide.
- `design/system/metrics.md`: S-0214 adds rows to § Claims and touches and to the Charts table. S-0213 adds rows under Cost of delay and to the same Charts table.

Whichever story is accepted second keeps both sets of rows when it rebases. S-0214 will resolve it at its next sync if S-0213 lands first. S-0214 also adds to `PLANNING_KINDS` in `flaiover/src/lib/viz/charts.ts` and to the Planning group on the Charts page, which S-0213 also extends. Those conflicts are additive too, and S-0214 resolves them the same way.

### 2026-10-07T08:28:40Z alex
Resolved.

### 2026-10-07T08:28:41Z alex
Resolved.
