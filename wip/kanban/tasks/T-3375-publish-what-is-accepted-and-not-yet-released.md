---
id: T-3375
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: backlog
parent: S-0791
owner: alex
created: 2026-11-15T07:18:49Z
updated: 2026-11-15T07:18:49Z
transitions: []
stream: S-0791
tags: []
after: [T-3374]
---
# T-3375 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0791 --step publish --exit 2`, giving it the step's output.

## Notes
