---
id: E-0028
type: epic
nature: improvement
title: Charts improvements
status: done
owner: alex
created: 2026-10-21T20:40:56Z
updated: 2026-10-22T21:15:43Z
transitions:
  - to: ready
    at: 2026-10-21T20:41:04Z
    by: alex
  - to: in-progress
    at: 2026-10-22T20:17:45Z
    by: agent-S-0484
  - to: review
    at: 2026-10-22T21:11:53Z
    by: agent-S-0485
  - to: done
    at: 2026-10-22T21:15:43Z
    by: orchestrator
tags: [dashboard]
usage:
  source: sum
  seconds: 2970
  priced_by: harness
  turns:
    - day: 2026-10-22
      ceremony: 4
      test_runs: 2
      hand_edits: 6
      work: 102
  models:
    - model: claude-opus-5-5
      input: 578
      output: 212039
      cache_read: 33379536
      cache_write: 1235752
      cost: 18.5364
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
      seconds: 86070
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 303
          output: 4789
          cache_read: 50620481
          cache_write: 91456
          cost: 11.198
        - model: claude-sonnet-5-5
          input: 30
          output: 192
          cache_read: 604441
          cache_write: 138989
          cost: 0.5145
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
- 2026-10-22T21:11:53Z: moved to review: follows S-0485, which moved to review
- 2026-10-22T21:15:43Z: moved to done: follows S-0485, which moved to done
