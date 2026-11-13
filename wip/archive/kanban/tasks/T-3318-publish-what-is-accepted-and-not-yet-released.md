---
id: T-3318
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0776
owner: alex
created: 2026-11-13T17:20:42Z
updated: 2026-11-13T17:22:38Z
transitions:
  - to: ready
    at: 2026-11-13T17:22:22Z
    by: agent-S-0776
  - to: in-progress
    at: 2026-11-13T17:22:22Z
    by: agent-S-0776
  - to: done
    at: 2026-11-13T17:22:38Z
    by: agent-S-0776
stream: S-0776
tags: []
after: [T-3317]
usage:
  source: log
  seconds: 16
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 946
      cache_read: 482742
      cache_write: 20724
      cost: 0.2813
---
# T-3318 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0776 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: exited 0 in the main checkout (flai 1.86.17), `pushed: true`. Tagged and pushed `flai/v1.94.1` (patch, from 1.94.0) and `flaiover/v0.60.0` (minor, from 0.59.0), both for S-0757; `git ls-remote` shows both tags on origin.
