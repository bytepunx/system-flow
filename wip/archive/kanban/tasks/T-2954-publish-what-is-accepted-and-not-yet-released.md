---
id: T-2954
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0704
owner: alex
created: 2026-11-10T08:56:47Z
updated: 2026-11-11T01:43:17Z
transitions:
  - to: ready
    at: 2026-11-11T01:42:53Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T01:42:53Z
    by: system-flow
  - to: done
    at: 2026-11-11T01:43:17Z
    by: agent-S-0704
stream: S-0704
tags: []
after: [T-2953]
usage:
  source: log
  seconds: 24
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 1338
      cache_read: 672666
      cache_write: 22756
      cost: 0.3434
---
# T-2954 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0704 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-11: ran as `scripts/flai.sh release --pending --require-pending --json` in the main checkout; exit 0, `pushed: true`.
- Released flai 1.86.13 → 1.86.14 (patch), tag `flai/v1.86.14`, bundling S-0665, S-0666, S-0675, S-0677, and S-0679, each a patch. No other component was released.
