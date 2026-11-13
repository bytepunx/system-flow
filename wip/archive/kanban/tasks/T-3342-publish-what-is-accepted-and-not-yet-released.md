---
id: T-3342
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: cancelled
parent: S-0783
owner: alex
created: 2026-11-13T22:47:23Z
updated: 2026-11-13T22:49:25Z
transitions:
  - to: cancelled
    at: 2026-11-13T22:49:25Z
    by: system-flow
stream: S-0783
tags: []
after: [T-3341]
---
# T-3342 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0783 --step publish --exit 2`, giving it the step's output.

## Notes
- 2026-11-13T22:49:25Z: moved to cancelled: step evaluate stopped the run with exit 2
