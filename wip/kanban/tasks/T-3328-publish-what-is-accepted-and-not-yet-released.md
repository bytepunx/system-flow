---
id: T-3328
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0778
owner: alex
created: 2026-11-13T20:23:36Z
updated: 2026-11-13T20:25:57Z
transitions:
  - to: ready
    at: 2026-11-13T20:25:56Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T20:25:57Z
    by: system-flow
stream: S-0778
tags: []
after: [T-3327]
---
# T-3328 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0778 --step publish --exit 2`, giving it the step's output.

## Notes
