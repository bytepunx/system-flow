---
id: T-0431
type: task
nature: improvement
title: Docs and design say what Retry does
status: done
parent: S-0118
owner: arobson
created: 2026-09-26T03:02:39Z
updated: 2026-09-26T03:14:40Z
transitions:
  - to: ready
    at: 2026-09-26T03:02:43Z
    by: agent-S-0118
  - to: in-progress
    at: 2026-09-26T03:13:33Z
    by: agent-S-0118
  - to: done
    at: 2026-09-26T03:14:40Z
    by: agent-S-0118
stream: S-0118
tags: []
touches: [design/system, docs]
---
# T-0431 Docs and design say what Retry does

## Work

Update these to describe Retry and queueing:

- `design/system/flaiover-dashboard.md`
- `design/system/flai-cli.md`
- `docs/users/flaiover.md`
- `docs/users/flai-reference.md` (regenerate it if it is generated)
- `docs/operators/index.md`

## Done when

- No document still names a Restart agent button.
- Every document above says what Retry does when the in-progress limit is full.
- `flai check --strict` is clean.

## Notes
