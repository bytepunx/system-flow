---
id: T-2833
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: backlog
parent: S-0656
owner: alex
created: 2026-11-10T06:45:10Z
updated: 2026-11-10T06:45:10Z
transitions: []
stream: S-0656
tags: []
after: [T-2832]
---
# T-2833 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes
