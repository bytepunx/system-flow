---
id: T-0598
type: task
nature: remediation
title: Remove aging and estimates, and retitle the charts
status: done
parent: S-0169
owner: alex
created: 2026-09-30T00:21:39Z
updated: 2026-09-30T00:24:34Z
transitions:
  - to: ready
    at: 2026-09-30T00:21:49Z
    by: agent-S-0169
  - to: in-progress
    at: 2026-09-30T00:21:49Z
    by: agent-S-0169
  - to: done
    at: 2026-09-30T00:24:34Z
    by: agent-S-0169
stream: S-0169
tags: []
touches: [flaiover/src]
usage:
  source: log
  seconds: 165
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 25
      output: 7421
      cache_read: 1886698
      cache_write: 26004
      cost: 0.7339
---
# T-0598 Remove aging and estimates, and retitle the charts

## Work

In flaiover/src/lib/viz/charts.ts drop the aging and estimates kinds and their builders, rename the titles the story lists, and drop their branches from the charts page. Update the unit and page tests.

## Done when

- The flow nav lists Cycle Time, Burn-up, Cumulative Flow, Time in State, Throughput only.
- The usage titles read Tokens / Min, Tokens / Day (by bucket), Tokens / $, $ / Day (by bucket), $ / Work Type, $ / Item.
- flaiover tests and lint pass.

## Notes
