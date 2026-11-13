---
id: T-3299
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: backlog
parent: S-0770
owner: alex
created: 2026-11-13T13:16:45Z
updated: 2026-11-13T13:16:45Z
transitions: []
stream: S-0770
tags: []
after: [T-3298]
---
# T-3299 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0770 --step publish --exit 2`, giving it the step's output.

## Notes
