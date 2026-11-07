---
id: T-2609
type: task
nature: research
title: "Explore and measure the chart data: stats.get over windows of 7, 30, 90 days and all history"
status: done
parent: S-0596
owner: alex
created: 2026-11-07T22:10:35Z
updated: 2026-11-07T22:36:16Z
transitions:
  - to: ready
    at: 2026-11-07T22:12:55Z
    by: agent-S-0596
  - to: in-progress
    at: 2026-11-07T22:12:58Z
    by: agent-S-0596
  - to: done
    at: 2026-11-07T22:36:16Z
    by: agent-S-0596
stream: S-0596
tags: [cli]
usage:
  source: log
  seconds: 1398
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 88
      output: 25447
      cache_read: 5756836
      cache_write: 109206
      cost: 2.344
---
# T-2609 Explore and measure the chart data: stats.get over windows of 7, 30, 90 days and all history

## Work

The serve log since 2026-11-04 shows `stats.get` at p50 1.2 s and up to 6.5 s, nearly all in `stats.compute`. Time `stats.get` with `flai hostapi --timing` for each window the dashboard offers, cold and warm, and profile it. Read `flai/internal/metrics` and what it reads (items, archive, threads, activity documents, event logs under `.flai-cache`, git). Say how the cost grows with the window and with the archive, which reads or computations run one after another where they could run at once, and what is recomputed on every request that could be kept. Waits for no task: it writes nothing and returns its findings to the story's agent.

## Done when

- [ ] The findings give stats.get's cost per window, cold and warm, where the time goes by path and line, and each serial step that could run in parallel or be kept, with a measurement.

## Notes
