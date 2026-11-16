---
id: T-3389
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0794
owner: alex
created: 2026-11-15T21:47:08Z
updated: 2026-11-16T00:38:46Z
transitions:
  - to: ready
    at: 2026-11-16T00:38:46Z
    by: agent-S-0794
  - to: in-progress
    at: 2026-11-16T00:38:46Z
    by: agent-S-0794
stream: S-0794
tags: []
after: [T-3388]
---
# T-3389 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0794 --step publish --exit 2`, giving it the step's output.

## Notes
