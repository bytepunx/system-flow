---
id: T-0740
type: task
nature: feature
title: flai check validates the manifest's planning and warns on a draft story in ready or later
status: done
parent: S-0199
owner: alex
created: 2026-10-03T05:41:11Z
updated: 2026-10-03T06:01:03Z
transitions:
  - to: ready
    at: 2026-10-03T05:48:08Z
    by: agent-S-0199
  - to: in-progress
    at: 2026-10-03T05:48:08Z
    by: agent-S-0199
  - to: done
    at: 2026-10-03T06:01:03Z
    by: agent-S-0199
stream: S-0199
tags: []
touches: [flai/internal/check]
after: [T-0739]
usage:
  source: log
  seconds: 775
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 42
      output: 15427
      cache_read: 2494432
      cache_write: 62928
      cost: 1.1788
---
# T-0740 flai check validates the manifest's planning and warns on a draft story in ready or later

## Work

In `flai/internal/check`, report the manifest's `planning` (currency code, hour rate, cycle duration) when it is wrong, and warn (`story.draft`) on a story with `draft: true` in `ready`, `in-progress`, `review`, or `done`. Item shapes reach check through `Validate` as `item.front-matter`. Waits for T-0739, whose types and manifest keys it reads.

## Done when

`flai check` warns on a draft story in ready or later and reports a bad `planning` block, each with a test.

## Notes
