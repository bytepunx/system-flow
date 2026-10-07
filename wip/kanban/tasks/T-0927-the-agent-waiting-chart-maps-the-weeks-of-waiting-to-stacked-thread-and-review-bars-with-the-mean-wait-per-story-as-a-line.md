---
id: T-0927
type: task
nature: feature
title: The agent-waiting chart maps the weeks of waiting to stacked thread and review bars with the mean wait per story as a line
status: done
parent: S-0215
owner: alex
created: 2026-10-05T05:45:02Z
updated: 2026-10-07T09:01:48Z
transitions:
  - to: ready
    at: 2026-10-07T08:54:32Z
    by: agent-S-0215
  - to: in-progress
    at: 2026-10-07T08:55:43Z
    by: agent-S-0215
  - to: done
    at: 2026-10-07T09:01:47Z
    by: agent-S-0215
stream: S-0215
tags: [dashboard]
touches: [flaiover/src/lib/viz, flai/internal/metrics/waiting.go, flai/internal/metrics/waiting_test.go]
after: [T-0919]
usage:
  source: log
  seconds: 364
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 83
      output: 28225
      cache_read: 3773337
      cache_write: 132082
      cost: 2.1173
---
# T-0927 The agent-waiting chart maps the weeks of waiting to stacked thread and review bars with the mean wait per story as a line

## Work

In `flaiover/src/lib/viz/charts.ts`:

- Give `Report` the `waiting` section, `weeks[]` and `longest[]`, and give `ItemMetrics` `wait_threads_seconds` and `wait_review_seconds`, as `metrics.md` defines them.
- Add `agent-waiting` to `FLOW_KINDS`, so the Charts menu lists it under Flow, and give it a title in `TITLES`.
- Write its builder: one stacked bar per week of the window, hours waited on threads and on review (`total_seconds` / 3600), and a line of the mean wait per story, `(threads.total_seconds + review.total_seconds) / items`, in hours. The line is absent in a week with no items. The weeks span the window, as ADR-0054 requires.

It waits for T-0919 for the shape of `waiting`. It does not need T-0923's flai code, so it runs beside it.

## Done when

- `/charts/agent-waiting` renders the bars and the line from `/api/stats`, and the Charts menu lists it under Flow.
- Tests in `flaiover/src/lib/viz/charts.test.ts` cover the mapping: a week's two bars, the mean line, a week with no items, and the time axis from the window's start.
- The dashboard's lint, type check, and unit tests pass.

## Notes
