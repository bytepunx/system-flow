---
id: T-1675
type: task
nature: feature
title: The workflow, the CLI and dashboard designs, and the user guides describe flai queue and the queue on the board
status: done
parent: S-0398
owner: alex
created: 2026-10-16T05:03:25Z
updated: 2026-10-19T15:28:18Z
transitions:
  - to: ready
    at: 2026-10-19T15:23:11Z
    by: agent-S-0398
  - to: in-progress
    at: 2026-10-19T15:23:12Z
    by: agent-S-0398
  - to: done
    at: 2026-10-19T15:28:18Z
    by: agent-S-0398
stream: S-0398
tags: [cli, dashboard]
touches: [design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md]
after: [T-1655, T-1667]
usage:
  source: log
  seconds: 306
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 79
      output: 26191
      cache_read: 3992579
      cache_write: 230809
      cost: 2.7703
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
