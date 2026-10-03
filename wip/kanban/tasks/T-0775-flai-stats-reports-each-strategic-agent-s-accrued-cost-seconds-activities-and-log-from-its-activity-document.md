---
id: T-0775
type: task
nature: feature
title: flai stats reports each strategic agent's accrued cost, seconds, activities, and log from its activity document
status: done
parent: S-0206
owner: alex
created: 2026-10-03T18:37:42Z
updated: 2026-10-03T18:59:43Z
transitions:
  - to: ready
    at: 2026-10-03T18:38:47Z
    by: agent-S-0206
  - to: in-progress
    at: 2026-10-03T18:46:12Z
    by: agent-S-0206
  - to: done
    at: 2026-10-03T18:59:43Z
    by: agent-S-0206
stream: S-0206
tags: []
touches: [flai/internal/metrics, flai/cmd/stats.go, design/system/metrics.md]
after: [T-0772]
usage:
  source: log
  seconds: 811
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 59
      output: 18767
      cache_read: 3999512
      cache_write: 71150
      cost: 1.6114
---
# T-0775 flai stats reports each strategic agent's accrued cost, seconds, activities, and log from its activity document

## Work

`flai stats --json` gains `strategic`: per kind with a document, its `cost`, `seconds`, `activities`, `last_run` from the front matter, and its log entries in the window (timestamp, seconds, cost, items), for S-0205's per-day charts; `flai stats` prints the totals. Define it in `design/system/metrics.md` with its precision, under the ADR T-0772 writes.

Waits for T-0772: it reads with its parser.

## Done when

- Tests pin the section on a fixture with two documents and one absent, entries outside the window left out
- metrics.md defines every field

## Notes
