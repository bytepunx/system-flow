---
id: T-2549
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: backlog
parent: S-0585
owner: alex
created: 2026-11-07T08:03:26Z
updated: 2026-11-07T08:03:26Z
transitions: []
stream: S-0585
tags: []
after: [T-2548]
---
# T-2549 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes
