---
id: T-0582
type: task
nature: feature
title: The dashboard builds the token and cost charts from the spend series
status: done
parent: S-0163
owner: alex
created: 2026-09-29T20:33:47Z
updated: 2026-09-29T20:50:30Z
transitions:
  - to: ready
    at: 2026-09-29T20:34:13Z
    by: agent-S-0163
  - to: in-progress
    at: 2026-09-29T20:40:56Z
    by: agent-S-0163
  - to: done
    at: 2026-09-29T20:50:30Z
    by: agent-S-0163
stream: S-0163
tags: []
touches: [flaiover/src/lib]
usage:
  source: log
  seconds: 574
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 54
      output: 389
      cache_read: 6583078
      cache_write: 87639
      cost: 10.8739
---
# T-0582 The dashboard builds the token and cost charts from the spend series

## Work

In `flaiover/src/lib/viz/charts.ts`, replace the token rate scatter with tokens per agent minute over time, per model, and add builders for tokens per bucket with its running mean, tokens per item by type or by model, tokens per dollar, cost per bucket with its running mean, and cost per item by type or by model. Keep the rules of `design/tech/charts.md`: one y-axis, thin marks, a legend for two or more series, a tooltip on every mark, models in their fixed slots. Item types get fixed slots too.

## Done when

- Unit tests in `charts.test.ts` cover each new builder's series and values against a report with two models, three types, and an empty bucket, and a report from a flai that sends no spend series.
- The fixture test through flai reads the spend series.
- Prettier, eslint, svelte-check, and vitest pass.

## Notes
