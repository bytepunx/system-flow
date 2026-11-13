---
id: T-3315
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0775
owner: alex
created: 2026-11-13T16:13:36Z
updated: 2026-11-13T16:15:41Z
transitions:
  - to: ready
    at: 2026-11-13T16:15:26Z
    by: agent-S-0775
  - to: in-progress
    at: 2026-11-13T16:15:26Z
    by: agent-S-0775
  - to: done
    at: 2026-11-13T16:15:41Z
    by: agent-S-0775
stream: S-0775
tags: []
after: [T-3314]
usage:
  source: log
  seconds: 15
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 587
      cache_read: 326372
      cache_write: 12786
      cost: 0.1793
---
# T-3315 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0775 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: exited 0 in the main checkout. It tagged `flai/v1.94.0` (minor, from 1.93.1, for S-0751) and pushed the branch and the tag (`pushed: true`).
