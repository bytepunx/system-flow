---
id: T-0599
type: task
nature: remediation
title: Replace Completion over time with Avg. Time / Model
status: done
parent: S-0169
owner: alex
created: 2026-09-30T00:21:39Z
updated: 2026-09-30T00:36:16Z
transitions:
  - to: ready
    at: 2026-09-30T00:21:49Z
    by: agent-S-0169
  - to: in-progress
    at: 2026-09-30T00:29:36Z
    by: agent-S-0169
  - to: done
    at: 2026-09-30T00:36:16Z
    by: agent-S-0169
stream: S-0169
tags: []
touches: [flaiover/src, flai]
usage:
  source: log
  seconds: 400
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 69
      output: 20614
      cache_read: 5240879
      cache_write: 72235
      cost: 2.0386
---
# T-0599 Replace Completion over time with Avg. Time / Model

## Work

Replace the completion-time chart with Avg. Time / Model: one line per model per bucket, the mean time to complete an item of the chosen type, as TH-0040 settles which time.

## Done when

- Completion over time is gone and Avg. Time / Model draws on a time axis spanning the window.
- Unit and page tests cover it, and it is checked in a browser against this repository's stats.

## Notes
