---
id: T-2091
type: task
nature: feature
title: docs/users/flai-reference.md is regenerated for flai dashboard's help on the allowed hosts
status: done
parent: S-0489
owner: alex
created: 2026-10-23T19:50:59Z
updated: 2026-10-23T19:51:42Z
transitions:
  - to: ready
    at: 2026-10-23T19:51:13Z
    by: agent-S-0489
  - to: in-progress
    at: 2026-10-23T19:51:14Z
    by: agent-S-0489
  - to: done
    at: 2026-10-23T19:51:42Z
    by: agent-S-0489
stream: S-0489
tags: []
touches: [docs/users/flai-reference.md]
usage:
  source: log
  seconds: 28
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 4357
      cache_read: 907713
      cache_write: 25351
      cost: 0.42
---

# T-2091 docs/users/flai-reference.md is regenerated for flai dashboard's help on the allowed hosts

## Work

The merge queue's batch 20261023-194240 failed S-0489 at `integration` (TH-0625). `TestReferenceIsCurrent` found `docs/users/flai-reference.md` stale. T-2084 added `dashboard.allowed_hosts` to `config.Keys` and changed the long help of `flai dashboard restart`, `status`, and `upgrade`, and the reference is generated from both. Regenerate it with `scripts/flai-reference.sh` after the sync.

This task waits for no other: every other task of S-0489 is done.

## Done when

- `docs/users/flai-reference.md` is regenerated, and `TestReferenceIsCurrent` passes.
- The close-out passes.

## Notes
