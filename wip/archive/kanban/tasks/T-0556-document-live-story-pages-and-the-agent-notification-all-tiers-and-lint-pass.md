---
id: T-0556
type: task
nature: feature
title: "Document live story pages and the agent notification; all tiers and lint pass"
status: done
parent: S-0154
owner: alex
created: 2026-09-29T07:08:18Z
updated: 2026-09-29T07:21:29Z
transitions:
  - to: ready
    at: 2026-09-29T07:08:32Z
    by: agent-S-0154
  - to: in-progress
    at: 2026-09-29T07:20:35Z
    by: agent-S-0154
  - to: done
    at: 2026-09-29T07:21:29Z
    by: agent-S-0154
stream: S-0154
tags: []
usage:
  source: log
  seconds: 54
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 21
      output: 5838
      cache_read: 1489664
      cache_write: 19458
      cost: 0.5704
---

# T-0556 Document live story pages and the agent notification; all tiers and lint pass

## Work

Update `docs/users/flaiover.md` (the item page) and `design/system/flaiover-dashboard.md` (live updates), run `make test`, `make integration`, `make smoke`, flaiover's tests and lint, and `flai check --strict`.

## Done when

Docs say story pages update live and show the working line, and every tier and lint passes.

## Notes
