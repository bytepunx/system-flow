---
id: T-3512
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: cancelled
parent: S-0821
owner: alex
created: 2026-11-17T04:29:08Z
updated: 2026-11-17T07:46:43Z
transitions:
  - to: cancelled
    at: 2026-11-17T07:46:43Z
    by: alex
stream: S-0821
tags: []
after: [T-3511]
---
# T-3512 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0821 --step publish --exit 2`, giving it the step's output.

## Notes
- 2026-11-17T07:46:43Z: moved to cancelled: S-0821 cancelled: not needed
