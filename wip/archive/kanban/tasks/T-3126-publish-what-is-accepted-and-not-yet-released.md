---
id: T-3126
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0727
owner: alex
created: 2026-11-11T23:22:51Z
updated: 2026-11-11T23:24:54Z
transitions:
  - to: ready
    at: 2026-11-11T23:24:38Z
    by: agent-S-0727
  - to: in-progress
    at: 2026-11-11T23:24:38Z
    by: agent-S-0727
  - to: done
    at: 2026-11-11T23:24:54Z
    by: agent-S-0727
stream: S-0727
tags: []
after: [T-3125]
usage:
  source: log
  seconds: 16
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 571
      cache_read: 317304
      cache_write: 12329
      cost: 0.1735
---
# T-3126 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0727 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-11: exited 0 in the main checkout, `pushed: true`. Tagged and pushed `flai/v1.86.19` (patch from 1.86.18) and `flaiover/v0.58.8` (patch from 0.58.7), both for S-0720.
