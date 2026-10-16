---
id: T-1675
type: task
nature: feature
title: The workflow, the CLI and dashboard designs, and the user guides describe flai queue and the queue on the board
status: backlog
parent: S-0398
owner: alex
created: 2026-10-16T05:03:25Z
updated: 2026-10-16T05:03:25Z
transitions: []
stream: S-0398
tags: [cli, dashboard]
touches: [design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md]
after: [T-1638, T-1655, T-1667]
---
# T-1675 The workflow, the CLI and dashboard designs, and the user guides describe flai queue and the queue on the board

## Work

- `design/system/workflow.md` § What waits on an acceptance: the queue, its eligibility, and its order, linking the ADR.
- `design/system/flai-cli.md` and `docs/users/flai.md`: `flai queue` and the board's `queue`. Regenerate `flai-reference.md` with `make flai-reference`.
- `design/system/flaiover-dashboard.md` § Views and `docs/users/flaiover.md`: the review lane's marks and the queue panel.
- Waits for T-1638, T-1655, and T-1667, whose decision and behaviour it describes.

## Done when

The documents say the above, the reference is regenerated, and the markdown lint and `flai check --strict` pass.

## Notes
