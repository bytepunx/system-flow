---
id: T-1680
type: task
nature: feature
title: The dashboard design, the CLI design, and the user guides describe the queue's figures and chart
status: done
parent: S-0404
owner: alex
created: 2026-10-16T05:03:45Z
updated: 2026-10-20T00:34:36Z
transitions:
  - to: ready
    at: 2026-10-20T00:31:04Z
    by: agent-S-0404
  - to: in-progress
    at: 2026-10-20T00:31:05Z
    by: agent-S-0404
  - to: done
    at: 2026-10-20T00:34:36Z
    by: agent-S-0404
stream: S-0404
tags: [cli, dashboard]
touches: [design/system/flaiover-dashboard.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, design/system/metrics.md]
after: [T-1673]
usage:
  source: log
  seconds: 211
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 49
      output: 17791
      cache_read: 2993799
      cache_write: 101117
      cost: 1.5236
---
# T-1680 The dashboard design, the CLI design, and the user guides describe the queue's figures and chart

## Work

- `design/system/flai-cli.md` and `docs/users/flai.md`: `flai stats`'s `queue` section and lines.
- `design/system/flaiover-dashboard.md` § Views and `docs/users/flaiover.md`: the queue chart.
- Waits for T-1673, whose chart it describes.

## Done when

The documents say the above, and the markdown lint and `flai check --strict` pass.

## Notes
