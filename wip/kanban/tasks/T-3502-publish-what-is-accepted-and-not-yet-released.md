---
id: T-3502
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0820
owner: alex
created: 2026-11-17T02:48:43Z
updated: 2026-11-17T03:26:15Z
transitions:
  - to: ready
    at: 2026-11-17T03:26:15Z
    by: agent-S-0820
  - to: in-progress
    at: 2026-11-17T03:26:15Z
    by: agent-S-0820
stream: S-0820
tags: []
after: [T-3501]
---
# T-3502 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0820 --step publish --exit 2`, giving it the step's output.

## Notes
