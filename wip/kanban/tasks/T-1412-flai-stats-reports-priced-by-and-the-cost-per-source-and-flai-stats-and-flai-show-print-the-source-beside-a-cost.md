---
id: T-1412
type: task
nature: feature
title: flai stats reports priced_by and the cost per source, and flai stats and flai show print the source beside a cost
status: backlog
parent: S-0358
owner: alex
created: 2026-10-08T08:51:43Z
updated: 2026-10-08T08:51:43Z
transitions: []
stream: S-0358
tags: [cli]
touches: [flai/internal/metrics/usage.go, flai/internal/metrics/usage_test.go, flai/internal/metrics/spend.go, flai/internal/metrics/spend_test.go, flai/internal/usage/usage.go, flai/internal/usage/usage_test.go, flai/cmd/stats.go, flai/cmd/items.go]
after: [T-1411]
---
# T-1412 flai stats reports priced_by and the cost per source, and flai stats and flai show print the source beside a cost

## Work

- `metrics/usage.go` and `metrics/spend.go`: `priced_by` on each item's usage in `flai stats --json`, and `cost_by_source` in the totals and in each spend bucket, as T-1411 defines them.
- `usage.Usage.Summary()` prints `(gateway)`, `(harness)`, or `(estimated)` after the cost; `flai stats` text and `flai show` (`cmd/items.go`) use it.
- Tests: an item of each source, an older item without `priced_by`, and a bucket mixing sources.

It waits for T-1411, whose fields it writes.

## Done when

- `flai test flai/internal/metrics/ flai/internal/usage/ flai/cmd/` passes.

## Notes

Layer 2 of S-0358.
