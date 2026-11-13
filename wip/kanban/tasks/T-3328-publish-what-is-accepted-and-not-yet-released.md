---
id: T-3328
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0778
owner: alex
created: 2026-11-13T20:23:36Z
updated: 2026-11-13T20:26:20Z
transitions:
  - to: ready
    at: 2026-11-13T20:25:56Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T20:25:57Z
    by: system-flow
  - to: done
    at: 2026-11-13T20:26:20Z
    by: agent-S-0778
stream: S-0778
tags: []
after: [T-3327]
usage:
  source: log
  seconds: 23
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 267
      cache_read: 163002
      cache_write: 5701
      cost: 0.0836
---
# T-3328 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0778 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: `scripts/flai.sh release --pending --require-pending --json` in the main checkout exited 0 and pushed: `flai/v1.95.1` (patch from 1.95.0, S-0760) and the template 1.1.39 (patch from 1.1.38, S-0760) published to `bytepunx/system-flow-template` at 88a498e.
