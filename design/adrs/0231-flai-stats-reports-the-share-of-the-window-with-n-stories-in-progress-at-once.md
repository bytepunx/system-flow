---
id: ADR-0231
title: "flai stats reports the share of the window with n stories in progress at once, and the ready wait, the queue time, at p50 and p85"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0113, ADR-0190]
---

# ADR-0231 flai stats reports the share of the window with n stories in progress at once, and the ready wait, the queue time, at p50 and p85

## Context

S-0551 found the board idle while stories waited: with an in-progress limit of 5, nothing was in progress for 42 percent of the fortnight to 2026-10-30, one story for 29 percent, two for 26 percent, three for 2 percent, and four or five never. The mean while anything ran was 1.6. The journal recorded 124 overlap holds and 113 `after` holds against 62 real conflicts, and the ready wait at p85 was 5.4 hours. [ADR-0230](0230-a-ready-story-whose-claim-overlaps-a-story-in-progress-is-held-only-when-a.md) changes the holds to raise those figures. They were worked out by hand from the transitions, and `flai stats` reports neither.

`claims.days[].in_progress` ([ADR-0113](0113-flai-stats-reports-held-stories-per-day-held-time-by-reason-per-week-and-the.md)) counts the items in progress at the end of each day. A count at one moment a day misses a story that started and moved to review within the day, and says nothing of the hours between. Queue time ([ADR-0190](0190-an-item-s-queue-time-is-the-seconds-it-spent-in-ready-before-it-started-so-a.md)) is the time an item waited in `ready` before it started. `summary.queue_time` gives its p50 and p85, but under a name the story's goal and the board do not use. `design/system/metrics.md` is the contract between `flai stats` and the dashboard, and changes only with an ADR.

## Decision

**`flai stats` reports, for the stories of the report's type, the share of the window with n of them in progress at once, for each n from 0 to the in-progress limit or the most seen, and the ready wait, the queue time of ADR-0190, at p50 and p85.** This refines ADR-0113 and ADR-0190, as `metrics.md` § Concurrency and the ready wait defines it.

- **`concurrency`.** Over the window, from `window_start` to `generated_at`, the seconds with exactly n items of the report's type in `in-progress` at once, from their state intervals, archived and cancelled items included. `seconds[n]` and `share[n]`, n from 0, through the larger of the board's `in-progress` limit and the most in progress at once in the window. `share[n]` is `seconds[n]` over the window's seconds, so the shares sum to 1. `mean_busy` is the mean n over the seconds with one or more in progress.
- **`ready_wait`.** `count`, `p50_seconds`, and `p85_seconds` of the queue time of the items `summary.queue_time` counts. It is that figure under the name the board uses, so the two agree.
- **The text report.** `flai stats` prints the queue time line as `ready wait`, and under it the shares as percents with the mean.
- **The dashboard.** The flow page charts the shares as bars per n, with the ready wait's p50 and p85 beside them.

## Consequences

- The effect of ADR-0230 can be read from `flai stats` and the dashboard: the share at 0 should fall and the shares above 1 rise, and the ready wait should shorten.
- The shares are exact to the second, unlike the day-end count of `claims.days[].in_progress`, which the Parallelism chart keeps.
- `ready_wait` repeats `summary.queue_time`'s p50 and p85. A reader of either finds the same numbers.
- A story in review is not in progress here, as it is not for the in-progress limit.

## Alternatives considered

- **Read concurrency from `claims.days[].in_progress`.** One sample a day cannot give a share of the hours.
- **Count review as in progress.** The limit counts `in-progress` alone, and a story in review uses no agent.
- **A new ready wait over the stories started in the window.** Its figure would differ from the queue time the dashboard already shows for the same stories.
- **Rename `summary.queue_time`.** It would break the dashboard and every script that reads it.
