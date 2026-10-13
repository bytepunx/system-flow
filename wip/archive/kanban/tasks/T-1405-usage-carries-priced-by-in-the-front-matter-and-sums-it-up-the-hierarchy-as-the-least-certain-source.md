---
id: T-1405
type: task
nature: feature
title: usage carries priced_by in the front matter and sums it up the hierarchy as the least certain source
status: done
parent: S-0357
owner: alex
created: 2026-10-08T08:50:31Z
updated: 2026-10-13T19:25:30Z
transitions:
  - to: ready
    at: 2026-10-13T19:17:07Z
    by: agent-S-0357
  - to: in-progress
    at: 2026-10-13T19:17:07Z
    by: agent-S-0357
  - to: done
    at: 2026-10-13T19:25:30Z
    by: agent-S-0357
stream: S-0357
tags: [cli]
touches: [flai/internal/usage/usage.go, flai/internal/usage/usage_test.go, flai/internal/workitem/usage.go, flai/internal/workitem/usage_test.go, flai/internal/workitem/front-matter-fields.txt, design/system/work-hierarchy.md]
usage:
  source: log
  seconds: 503
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 24
      output: 12568
      cache_read: 1543081
      cache_write: 46927
      cost: 0.8422
---
# T-1405 usage carries priced_by in the front matter and sums it up the hierarchy as the least certain source

## Work

- `usage.Usage` gains `PricedBy string` (`priced_by`, omitempty), one of the constants `PricedGateway` (`gateway`), `PricedHarness` (`harness`), `PricedEstimate` (`estimate`); empty on usage written before it. Strategic entries do not carry it: they are apportioned, always `estimated`.
- A function such as `LeastCertain(a, b string) string`, and an accessor `(u *Usage) Priced() string` that reads an empty value as `harness` when `estimated` is false and `estimate` when it is true, so old items sum as they read. Only usage with a cost or models counts: an empty usage adds nothing to the sum.
- `Add` and `Sum` set `PricedBy` to the least certain of the parts: `estimate` below `harness` below `gateway`. `Split` keeps it on each share. `Same` compares it.
- `workitem/usage.go` writes `priced_by` in the front matter beside `estimated`, reads it, and refuses a value outside the three (an error in `usageErrors`, as `source` is checked); `Summed` and `RollUp` carry it through `usage.Sum`.
- `design/system/work-hierarchy.md`: the usage example and the bullet list gain `priced_by` and its sum, citing ADR-0132.
- Tests: the order of the sum, an empty value from an older item read by `estimated`, a refusal of an unknown value, and a round trip through front matter.

First layer: it runs together with T-1406 and T-1408, which touch other files.

## Done when

- `flai test flai/internal/usage/ flai/internal/workitem/ design/system/work-hierarchy.md` passes.

## Notes

Layer 1 of S-0357. Meets part of criterion 4.
