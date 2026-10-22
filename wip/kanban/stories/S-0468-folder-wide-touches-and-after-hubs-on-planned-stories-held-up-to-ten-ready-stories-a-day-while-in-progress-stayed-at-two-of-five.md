---
id: S-0468
type: story
nature: improvement
title: Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five
status: backlog
owner: alex
created: 2026-10-21T11:58:08Z
updated: 2026-10-22T20:51:09Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: analyzer
      seconds: 113
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 97
          output: 6305
          cache_read: 370755
          cache_write: 41688
          cost: 1.2427
        - model: claude-haiku-5-5
          input: 8
          output: 6718
          cache_read: 204809
          cache_write: 23259
          cost: 0.0083
cost_of_delay:
  inputs:
    time_lost_per_cycle: 44h
    by: flai
    at: 2026-10-21T11:58:08Z
finalized:
  by: alex
  at: 2026-10-22T20:51:09Z
---
# S-0468 Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five

## Goal

This story remediates [I-0152](../../../design/issues/I-0152-folder-wide-touches-and-after-hubs-on-planned-stories-held-up-to-ten-ready-stories-a-day-while-in-progress-stayed-at-two-of-five.md), "Folder-wide touches and after hubs on planned stories held up to ten ready stories a day while in-progress stayed at two of five", which the analysis in [2026-10-15-bottlenecks.md](../../../design/analysis/2026-10-15-bottlenecks.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0152 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0152 is closed with `flai issue close I-0152 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0152: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0152. time_lost_per_cycle 44h carried over from I-0152's Impact section. Its Impact time_lost_per_cycle was taken rather than the 44h derived from its cost and count.

Strategic usage carried over from I-0152: analyzer 1m53s, 1.2510 USD.
