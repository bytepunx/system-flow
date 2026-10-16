---
id: T-1680
type: task
nature: feature
title: The dashboard design, the CLI design, and the user guides describe the queue's figures and chart
status: backlog
parent: S-0404
owner: alex
created: 2026-10-16T05:03:45Z
updated: 2026-10-16T05:03:45Z
transitions: []
stream: S-0404
tags: [cli, dashboard]
touches: [design/system/flaiover-dashboard.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md]
after: [T-1673]
---
# T-1680 The dashboard design, the CLI design, and the user guides describe the queue's figures and chart

## Work

- `design/system/flai-cli.md` and `docs/users/flai.md`: `flai stats`'s `queue` section and lines.
- `design/system/flaiover-dashboard.md` § Views and `docs/users/flaiover.md`: the queue chart.
- Waits for T-1673, whose chart it describes.

## Done when

The documents say the above, and the markdown lint and `flai check --strict` pass.

## Notes
