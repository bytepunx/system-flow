---
id: T-3286
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0766
owner: alex
created: 2026-11-13T11:40:39Z
updated: 2026-11-13T11:42:38Z
transitions:
  - to: ready
    at: 2026-11-13T11:42:17Z
    by: agent-S-0766
  - to: in-progress
    at: 2026-11-13T11:42:17Z
    by: agent-S-0766
  - to: done
    at: 2026-11-13T11:42:38Z
    by: agent-S-0766
stream: S-0766
tags: []
after: [T-3285]
usage:
  source: log
  seconds: 21
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 737
      cache_read: 341362
      cache_write: 15395
      cost: 0.2062
---
# T-3286 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0766 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: ran as `scripts/flai.sh release --pending --require-pending --json` in the main checkout; exit 0, `pushed: true`.
- Released flai 1.88.0 → 1.89.0 (minor), tag `flai/v1.89.0` (on `origin` at 3b6f667), for S-0734 (minor). No other component was released; the template was not published.
