---
id: T-0543
type: task
nature: feature
title: The dashboard charts token rate, cost, and completion against time and cost by model
status: done
parent: S-0143
owner: alex
created: 2026-09-29T06:08:16Z
updated: 2026-09-29T06:31:35Z
transitions:
  - to: ready
    at: 2026-09-29T06:23:26Z
    by: agent-S-0143
  - to: in-progress
    at: 2026-09-29T06:23:26Z
    by: agent-S-0143
  - to: done
    at: 2026-09-29T06:31:35Z
    by: agent-S-0143
stream: S-0143
tags: []
touches: [flaiover/src, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
---
# T-0543 The dashboard charts token rate, cost, and completion against time and cost by model

## Work

- Charts, each grouping by model and selectable by type (epic, story, task): token rate per hour of agent work per item; cost per item, marking estimates; items done cumulatively against time and against cumulative cost.
- The item page shows an item's usage.

## Done when

- Component and chart tests pass (`pnpm test`), with `pnpm check` and lint clean.
- `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md` describe the charts.

## Notes
