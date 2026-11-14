---
id: T-3363
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0787
owner: alex
created: 2026-11-14T18:00:18Z
updated: 2026-11-14T18:02:06Z
transitions:
  - to: ready
    at: 2026-11-14T18:02:05Z
    by: system-flow
  - to: in-progress
    at: 2026-11-14T18:02:06Z
    by: system-flow
stream: S-0787
tags: []
after: [T-3362]
---
# T-3363 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0787 --step publish --exit 2`, giving it the step's output.

## Notes
