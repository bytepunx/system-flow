---
id: T-0580
type: task
nature: feature
title: The usage metrics define spend per bucket of time, per item, per agent minute, and per dollar, with an ADR
status: done
parent: S-0163
owner: alex
created: 2026-09-29T20:33:47Z
updated: 2026-09-29T20:36:16Z
transitions:
  - to: ready
    at: 2026-09-29T20:34:13Z
    by: agent-S-0163
  - to: in-progress
    at: 2026-09-29T20:34:13Z
    by: agent-S-0163
  - to: done
    at: 2026-09-29T20:36:16Z
    by: agent-S-0163
stream: S-0163
tags: []
touches: [design/system/metrics.md, design/adrs]
usage:
  source: log
  seconds: 123
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 10
      output: 6572
      cache_read: 1114175
      cache_write: 18893
      cost: 0.9851
---
# T-0580 The usage metrics define spend per bucket of time, per item, per agent minute, and per dollar, with an ADR

## Work

Write the ADR that changes the usage metrics of `design/system/metrics.md` (ADR-0051 point 3 is refined, not reversed): agent time is counted in minutes, not hours; spend is laid out over time in buckets of an hour, a day, or a week, by the moment an item entered done, for each item type at once; each bucket carries the items, tokens, cost, and agent seconds, their averages per item, the tokens per agent minute and per dollar, the same per model, and the running mean of tokens and cost per bucket. Rewrite the `## Usage` section and the usage rows of `## Charts` in `design/system/metrics.md` to say so, naming the JSON keys.

## Done when

- The ADR is in `design/adrs` with its row in the index, written with `flai adr new`.
- `design/system/metrics.md` defines every value the new charts plot and links the ADR.
- `flai check --strict` and the markdown lint pass.

## Notes
