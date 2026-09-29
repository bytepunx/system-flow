---
id: E-0012
type: epic
nature: improvement
title: Performance Analysis and Improvements
status: in-progress
owner: alex
created: 2026-09-29T05:26:40Z
updated: 2026-09-29T07:07:30Z
transitions:
  - to: ready
    at: 2026-09-29T07:07:28Z
    by: alex
  - to: in-progress
    at: 2026-09-29T07:07:30Z
    by: alex
tags: [dashboard, cli]
topics: [front-end, back-end]
touches: [flaiover/src, flai/cmd]
usage:
  source: sum
  seconds: 4249
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 918
      output: 290747
      cache_read: 63651208
      cache_write: 1025050
      cost: 26.7599
---
# E-0012 Performance Analysis and Improvements

## Outcome

Several interactions (loading the board) have become slower over time. We need to look into performance analysis to determine where these issues are occurring, file them and address them.

## Stories
- S-0152 Identify issues with responses from the server side
- S-0156 flai serve and flai mcp keep parsed work items between requests and read again only the files that changed
- S-0157 Pending releases are worked out without starting a git process per accepted item
- S-0158 The designer's inbox finds overlapping touches without running the whole check
- S-0159 The dashboard's reads are answered in flai serve's process, not by starting flai
- S-0160 A flai process starts in milliseconds whatever the host's PATH
- S-0161 The dashboard forgets only the answers a changed file affects, and gathers changes that arrive together
- S-0162 The items and documents pages ask for what they show, not the whole archive and every file's front matter

## Notes
