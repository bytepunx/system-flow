---
id: T-3519
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0826
owner: alex
created: 2026-11-17T08:16:46Z
updated: 2026-11-17T08:19:14Z
transitions:
  - to: ready
    at: 2026-11-17T08:19:14Z
    by: agent-S-0826
  - to: in-progress
    at: 2026-11-17T08:19:14Z
    by: agent-S-0826
stream: S-0826
tags: []
after: [T-3518]
---
# T-3519 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0826 --step publish --exit 2`, giving it the step's output.

## Notes
