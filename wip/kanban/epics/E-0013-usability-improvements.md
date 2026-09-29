---
id: E-0013
type: epic
nature: improvement
title: Usability Improvements
status: in-progress
owner: alex
created: 2026-09-29T05:28:11Z
updated: 2026-09-29T23:43:11Z
transitions:
  - to: ready
    at: 2026-09-29T07:07:18Z
    by: alex
  - to: in-progress
    at: 2026-09-29T07:07:20Z
    by: alex
tags: [dashboard]
topics: [front-end]
touches: [flaiover/src]
usage:
  source: sum
  seconds: 4306
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 1160
      output: 242648
      cache_read: 70815780
      cache_write: 1228464
      cost: 29.9776
---
# E-0013 Usability Improvements

## Outcome

There are a number of odd patterns that cause confusion in the dashboard UI/UX. We need to catalog and address these to improve the usability and utility of the system.

## Stories
- S-0150 Agent status and updates on board need to be collapsible
- S-0151 Temporary banner notifications should be dismissible
- S-0153 Threads in the story pages are supposed to be paged
- S-0154 Story pages receive live updates
- S-0155 Inbox items take users to the story page, not document pages
- S-0165 Tokens rate chart should express token use in minutes
- S-0166 Charts ignore the window drop-down
- S-0167 Right click menu for the board
- S-0168 Several Charts Don't Use Window correctly

## Notes
