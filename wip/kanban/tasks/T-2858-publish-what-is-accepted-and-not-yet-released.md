---
id: T-2858
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: backlog
parent: S-0658
owner: alex
created: 2026-11-10T07:56:53Z
updated: 2026-11-10T07:56:53Z
transitions: []
stream: S-0658
tags: []
after: [T-2857]
---
# T-2858 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0658 --step publish --exit 2`, giving it the step's output.

## Notes
