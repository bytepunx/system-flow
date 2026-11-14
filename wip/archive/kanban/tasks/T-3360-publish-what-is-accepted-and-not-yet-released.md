---
id: T-3360
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0786
owner: alex
created: 2026-11-14T17:32:28Z
updated: 2026-11-14T17:34:22Z
transitions:
  - to: ready
    at: 2026-11-14T17:34:05Z
    by: agent-S-0786
  - to: in-progress
    at: 2026-11-14T17:34:05Z
    by: agent-S-0786
  - to: done
    at: 2026-11-14T17:34:22Z
    by: agent-S-0786
stream: S-0786
tags: []
after: [T-3359]
usage:
  source: log
  seconds: 17
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 865
      cache_read: 496749
      cache_write: 19815
      cost: 0.2752
---
# T-3360 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0786 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-14: the command ran in the main checkout and exited 0, with `pushed: true`. It released S-0785 as `flai/v1.95.5` (patch, from 1.95.4) and `flaiover/v0.61.0`. Afterwards the main checkout's `main` was even with `origin/main`.
