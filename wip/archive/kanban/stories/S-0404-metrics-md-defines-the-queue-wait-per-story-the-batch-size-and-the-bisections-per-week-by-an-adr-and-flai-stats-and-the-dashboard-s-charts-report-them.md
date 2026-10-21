---
id: S-0404
type: story
nature: feature
title: metrics.md defines the queue wait per story, the batch size, and the bisections per week by an ADR, and flai stats and the dashboard's charts report them
status: done
parent: E-0022
owner: alex
created: 2026-10-16T04:58:58Z
updated: 2026-10-20T19:00:27Z
transitions:
  - to: ready
    at: 2026-10-19T04:53:59Z
    by: alex
  - to: in-progress
    at: 2026-10-19T23:41:29Z
    by: agent-S-0404
  - to: review
    at: 2026-10-20T00:55:34Z
    by: agent-S-0404
  - to: done
    at: 2026-10-20T19:00:27Z
    by: alex
tags: [cli, dashboard]
topics: [analysis, orchestration]
touches: [design/adrs, design/system/metrics.md, flai/internal/mergequeue/log.go, flai/internal/mergequeue/log_test.go, flai/internal/mergequeue/state.go, flai/internal/metrics/queue.go, flai/internal/metrics/queue_test.go, flai/internal/metrics/metrics.go, flai/cmd/stats.go, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", design/system/flaiover-dashboard.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, flai/cmd/accept_queue_test.go, flai/cmd/stats_test.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, design/issues/summary.md, design/issues/I-0167-task-done-and-story-start-answers-exceed-claude-code-s-tool-output-limit-because-the-inbox-they-carry-lists-every-thread-awaiting-the-operator.md, docs/users/flai-reference.md]
after: [S-0400]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4464
  priced_by: harness
  turns:
    - day: 2026-10-19
      ceremony: 2
      hand_edits: 1
      work: 35
    - day: 2026-10-20
      work: 30
  models:
    - model: claude-opus-5-5
      input: 516
      output: 186648
      cache_read: 31407708
      cache_write: 1060814
      cost: 15.9835
  strategic:
    - kind: orchestrator
      seconds: 3245
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 117
          output: 2427
          cache_read: 26110493
          cache_write: 45771
          cost: 6.5575
        - model: claude-sonnet-5-5
          input: 20
          output: 119
          cache_read: 358690
          cache_write: 57866
          cost: 0.3409
forecast:
  duration: 37m
  delivery: 2026-10-20T00:22:00Z
  basis: "Its own forecast of 37m; 1st in the pull order with an in-progress limit of 5, with nothing ahead of it."
  by: flai
  at: 2026-10-19T23:41:20Z
finalized:
  by: alex
  at: 2026-10-19T04:53:51Z
---
# S-0404 metrics.md defines the queue wait per story, the batch size, and the bisections per week by an ADR, and flai stats and the dashboard's charts report them

## Goal

The analyzer and the operator can see whether the queue removed the review waits I-0150 counts. `flai stats` reports, per ISO week, how long stories waited in the queue, how large the batches were, and how often a batch was bisected, and the dashboard charts them beside the review wait.

## Acceptance criteria

- [x] An ADR adds to `design/system/metrics.md` the queue wait per story (from its enqueueing to its landing or its send-back), the batch size, and the bisections, each defined, with what each is read from.
- [x] The queue appends each event it acts on, enqueued, batch built, verified, bisected, landed, sent back, dequeued, to a log under `.flai-cache`, one line per event with the time, the batch, and the stories, as the conflict log of ADR-0138 is kept.
- [x] `flai stats` and `flai stats --json` report, under `queue` and per ISO week of the window, the stories landed by the queue with the mean and longest queue wait, the batches with their mean and largest size, and the bisections, with each item's queue wait in `items[]`. With no log, the section is absent.
- [x] The dashboard's charts draw the queue's weekly figures beside the agent-waiting chart, and say so when the flai on the host sends none. The port fixture test covers the new fields.
- [x] `design/system/flaiover-dashboard.md` and `docs/users/flai.md` describe the new figures and chart.

## Tasks
- T-1651 An ADR adds the queue wait per story, the batch size, and the bisections to metrics.md
- T-1652 The queue appends each event it acts on to a log under .flai-cache
- T-1664 flai stats reports the queue per week from its log, and each item's queue wait
- T-1673 The dashboard charts the queue's weekly wait, batch size, and bisections beside the agent-waiting chart
- T-1680 The dashboard design, the CLI design, and the user guides describe the queue's figures and chart

## Notes

- Drafted by the planner for E-0022, after S-0400, since landing and bisection are what it counts.
- Metrics definitions are the contract between `flai stats` and the dashboard, changed only with an ADR (CLAUDE.md).

### Planning

- **Touches.** None were declared; every touch is the planner's.
  - Layout, from the explorer's map:
    - waits and coordination: `flai/internal/metrics/waiting.go` and `coordination.go`, read and not changed, and the report in `metrics.go`
    - the dashboard: `flaiover/src/lib/viz/charts.ts`, its port fixture `charts.test.ts`, and the charts page.
  - The log: new files in `flai/internal/mergequeue`, written from `serve/mergequeue.go` and the state's writer.
  - Design: the epic's criterion 8; ADR-0138's conflict log as the pattern.
  - Co-change: `flai touches suggest` added `design/system/flai-cli.md` and `docs/users/flaiover.md`.
  - Folder touch kept: `design/adrs`, for the new ADR.
- **Forecast.** flai's 37m is kept: it follows ADR-0138's pattern of a log, `flai stats`, and a chart.
- **Cost of delay.** Waits on TH-0465 question 1.
