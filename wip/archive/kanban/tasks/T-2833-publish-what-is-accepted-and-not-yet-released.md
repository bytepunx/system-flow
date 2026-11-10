---
id: T-2833
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0656
owner: alex
created: 2026-11-10T06:45:10Z
updated: 2026-11-10T07:02:20Z
transitions:
  - to: ready
    at: 2026-11-10T07:01:14Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T07:01:15Z
    by: system-flow
  - to: done
    at: 2026-11-10T07:02:20Z
    by: agent-S-0656
stream: S-0656
tags: []
after: [T-2832]
usage:
  source: log
  seconds: 65
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 563
      cache_read: 321653
      cache_write: 12253
      cost: 0.1736
---
# T-2833 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes

- `flai release --pending --json` exited 0 in the main checkout, with `pushed: true`. It released S-0625, S-0630, and S-0631: tags `flai/v1.86.9` and `flaiover/v0.58.3`, and the template at 1.1.31, published to the template repository as v1.1.31 (0f922e6).
