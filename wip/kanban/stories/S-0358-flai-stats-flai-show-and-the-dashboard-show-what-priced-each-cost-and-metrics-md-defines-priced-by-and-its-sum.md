---
id: S-0358
type: story
nature: feature
title: flai stats, flai show, and the dashboard show what priced each cost, and metrics.md defines priced_by and its sum
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:51:30Z
updated: 2026-10-08T08:52:51Z
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

Drafted by the planner; see the children.
- T-1411 metrics.md defines priced_by, its sum, how an older item reads, and the cost per source
- T-1412 flai stats reports priced_by and the cost per source, and flai stats and flai show print the source beside a cost
- T-1413 The dashboard's usage lines and spend table show the source beside each cost
- T-1414 The user guide says what gateway, harness, and estimated mean beside a cost

## Notes

- `metrics.md` changes under ADR-0132, which names this consequence; no further ADR is needed.
