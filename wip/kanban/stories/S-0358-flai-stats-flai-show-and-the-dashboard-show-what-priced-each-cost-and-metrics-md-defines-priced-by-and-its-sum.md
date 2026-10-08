---
id: S-0358
type: story
nature: feature
title: flai stats, flai show, and the dashboard show what priced each cost, and metrics.md defines priced_by and its sum
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:51:30Z
updated: 2026-10-08T09:30:45Z
transitions: []
tags: [cli, dashboard]
topics: [usage, metrics]
touches: [design/system/metrics.md, flai/internal/metrics/usage.go, flai/internal/metrics/usage_test.go, flai/internal/metrics/spend.go, flai/internal/metrics/spend_test.go, flai/internal/usage/usage.go, flai/internal/usage/usage_test.go, flai/cmd/stats.go, flai/cmd/items.go, flaiover/src/lib/usage.ts, flaiover/src/lib/usage.test.ts, flaiover/src/lib/components/SpendTable.svelte, flaiover/src/lib/components/SpendTable.svelte.test.ts, docs/users/flai.md]
after: [S-0357]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 33.19
  by: planner-E-0019
  at: 2026-10-08T09:00:26Z
forecast:
  duration: 30m
  delivery: 2026-10-09T01:21:00Z
  basis: "Its own forecast of 30m; 30th in the pull order with an in-progress limit of 5, behind S-0232, S-0341, S-0344, S-0346, S-0348, S-0338, S-0342, S-0337, S-0334, S-0343, S-0297, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0289, S-0304, S-0305, S-0306, S-0313, S-0347, S-0349, S-0350, S-0351, S-0352, S-0353, S-0354, S-0355, S-0356 and S-0357."
  by: flai
  at: 2026-10-08T09:30:45Z
finalized:
  by: alex
  at: 2026-10-08T08:52:51Z
---
# S-0358 flai stats, flai show, and the dashboard show what priced each cost, and metrics.md defines priced_by and its sum

## Goal

S-0357 records `priced_by` on each item's usage. [ADR-0132](../../../design/adrs/0132-a-story-s-cost-comes-from-its-provider-s-spend-log-when-it-has-a-provider-from.md) has `flai stats` and the dashboard show it beside the cost, and `design/system/metrics.md`, the contract between the two, define it and its sum. Today a cost is either reported or "(estimated)"; this story names three sources, so that a reader can tell what the gateway charged from what Claude Code estimated.

## Acceptance criteria

- [ ] `design/system/metrics.md` defines `priced_by` (`gateway`, `harness`, `estimate`), its sum up the hierarchy as the least certain source, how an item written before it reads, and the cost per source `flai stats` adds to its usage totals and to spend over time.
- [ ] `flai stats --json` carries `priced_by` on each item's usage and the cost per source in the window's totals and in each spend bucket; `flai stats` and `flai show` print the source beside a cost, as `(gateway)`, `(harness)`, or `(estimated)`.
- [ ] The dashboard's usage lines and spend table show the source beside each cost, and an item with no `priced_by` shows as today.
- [ ] Tests in `flai/internal/metrics`, `flai/internal/usage`, and the two flaiover files cover the three sources and an older item; `docs/users/flai.md` says what each source means.

## Tasks

- T-1411 metrics.md defines priced_by, its sum, how an older item reads, and the cost per source
- T-1412 flai stats reports priced_by and the cost per source, and flai stats and flai show print the source beside a cost
- T-1413 The dashboard's usage lines and spend table show the source beside each cost
- T-1414 The user guide says what gateway, harness, and estimated mean beside a cost

## Notes

- `metrics.md` changes under ADR-0132, which names this consequence; no further ADR is needed.

### Planning

Planned by planner-E-0019 on 2026-10-08. Every touch is a file; no folder touch is kept. Three layers: the contract (T-1411), then flai's side (T-1412), then the dashboard (T-1413) and the user guide (T-1414) together.

| Touch | Source | Why |
|-------|--------|-----|
| `design/system/metrics.md` | design | The contract between `flai stats` and the dashboard (T-1411) |
| `flai/internal/metrics/usage.go`, `usage_test.go`, `spend.go`, `spend_test.go` | layout | The usage totals and spend over time (T-1412) |
| `flai/internal/usage/usage.go`, `usage_test.go`, `flai/cmd/stats.go`, `flai/cmd/items.go` | layout | `Summary()` prints `(estimated)` today, used by `flai stats` and `flai show` (T-1412) |
| `flaiover/src/lib/usage.ts`, `usage.test.ts`, `flaiover/src/lib/components/SpendTable.svelte`, `SpendTable.svelte.test.ts` | layout | Where the dashboard prints a cost's mark (T-1413) |
| `docs/users/flai.md` | design | What each source means (T-1414) |

`touches suggest` listed `design/system/flai-cli.md`, `docs/users/flai-reference.md`, and the dashboard's documents. None is taken: no command or flag changes, and the dashboard's design names no cost mark.

Forecast: 30m, as `flai forecast` gives it, 100 s per unit over 34 done feature stories on claude-opus-5-5 in the large band, times size 18. It stands.

Cost of delay: 33.19 USD a week, as `flai cod` works it out: 30m of the 7h32m forecast over E-0019's 12 stories, of the operator's 500 USD a week penalty. It stands.
