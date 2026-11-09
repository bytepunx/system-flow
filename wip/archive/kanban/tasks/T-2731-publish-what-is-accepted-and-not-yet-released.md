---
id: T-2731
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0639
owner: alex
created: 2026-11-09T11:24:54Z
updated: 2026-11-09T11:26:15Z
transitions:
  - to: ready
    at: 2026-11-09T11:26:03Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T11:26:03Z
    by: system-flow
  - to: done
    at: 2026-11-09T11:26:15Z
    by: agent-S-0639
stream: S-0639
tags: []
after: [T-2730]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 628
      cache_read: 299172
      cache_write: 12380
      cost: 0.1715
---
# T-2731 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes

- 2026-11-09: `flai release --pending --json` (as `scripts/flai.sh`) ran in the main checkout and exited 0: tagged `flai/v1.86.4` (patch from 1.86.3, S-0618), `pushed: true`, tags `["flai/v1.86.4"]`.
