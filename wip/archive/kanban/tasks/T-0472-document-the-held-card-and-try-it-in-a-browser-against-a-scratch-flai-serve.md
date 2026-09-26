---
id: T-0472
type: task
nature: feature
title: Document the held card and try it in a browser against a scratch flai serve
status: done
parent: S-0129
owner: alex
created: 2026-09-26T17:54:08Z
updated: 2026-09-26T18:03:12Z
transitions:
  - to: ready
    at: 2026-09-26T17:54:12Z
    by: agent-S-0129
  - to: in-progress
    at: 2026-09-26T17:58:29Z
    by: agent-S-0129
  - to: done
    at: 2026-09-26T18:03:12Z
    by: agent-S-0129
stream: S-0129
tags: []
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/operators/index.md]
---
# T-0472 Document the held card and try it in a browser against a scratch flai serve

## Work

- Describe the held card and page in `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`, and `docs/operators/index.md`.
- Run `scripts/flaiover-test.sh`.
- Try it in a browser against a scratch project served by a scratch `flai serve` (never the operator's), with two stories whose touches overlap, one in progress and one in ready; then move the open one on and see the hold clear without a reload.

## Done when

- Docs describe it, flaiover's lint and tests pass, and the browser trial is recorded in the narrative.

## Notes
