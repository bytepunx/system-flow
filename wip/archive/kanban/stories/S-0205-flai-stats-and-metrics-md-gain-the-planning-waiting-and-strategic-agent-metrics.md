---
id: S-0205
type: story
nature: feature
title: flai stats and metrics.md gain the planning, waiting, and strategic-agent metrics
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:13Z
updated: 2026-10-04T00:39:24Z
transitions:
  - to: ready
    at: 2026-10-03T05:34:33Z
    by: alex
  - to: in-progress
    at: 2026-10-03T20:33:46Z
    by: system-flow
  - to: review
    at: 2026-10-03T21:43:22Z
    by: agent-S-0205
  - to: done
    at: 2026-10-04T00:39:24Z
    by: alex
tags: [flai]
touches: [flai/internal/metrics, flai/cmd/stats.go, flai/cmd/check_stats_test.go, flai/cmd/hostapi_reads_test.go, flai/internal/hostapi/reads.go, flai/internal/statsread, design/system/metrics.md, flai/internal/usage, flai/internal/storygit/committed.go, flai/internal/storygit/committed_test.go, design/adrs, docs/users/flai.md, docs/users/flai-reference.md, design/issues/summary.md, design/issues/I-0068-the-mcp-prime-pack-for-a-story-is-larger-than-claude-code-s-tool-result-limit-so-the-agent-reads-it-back-from-a-saved-file.md, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md]
after: [S-0199, S-0206]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4208
  models:
    - model: claude-haiku-4-5-20251001
      input: 178
      output: 8046
      cache_read: 1036183
      cache_write: 73170
      cost: 0.2355
    - model: claude-opus-5-5
      input: 462
      output: 193372
      cache_read: 27921487
      cache_write: 755842
      cost: 13.9815
    - model: claude-sonnet-5
      input: 126
      output: 29063
      cache_read: 5220941
      cache_write: 349253
      cost: 2.2082
---
# S-0205 flai stats and metrics.md gain the planning, waiting, and strategic-agent metrics

## Goal

The charts E-0016 asks for need numbers `flai stats` does not compute: forecast against actual, cost of delay outstanding and incurred, time agents spend waiting, claim and touches drift, and what the planner, orchestrator, and analyzer cost against delivery. `design/system/metrics.md` is the contract between `flai stats` and the dashboard, so this is an ADR.

## Acceptance criteria
- [x] Per story: `forecast_error` (actual cycle time minus `forecast.duration`), `delivery_error` (completed minus `forecast.delivery`), `estimate_error` (against the human `estimate`), each absent when the input is; aggregates p50 and p85 of absolute error, by nature and by the story's agent model
- [x] Cost of delay: per item its `value`; per day of the window the value outstanding per column (summing open items' values) and the value incurred (value × days waited in backlog and ready); per week the total incurred
- [x] Waiting: per story the time its agent waited, from each thread's opened-to-answered interval while the story was in progress, plus the review interval (review to done); aggregates mean and total per week, split into `threads` and `review`
- [x] Claims and touches: per story the time it was held (from the launcher's hold reasons where recorded), stories in progress per day against the WIP limit, and touches drift: files committed outside its declared touches and declared touches never changed, as counts and paths
- [x] Strategic agents: per day of the window the cost and agent seconds of the planner, orchestrator, and analyzer (from their activity logs' front matter), beside the mean cost and cycle time per story completed that day, so the dashboard can chart Strategic Cost and Strategic Use
- [x] `flai stats --json` carries all of it; `flai stats` prints the aggregates; precision rules are written for each in `metrics.md`, and an ADR records the additions
- [x] Tests pin each metric on fixtures, to the second

## Tasks
- T-0777 metrics.md defines the planning, waiting, claims, and strategic-use metrics, and an ADR records them
- T-0778 flai stats reports forecast, delivery, and estimate error per story and their p50 and p85
- T-0779 flai stats reports cost of delay outstanding per column and incurred per day and week
- T-0780 flai stats reports the time a story's agent waited on threads and on review
- T-0781 flai stats reports held time, stories in progress against the limit, and touches drift
- T-0782 flai stats reports the strategic agents' cost and seconds per day beside delivery
- T-0783 flai stats reads threads, the limit, and committed files, and prints the new aggregates

## Notes
- The errors are keyed with their unit, `forecast_error_seconds`, `delivery_error_seconds`, and `estimate_error_seconds`, as every duration in `flai stats --json` is. The relative `estimate_error` stays beside them (ADR-0081).
- flai records no hold, so held time is replayed from the hold rules over the states of the time, with today's touches and `after` (ADR-0081).
- An activity document's front matter holds only all-time totals, so the strategic agents' use per day sums its log entries by the day they ended.
- The dashboard's `stats.get` and `flai stats` read their inputs through one reader, `flai/internal/statsread`, so both carry the new sections.
