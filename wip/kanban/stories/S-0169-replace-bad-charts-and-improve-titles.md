---
id: S-0169
type: story
nature: remediation
title: Replace Bad Charts and Improve Titles
status: ready
owner: alex
created: 2026-09-29T23:50:42Z
updated: 2026-09-29T23:52:21Z
transitions:
  - to: ready
    at: 2026-09-29T23:50:44Z
    by: alex
tags: [dashboard, cli]
topics: [server-side, client-side]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0169 Replace Bad Charts and Improve Titles

## Goal

Several of the charts built for analyzing usage and flow do not actually show what was requested. Two replacements are proposed: 

"Avg. Time / Model" instead of "Completion over time"
"Avg. Cost / Model" instead of "Completion against cost"

Some additional title changes are requested and listed in acceptance criteria.

## Acceptance criteria
- [ ] Remove the "aging work in progress" and "estimate versus actual" charts from flow analysis
- [ ] Remove Completion over time and replace it with one that shows average time to complete a work type per model over time
- [ ] Remove Completion against cost and replace it with one that shows average dollars spent per model over time
- [ ] Change "Cycle time" to "Cycle Time"
- [ ] Change "Cumulative flow" to "Cumulative Flow"
- [ ] Change "Time in state" to "Time in State"
- [ ] Change "Token rate" to "Tokens / Min"
- [ ] Change "Tokens per day" to "Tokens / Day"
- [ ] Change "Tokens per dollar" to "Tokens / $"
- [ ] Change "Cost per day" to "$ / Day"
- [ ] Change "Cost per item" to "$ / Work Type"
- [ ] Change "Cost by item" to "$ / Item"

## Tasks

## Notes
