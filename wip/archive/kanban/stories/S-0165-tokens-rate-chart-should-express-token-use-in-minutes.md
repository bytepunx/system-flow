---
id: S-0165
type: story
nature: remediation
title: Tokens rate chart should express token use in minutes
status: done
parent: E-0013
owner: alex
created: 2026-09-29T22:37:02Z
updated: 2026-09-29T22:41:37Z
transitions:
  - to: ready
    at: 2026-09-29T22:37:06Z
    by: alex
  - to: in-progress
    at: 2026-09-29T22:41:03Z
    by: agent-S-0165
  - to: review
    at: 2026-09-29T22:41:14Z
    by: agent-S-0165
  - to: done
    at: 2026-09-29T22:41:37Z
    by: alex
tags: [dashboard, cli]
topics: [client-side, server-side]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 251
  models:
    - model: claude-opus-5-5
      input: 74
      output: 11746
      cache_read: 2368511
      cache_write: 76193
      cost: 1.3185
---
# S-0165 Tokens rate chart should express token use in minutes

## Goal

The token rate chart is showing tokens / agent hour for stories but this is a bad metric - most stories are completed in minutes, not hours, the chosen unit of time measure makes token expenditure seem much higher. change this to agent minutes, not hours.

## Acceptance criteria
- [x] token rate chart uses agent minutes, not agent hours

## Tasks
- T-0588 Confirm the token rate is charted per agent minute

## Notes

Met by S-0163 (ADR-0053 point 4) before this story was pulled: the chart plots `tokens_per_minute`. The operator's tab still ran the bundle from before the dashboard's restart; a hard refresh showed minutes (TH-0038). Nothing changed in code.
