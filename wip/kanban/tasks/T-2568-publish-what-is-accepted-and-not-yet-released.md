---
id: T-2568
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: backlog
parent: S-0588
owner: alex
created: 2026-11-07T09:24:05Z
updated: 2026-11-07T09:24:05Z
transitions: []
stream: S-0588
tags: []
after: [T-2567]
---
# T-2568 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes
