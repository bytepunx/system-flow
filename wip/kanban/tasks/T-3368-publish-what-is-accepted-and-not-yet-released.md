---
id: T-3368
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: backlog
parent: S-0789
owner: alex
created: 2026-11-14T23:49:29Z
updated: 2026-11-14T23:49:29Z
transitions: []
stream: S-0789
tags: []
after: [T-3367]
---
# T-3368 Publish what is accepted and not yet released

## Work

Run `scripts/flai.sh release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0789 --step publish --exit 2`, giving it the step's output.

## Notes
