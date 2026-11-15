---
id: T-3375
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: done
parent: S-0791
owner: alex
created: 2026-11-15T07:18:49Z
updated: 2026-11-15T21:01:26Z
transitions:
  - to: ready
    at: 2026-11-15T21:00:58Z
    by: agent-S-0791
  - to: in-progress
    at: 2026-11-15T21:00:58Z
    by: agent-S-0791
  - to: done
    at: 2026-11-15T21:01:26Z
    by: agent-S-0791
stream: S-0791
tags: []
after: [T-3374]
usage:
  source: log
  seconds: 28
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1060
      cache_read: 495075
      cache_write: 20500
      cost: 0.2842
---
# T-3375 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0791 --step publish --exit 2`, giving it the step's output.

## Notes

- 2026-11-15: exit 0. The release published flai 1.97.0 → 1.98.0, a minor bump, as the tag `flai/v1.98.0`, with S-0743, S-0744, and S-0749. The tags and the branch were pushed (`pushed: true`).
