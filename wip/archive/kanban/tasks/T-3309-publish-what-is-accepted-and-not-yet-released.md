---
id: T-3309
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0773
owner: alex
created: 2026-11-13T15:20:30Z
updated: 2026-11-13T15:23:49Z
transitions:
  - to: ready
    at: 2026-11-13T15:23:18Z
    by: agent-S-0773
  - to: in-progress
    at: 2026-11-13T15:23:18Z
    by: agent-S-0773
  - to: done
    at: 2026-11-13T15:23:49Z
    by: agent-S-0773
stream: S-0773
tags: []
after: [T-3308]
usage:
  source: log
  seconds: 31
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 1364
      cache_read: 669126
      cache_write: 26196
      cost: 0.3707
---
# T-3309 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0773 --step publish --exit 2`, giving it the step's output.

## Notes

Exited 0 in the main checkout, with `pushed: true`. It released S-0662:

| Component | From | To | Tag |
|-----------|------|----|-----|
| flai | 1.92.0 | 1.93.0 | `flai/v1.93.0` |
| flaiover | 0.58.8 | 0.58.9 | `flaiover/v0.58.9` |
| template | 1.1.37 | 1.1.38 | `v1.1.38` in `bytepunx/system-flow-template` (a7ea609) |
