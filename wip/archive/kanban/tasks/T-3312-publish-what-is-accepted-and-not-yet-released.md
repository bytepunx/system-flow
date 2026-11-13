---
id: T-3312
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0774
owner: alex
created: 2026-11-13T15:43:49Z
updated: 2026-11-13T15:45:38Z
transitions:
  - to: ready
    at: 2026-11-13T15:45:14Z
    by: agent-S-0774
  - to: in-progress
    at: 2026-11-13T15:45:14Z
    by: agent-S-0774
  - to: done
    at: 2026-11-13T15:45:38Z
    by: agent-S-0774
stream: S-0774
tags: []
after: [T-3311]
usage:
  source: log
  seconds: 24
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 864
      cache_read: 488506
      cache_write: 20087
      cost: 0.2757
---
# T-3312 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0774 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: exited 0 in the main checkout (flai 1.86.17), `pushed: true`. Tagged and pushed `flai/v1.93.1` (patch, from 1.93.0) and `flaiover/v0.59.0` (minor, from 0.58.9), both for S-0663.
