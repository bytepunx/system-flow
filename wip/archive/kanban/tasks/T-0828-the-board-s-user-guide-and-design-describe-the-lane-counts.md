---
id: T-0828
type: task
nature: improvement
title: The board's user guide and design describe the lane counts
status: done
parent: S-0256
owner: alex
created: 2026-10-04T21:24:02Z
updated: 2026-10-04T21:27:27Z
transitions:
  - to: ready
    at: 2026-10-04T21:27:07Z
    by: agent-S-0256
  - to: in-progress
    at: 2026-10-04T21:27:07Z
    by: agent-S-0256
  - to: done
    at: 2026-10-04T21:27:27Z
    by: agent-S-0256
stream: S-0256
tags: []
touches: [docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [T-0827]
---
# T-0828 The board's user guide and design describe the lane counts

## Work

Describe the lane counts and their tooltip in the board section of `docs/users/flaiover.md` and the `/board` row of `design/system/flaiover-dashboard.md`, naming the function in `$lib/lanes.ts`. Waits for T-0827 so that the documents describe what was built.

## Done when

- [x] Both documents say what the count shows, what it counts, and what the tooltip says, with `updated` bumped.

## Notes
