---
id: T-0581
type: task
nature: feature
title: flai stats lays spend out per bucket and per item type, and counts agent time in minutes
status: done
parent: S-0163
owner: alex
created: 2026-09-29T20:33:47Z
updated: 2026-09-29T20:40:56Z
transitions:
  - to: ready
    at: 2026-09-29T20:34:13Z
    by: agent-S-0163
  - to: in-progress
    at: 2026-09-29T20:36:16Z
    by: agent-S-0163
  - to: done
    at: 2026-09-29T20:40:56Z
    by: agent-S-0163
stream: S-0163
tags: []
touches: [flai/internal/metrics, flai/internal/hostapi, flai/cmd, docs/users]
usage:
  source: log
  seconds: 280
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 24
      output: 250
      cache_read: 2242934
      cache_write: 37468
      cost: 3.7175
---
# T-0581 flai stats lays spend out per bucket and per item type, and counts agent time in minutes

## Work

In `flai/internal/metrics`, compute the spend series the ADR defines: for each of epic, story, and task, whatever `--type` is, one point per bucket from the first with spend in the window to now, with totals, averages per item, tokens per agent minute, tokens per dollar, the same per model, and the running means. Add tokens per agent minute beside tokens per hour on items and models. `flai stats` takes `--bucket hour|day|week` (day by default; an hour only with a window of 31 days or less) and its table prints the averages per item, per agent minute, and per dollar. `stats.get` takes `bucket` and answers as the command does. Update `docs/users/flai.md` and regenerate the reference.

## Done when

- Behaviour tests in `flai/internal/metrics` fail without the change and pass with it: buckets, averages, per model, running mean, an empty bucket between two with spend, a window with none.
- The contract test in `flai/cmd` holds `stats.get` with a bucket to `flai stats --bucket`.
- `make test` and golangci-lint pass.

## Notes
