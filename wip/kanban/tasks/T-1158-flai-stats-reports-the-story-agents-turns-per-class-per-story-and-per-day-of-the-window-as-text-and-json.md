---
id: T-1158
type: task
nature: feature
title: flai stats reports the story agents' turns per class, per story and per day of the window, as text and --json
status: done
parent: S-0293
owner: alex
created: 2026-10-07T09:28:21Z
updated: 2026-10-07T09:42:57Z
transitions:
  - to: ready
    at: 2026-10-07T09:28:45Z
    by: agent-S-0293
  - to: in-progress
    at: 2026-10-07T09:37:05Z
    by: agent-S-0293
  - to: done
    at: 2026-10-07T09:42:57Z
    by: agent-S-0293
stream: S-0293
tags: []
touches: [flai/internal/metrics/turns.go, flai/internal/metrics/turns_test.go, flai/internal/metrics/metrics.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go]
after: [T-1157]
usage:
  source: log
  seconds: 352
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 24
      output: 11296
      cache_read: 1602715
      cache_write: 47548
      cost: 0.856
---
# T-1158 flai stats reports the story agents' turns per class, per story and per day of the window, as text and --json

## Work

`metrics.Compute` reads each story's `usage.turns`, archived stories included, whatever `--type` is, and reports under `turns` in `--json`: the classes in order, the totals over the window, every day of the window with each class's count over every story, and each story with turns in the window, with its ID, title, status, and its counts over the days of the window. `flai stats` prints a section with the totals, the days that have turns, and the stories. Waits for T-1157, whose `usage.turns` it reads.

## Done when

- `flai stats --json` carries `turns` with `classes`, `total`, `days`, and `stories`, and a metrics test shows days outside the window left out and two stories summed on a day.
- `flai stats` prints the section only when some story has turns in the window, and a command test in `check_stats_test.go` pins it.

## Notes
