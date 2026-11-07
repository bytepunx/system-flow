---
id: T-2547
type: task
nature: improvement
title: List the open issues the analyzer's report names
status: backlog
parent: S-0584
owner: alex
created: 2026-11-07T08:03:23Z
updated: 2026-11-07T08:03:23Z
transitions: []
stream: S-0584
tags: []
after: [T-2545]
---
# T-2547 List the open issues the analyzer's report names

## Work

Run `flai issue list --report latest --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai issue list --report latest --json` exits 0 in the main checkout.

## Notes
