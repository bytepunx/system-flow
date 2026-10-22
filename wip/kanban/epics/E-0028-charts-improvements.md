---
id: E-0028
type: epic
nature: improvement
title: Charts improvements
status: in-progress
owner: alex
created: 2026-10-21T20:40:56Z
updated: 2026-10-22T20:17:45Z
transitions:
  - to: ready
    at: 2026-10-21T20:41:04Z
    by: alex
  - to: in-progress
    at: 2026-10-22T20:17:45Z
    by: agent-S-0484
tags: [dashboard]
usage:
  source: sum
  seconds: 1681
  priced_by: harness
  turns:
    - day: 2026-10-22
      ceremony: 2
      test_runs: 2
      hand_edits: 4
      work: 54
  models:
    - model: claude-opus-5-5
      input: 332
      output: 103932
      cache_read: 16266309
      cache_write: 583183
      cost: 8.8835
  strategic:
    - kind: planner
      seconds: 84620
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 2338
          output: 57932
          cache_read: 18779377
          cache_write: 314269
          cost: 13.9002
        - model: claude-haiku-5-5
          input: 64
          output: 32014
          cache_read: 2735337
          cache_write: 155410
          cost: 0.197
    - kind: orchestrator
      seconds: 85446
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 220
          output: 3583
          cache_read: 33067638
          cache_write: 47919
          cost: 7.3125
        - model: claude-sonnet-5-5
          input: 14
          output: 85
          cache_read: 255175
          cache_write: 65016
          cost: 0.2216
---
# E-0028 Charts improvements

## Outcome

At the moment, each type of chart has a label and then a row of chart names out to the side (with some rows overflowing). Filters at the top of the page give the operator ways to change what is displayed.

Instead, I would like the following changes:

- The charts sub-menu item should expand when clicked to reveal the list of charts, subdivided by type in a vertical list. Clicking each chart is what displays it on the chart page.
- The title and chart should be at the top of the chart page, followed by a collapsible paragraph explaining what the chart shows and helpful details (i.e. not just what Cycle Time is but include how that time is measured) with filtering controls beneath the explainer, followed by the expandable table last.

## Stories
- S-0484 The rail's Charts item expands to list the charts by group, each a link that opens it, and the charts page keeps no chart list of its own
- S-0485 The chart page shows the title and the chart first, then a collapsible explainer of what the chart shows and how it is measured, then the controls, then the table view

## Notes
- 2026-10-22T20:17:45Z: moved to in-progress: follows S-0484, which moved to in-progress
