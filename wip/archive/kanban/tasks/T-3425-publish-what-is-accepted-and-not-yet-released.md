---
id: T-3425
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0802
owner: alex
created: 2026-11-16T04:12:37Z
updated: 2026-11-16T04:14:30Z
transitions:
  - to: ready
    at: 2026-11-16T04:14:14Z
    by: agent-S-0802
  - to: in-progress
    at: 2026-11-16T04:14:14Z
    by: agent-S-0802
  - to: done
    at: 2026-11-16T04:14:30Z
    by: agent-S-0802
stream: S-0802
tags: []
after: [T-3424]
usage:
  source: log
  seconds: 16
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 951
      cache_read: 495870
      cache_write: 21452
      cost: 0.2898
---
# T-3425 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0802 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-16: exited 0 in the main checkout, `pushed: true`. Released S-0799 as patches: flai 1.101.1 → 1.101.2 (`flai/v1.101.2`), flaiover 0.61.1 → 0.61.2 (`flaiover/v0.61.2`).
