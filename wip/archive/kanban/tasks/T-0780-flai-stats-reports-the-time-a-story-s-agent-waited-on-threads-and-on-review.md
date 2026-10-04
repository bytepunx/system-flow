---
id: T-0780
type: task
nature: feature
title: flai stats reports the time a story's agent waited on threads and on review
status: done
parent: S-0205
owner: alex
created: 2026-10-03T20:40:04Z
updated: 2026-10-03T20:58:33Z
transitions:
  - to: ready
    at: 2026-10-03T20:52:52Z
    by: agent-S-0205
  - to: in-progress
    at: 2026-10-03T20:52:52Z
    by: agent-S-0205
  - to: done
    at: 2026-10-03T20:58:33Z
    by: agent-S-0205
stream: S-0205
tags: []
touches: [flai/internal/metrics]
after: [T-0779]
usage:
  source: log
  seconds: 341
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 28
      output: 11544
      cache_read: 1666844
      cache_write: 45122
      cost: 0.8347
---
# T-0780 flai stats reports the time a story's agent waited on threads and on review

## Work

In `flai/internal/metrics`, per item `wait_threads_seconds` and `wait_review_seconds`, and `waiting.weeks`, from the threads `Options` carries, as `metrics.md` § Waiting defines them. Waits for the cost of delay task: both change `metrics.go`.

## Done when

- [ ] `Compute` reports each value from threads given in `Options`
- [ ] Tests pin them on a fixture of their own, to the second

## Notes
