---
id: T-0589
type: task
nature: feature
title: flai stats lays out burn-up, cumulative flow, and throughput over the window's days and weeks
status: done
parent: S-0166
owner: alex
created: 2026-09-29T22:47:47Z
updated: 2026-09-29T22:51:42Z
transitions:
  - to: ready
    at: 2026-09-29T22:48:15Z
    by: agent-S-0166
  - to: in-progress
    at: 2026-09-29T22:48:15Z
    by: agent-S-0166
  - to: done
    at: 2026-09-29T22:51:42Z
    by: agent-S-0166
stream: S-0166
tags: []
touches: [flai/internal/metrics]
usage:
  source: log
  seconds: 207
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 36
      output: 10552
      cache_read: 2608425
      cache_write: 41206
      cost: 1.0625
---

# T-0589 flai stats lays out burn-up, cumulative flow, and throughput over the window's days and weeks

## Work

In `flai/internal/metrics`, burn-up and cumulative flow run from the day that holds the window's start (or the first item's creation, if later) to today, instead of from the first item ever created. Throughput has one bucket per ISO week from the week that holds the window's start to the week that holds now, a week with nothing done carrying zero, instead of only the weeks something was done.

## Done when

- `flai stats --since 1d --json` and `--since 365d --json` give different `cfd` and `burnup` ranges on this repository.
- A behaviour test pins the new ranges and the zero weeks; `make flai-test` passes.

## Notes
