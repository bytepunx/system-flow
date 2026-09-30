---
id: S-0169
type: story
nature: remediation
title: Replace Bad Charts and Improve Titles
status: done
owner: alex
created: 2026-09-29T23:50:42Z
updated: 2026-09-30T00:46:13Z
transitions:
  - to: ready
    at: 2026-09-29T23:50:44Z
    by: alex
  - to: in-progress
    at: 2026-09-30T00:20:13Z
    by: agent-S-0169
  - to: review
    at: 2026-09-30T00:38:38Z
    by: agent-S-0169
  - to: done
    at: 2026-09-30T00:46:13Z
    by: alex
tags: [dashboard, cli]
topics: [server-side, client-side]
touches: [flaiover/src, flai/cmd, flai/internal/metrics, design/system, design/adrs, docs/users]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1153
  models:
    - model: claude-opus-5-5
      input: 216
      output: 64617
      cache_read: 16427809
      cache_write: 226424
      cost: 6.3902
---
# S-0169 Replace Bad Charts and Improve Titles

## Goal

Several of the charts built for analyzing usage and flow do not actually show what was requested. Two replacements are proposed:

"Avg. Time / Model" instead of "Completion over time"
"Avg. Cost / Model" instead of "Completion against cost"

Some additional title changes are requested and listed in acceptance criteria.

## Acceptance criteria
- [x] Remove the "aging work in progress" and "estimate versus actual" charts from flow analysis
- [x] Remove Completion over time and replace it with one that shows average time to complete a work type per model over time
- [x] Remove Completion against cost and replace it with one that shows average dollars spent per model over time
- [x] Change "Cycle time" to "Cycle Time"
- [x] Change "Cumulative flow" to "Cumulative Flow"
- [x] Change "Time in state" to "Time in State"
- [x] Change "Token rate" to "Tokens / Min"
- [x] Change "Tokens per day" to "Tokens / Day"
- [x] Change "Tokens per dollar" to "Tokens / $"
- [x] Change "Cost per day" to "$ / Day"
- [x] Change "Cost per item" to "$ / Work Type"
- [x] Change "Cost by item" to "$ / Item"

## Tasks
- T-0598 Remove aging and estimates, and retitle the charts
- T-0599 Replace Completion over time with Avg. Time / Model
- T-0600 Replace Completion against cost with Avg. Cost / Model
- T-0601 Record the chart changes in an ADR and the design and user docs

## Notes

- Avg. Time / Model plots agent minutes per item, and Avg. Cost / Model replaces the compare control on Tokens per item and $ / Work Type (TH-0040, ADR-0057).
- Tokens / Day and $ / Day follow the bucket chosen: Tokens / Hour, $ / Week.
