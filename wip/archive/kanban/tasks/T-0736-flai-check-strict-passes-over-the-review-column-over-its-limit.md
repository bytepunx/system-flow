---
id: T-0736
type: task
nature: improvement
title: flai check --strict passes over the review column over its limit
status: done
parent: S-0243
owner: arobson
created: 2026-10-03T02:58:39Z
updated: 2026-10-03T03:04:56Z
transitions:
  - to: ready
    at: 2026-10-03T02:59:05Z
    by: agent-S-0243
  - to: in-progress
    at: 2026-10-03T02:59:05Z
    by: agent-S-0243
  - to: done
    at: 2026-10-03T03:04:56Z
    by: agent-S-0243
stream: S-0243
tags: []
touches: [flai/internal/check/check.go, flai/internal/check/check_test.go, flai/cmd/check.go, flai/cmd/check_stats_test.go]
usage:
  source: log
  seconds: 351
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 34
      output: 10126
      cache_read: 1292123
      cache_write: 36864
      cost: 0.6894
    - model: claude-sonnet-5-5
      input: 12
      output: 3133
      cache_read: 100544
      cache_write: 30263
      cost: 0.1271
---
# T-0736 flai check --strict passes over the review column over its limit

## Work

`board.wip-limit` for the review column stays a warning in `flai check`, but `--strict` no longer fails on it (TH-0078, option A): review over its limit is the operator's queue, which no story's agent can clear. The ready and in-progress limits still fail `--strict`. The result counts the warnings `--strict` passes over, and the text summary says so.

## Done when

- A test reproduces I-0007: four stories in review with a limit of 3 pass `flai check --strict`, with the warning still reported.
- Ready and in-progress over their limits still fail `--strict`.

## Notes
