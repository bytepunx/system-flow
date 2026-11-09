---
id: T-2762
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: backlog
parent: S-0642
owner: alex
created: 2026-11-09T17:28:50Z
updated: 2026-11-09T17:28:50Z
transitions: []
stream: S-0642
tags: []
after: [T-2761]
---
# T-2762 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes
