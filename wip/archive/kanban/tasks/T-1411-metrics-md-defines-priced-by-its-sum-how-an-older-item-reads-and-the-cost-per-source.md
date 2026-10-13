---
id: T-1411
type: task
nature: feature
title: metrics.md defines priced_by, its sum, how an older item reads, and the cost per source
status: done
parent: S-0358
owner: alex
created: 2026-10-08T08:51:37Z
updated: 2026-10-13T21:31:37Z
transitions:
  - to: ready
    at: 2026-10-13T21:30:28Z
    by: system-flow
  - to: in-progress
    at: 2026-10-13T21:30:29Z
    by: system-flow
  - to: done
    at: 2026-10-13T21:31:37Z
    by: system-flow
stream: S-0358
tags: [cli]
touches: [design/system/metrics.md]
usage:
  source: log
  seconds: 68
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 17
      output: 5652
      cache_read: 916938
      cache_write: 32670
      cost: 0.5123
---
# T-1411 metrics.md defines priced_by, its sum, how an older item reads, and the cost per source

## Work

- In `design/system/metrics.md` § Usage: `priced_by` and its three values; the sum up the hierarchy as the least certain (`estimate` below `harness` below `gateway`); an item with none reads as `estimate` when `estimated` is true and `harness` otherwise; and `cost_by_source`, the cost per source, in the window's usage totals and in each spend-over-time bucket. Cite ADR-0132.
- Name the JSON fields exactly, since the dashboard reads them.

First layer: the contract is written before the code on both sides of it.

## Done when

- The section names every field T-1412 and T-1413 will use.
- `flai test design/system/metrics.md` passes.

## Notes

Layer 1 of S-0358.
