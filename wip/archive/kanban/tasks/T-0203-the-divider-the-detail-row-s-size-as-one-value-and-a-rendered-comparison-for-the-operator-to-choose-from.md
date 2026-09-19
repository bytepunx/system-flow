---
id: T-0203
type: task
nature: improvement
title: The divider, the detail row's size as one value, and a rendered comparison for the operator to choose from
status: done
parent: S-0054
owner: alex
created: 2026-09-19T05:44:37Z
updated: 2026-09-19T05:51:22Z
transitions:
  - to: ready
    at: 2026-09-19T05:45:03Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T05:45:03Z
    by: system-flow
  - to: done
    at: 2026-09-19T05:51:22Z
    by: system-flow
stream: S-0054
tags: []
touches: [flaiover/src/lib/components]
---

# T-0203 The divider, the detail row's size as one value, and a rendered comparison for the operator to choose from

## Work
In `BoardCard.svelte`, add the thin divider between the title and the bottom row: a top border on the row in the theme's line colour, with the same space above and below it. Put the row's font size in one place so every item in it, the parent ID included, follows it. Render the board's cards at three candidate sizes (11, 12, and 13 pixels; today is 10 and the title is 12) in the light and the dark theme, from a real page and not a mock-up, put the comparison in front of the operator, and ask which they want, and whether the top row (ID and age) should follow. Do not pick for them: the criterion says the size is chosen with the operator.

## Done when
- The divider is in the component in both themes
- The comparison image was sent to the operator and their choice is recorded in the narrative's Decisions

## Notes
