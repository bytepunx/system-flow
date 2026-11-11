---
id: T-3129
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0728
owner: alex
created: 2026-11-11T23:40:03Z
updated: 2026-11-11T23:42:02Z
transitions:
  - to: ready
    at: 2026-11-11T23:41:50Z
    by: agent-S-0728
  - to: in-progress
    at: 2026-11-11T23:41:50Z
    by: agent-S-0728
  - to: done
    at: 2026-11-11T23:42:02Z
    by: agent-S-0728
stream: S-0728
tags: []
after: [T-3128]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 284
      cache_read: 158295
      cache_write: 6382
      cost: 0.0884
---
# T-3129 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0728 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-11: exited 0 in the main checkout, `pushed: true`. Tagged and pushed `flai/v1.86.20` (patch from 1.86.19) for S-0703.
