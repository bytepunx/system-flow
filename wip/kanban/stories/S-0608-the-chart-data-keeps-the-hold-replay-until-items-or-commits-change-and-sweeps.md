---
id: S-0608
type: story
nature: improvement
title: The chart data keeps the hold replay until items or commits change and sweeps the days once, so a 365-day window costs what a 30-day one does
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:41:41Z
updated: 2026-11-07T22:41:41Z
transitions: []
tags: [cli]
touches: [flai/internal/metrics/claims.go, flai/internal/metrics/costofdelay.go, flai/internal/metrics/metrics.go, flai/internal/metrics/costorder.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0608 The chart data keeps the hold replay until items or commits change and sweeps the days once, so a 365-day window costs what a 30-day one does

## Goal

Cause 18 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `replayHolds` replays all history on every call, whatever the window (`flai/internal/metrics/claims.go`), 175 to 310 ms and about 60 MB. `costOfDelay` and `claims` loop day by day over every item, parse every transition's time again each day through `stateAt`, and start at the window's start even before the first item existed: `costOfDelay` takes 2 ms at 1 day and 170 to 217 ms at 365 days. Thread entries are parsed twice, and `costOrder` builds its forecaster three times. A warm 365-day `stats.get` takes 1,067 to 1,118 ms against 672 to 816 ms at 30 days.

## Acceptance criteria

- [ ] `replayHolds` and `touchDrift` are kept between requests until the work items or the stories' commits change, and only the slicing to the window is done per request.
- [ ] Each item's transition times are parsed once a call, the day-by-day sections sweep the days once, and the days before the first item are answered without a loop; thread entries are parsed once and the forecaster built once.
- [ ] The report is the same as before for windows of 1, 7, 30, 90, and 365 days, which a test shows, and a warm 365-day `stats.get` on this repository takes under 1.2 times a warm 30-day one; server-performance.md records the before and after.

## Tasks

## Notes

- The answer and `metrics.md` do not change.
- Keeping needs a generation count on the item store (`flai/internal/workitem/store.go`), in memory.
