---
id: S-0607
type: story
nature: improvement
title: stats.get runs its reads at once and its sections at once, and answers identical requests in flight once
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:41:34Z
updated: 2026-11-07T22:41:34Z
transitions: []
tags: [cli]
touches: [flai/internal/statsread/statsread.go, flai/internal/metrics/metrics.go, flai/internal/hostapi/hostapi.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0607 stats.get runs its reads at once and its sections at once, and answers identical requests in flight once

## Goal

Cause 17 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `stats.get` runs its 11 reads one after another (`flai/internal/statsread/statsread.go` `Read`), 230 to 300 ms warm, then its sections one after another (`flai/internal/metrics/metrics.go` `Compute`). It takes 0.63 to 1.1 s warm and 1.8 to 4 s cold, allocating about 200 MB a warm call. The dashboard asks again on every item change, and of the 120 slow calls in the serve log, 73 came within 3 s of the one before, each computed in full.

## Acceptance criteria

- [ ] `statsread.Read` runs its independent reads at once, and `Compute` runs its sections at once, with `claims` after `replayHolds`, which it reads; the report is the same as before, which a test over `flai/internal/metrics/testdata/good` shows, run under `-race`.
- [ ] Two `stats.get` requests with the same parameters while one is being computed get one computation's answer.
- [ ] A warm 30-day `stats.get` on this repository takes at least 250 ms less than before, and server-performance.md records the before and after per window.

## Tasks

## Notes

- Estimated from the profile: reads about 230 to 90 ms, the slowest read; compute about 450 to 280 ms at 30 days.
- The answer and `metrics.md` do not change.
