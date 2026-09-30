---
id: T-0600
type: task
nature: remediation
title: Replace Completion against cost with Avg. Cost / Model
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
    at: 2026-09-30T00:24:52Z
    by: agent-S-0169
  - to: done
    at: 2026-09-30T00:36:16Z
    by: agent-S-0169
stream: S-0169
tags: []
touches: [flaiover/src]
usage:
  source: log
  seconds: 684
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 110
      output: 33011
      cache_read: 8392574
      cache_write: 115675
      cost: 3.2646
---
# T-0600 Replace Completion against cost with Avg. Cost / Model

## Work

Replace the completion-cost chart with Avg. Cost / Model: one line per model per bucket, the mean dollars per item of the chosen type, and settle the compare control as TH-0040 answers.

## Done when

- Completion against cost is gone and Avg. Cost / Model draws on a time axis spanning the window.
- Unit and page tests cover it, and it is checked in a browser.

## Notes
