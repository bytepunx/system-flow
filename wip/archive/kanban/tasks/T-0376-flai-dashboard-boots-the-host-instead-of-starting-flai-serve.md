---
id: T-0376
type: task
nature: feature
title: flai dashboard boots the host instead of starting flai serve
status: done
parent: S-0106
owner: alex
created: 2026-09-24T01:30:20Z
updated: 2026-09-24T01:42:35Z
transitions:
  - to: ready
    at: 2026-09-24T01:42:34Z
    by: system-flow
  - to: in-progress
    at: 2026-09-24T01:42:35Z
    by: system-flow
  - to: done
    at: 2026-09-24T01:42:35Z
    by: system-flow
stream: S-0106
tags: []
touches: [flai/cmd]
---
# T-0376 flai dashboard boots the host instead of starting flai serve

## Work

- `flai dashboard` (project and folder) registers as now, then makes sure the host runs rather than starting `flai serve`; `--no-serve` becomes "do not start the host".
- `flai dashboard status` and the start message report the host and its serve.

## Done when

- Dashboard tests pass with the host starter; no path in `flai dashboard` starts `flai serve` or `flai mcp` itself.

## Notes
