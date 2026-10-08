---
id: T-1405
type: task
nature: feature
title: usage carries priced_by in the front matter and sums it up the hierarchy as the least certain source
status: backlog
parent: S-0357
owner: alex
created: 2026-10-08T08:50:31Z
updated: 2026-10-08T08:50:31Z
transitions: []
stream: S-0357
tags: [cli]
touches: [flai/internal/usage/usage.go, flai/internal/usage/usage_test.go, flai/internal/workitem/usage.go, flai/internal/workitem/usage_test.go]
---
# T-1405 usage carries priced_by in the front matter and sums it up the hierarchy as the least certain source

## Work

- `usage.Usage` gains `PricedBy string` (`priced_by`, one of `gateway`, `harness`, `estimate`, empty on usage written before it), and the strategic entries the same.
- `Add` and `Sum` take the least certain of the two: `estimate` below `harness` below `gateway`; an empty one counts as `harness` when `estimated` is false and `estimate` when it is true, so old items sum as they read.
- `workitem/usage.go` writes and reads `priced_by` in the front matter in its place beside `estimated`, and `Summed` carries it.
- Tests: the order of the sum, an empty value from an older item, and a round trip through front matter.

First layer: it runs together with T-1406 and T-1408, which touch other files.

## Done when

- `flai test flai/internal/usage/ flai/internal/workitem/` passes.

## Notes

Layer 1 of S-0357.
