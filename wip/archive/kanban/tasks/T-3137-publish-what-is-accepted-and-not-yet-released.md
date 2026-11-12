---
id: T-3137
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0731
owner: alex
created: 2026-11-12T03:44:58Z
updated: 2026-11-12T03:47:24Z
transitions:
  - to: ready
    at: 2026-11-12T03:47:08Z
    by: agent-S-0731
  - to: in-progress
    at: 2026-11-12T03:47:08Z
    by: agent-S-0731
  - to: done
    at: 2026-11-12T03:47:24Z
    by: agent-S-0731
stream: S-0731
tags: []
after: [T-3136]
usage:
  source: log
  seconds: 16
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 661
      cache_read: 320452
      cache_write: 12981
      cost: 0.1812
---
# T-3137 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0731 --step publish --exit 2`, giving it the step's output.

## Notes

`flai release --pending --require-pending --json` exited 0 in the main checkout (2026-11-12), with flai 1.86.17: it tagged `flai/v1.86.22` (patch, from 1.86.21) for S-0717, S-0701, and S-0681, and pushed the branch and tag (`pushed: true`).
