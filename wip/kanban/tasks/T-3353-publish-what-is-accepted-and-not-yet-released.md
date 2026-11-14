---
id: T-3353
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0784
owner: alex
created: 2026-11-14T07:05:35Z
updated: 2026-11-14T07:07:49Z
transitions:
  - to: ready
    at: 2026-11-14T07:07:48Z
    by: agent-S-0784
  - to: in-progress
    at: 2026-11-14T07:07:49Z
    by: agent-S-0784
stream: S-0784
tags: []
after: [T-3352]
---
# T-3353 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0784 --step publish --exit 2`, giving it the step's output.

## Notes
