---
id: T-3290
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: cancelled
parent: S-0768
owner: alex
created: 2026-11-13T11:43:32Z
updated: 2026-11-13T11:45:33Z
transitions:
  - to: ready
    at: 2026-11-13T11:45:13Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T11:45:13Z
    by: system-flow
  - to: cancelled
    at: 2026-11-13T11:45:33Z
    by: alex
stream: S-0768
tags: []
after: [T-3289]
usage:
  source: log
  seconds: 20
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 510
      cache_read: 167528
      cache_write: 10432
      cost: 0.1272
---
# T-3290 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0768 --step publish --exit 2`, giving it the step's output.

## Notes
- 2026-11-13T11:45:33Z: moved to cancelled: S-0768 cancelled: REDUNDANT
