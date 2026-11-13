---
id: T-3312
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: backlog
parent: S-0774
owner: alex
created: 2026-11-13T15:43:49Z
updated: 2026-11-13T15:43:49Z
transitions: []
stream: S-0774
tags: []
after: [T-3311]
---
# T-3312 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0774 --step publish --exit 2`, giving it the step's output.

## Notes
