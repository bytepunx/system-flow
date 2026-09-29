---
id: E-0014
type: epic
nature: feature
title: Build real, collapsible diffs in the review screen
status: in-progress
owner: alex
created: 2026-09-29T19:57:06Z
updated: 2026-09-29T21:10:42Z
transitions:
  - to: ready
    at: 2026-09-29T21:02:58Z
    by: alex
  - to: in-progress
    at: 2026-09-29T21:05:12Z
    by: alex
tags: [dashboard, cli]
topics: [client-side, server-side]
touches: [flaiover/src, flai/cmd]
usage:
  source: sum
  seconds: 741
  models:
    - model: claude-fable-5-1
      input: 96
      output: 39175
      cache_read: 5516982
      cache_write: 171726
      cost: 6.7735
---
# E-0014 Build real, collapsible diffs in the review screen

## Outcome

The operator should be able to see clear diffs in the review pane by clicking the summary line for each file (currently shows the lines added, lines subtracted) to expand or collapse the git diff.

## Stories
- S-0164 Reviews need real diffs, not just line counts

## Notes
