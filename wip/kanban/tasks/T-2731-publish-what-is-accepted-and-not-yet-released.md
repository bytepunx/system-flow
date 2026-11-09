---
id: T-2731
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: backlog
parent: S-0639
owner: alex
created: 2026-11-09T11:24:54Z
updated: 2026-11-09T11:24:54Z
transitions: []
stream: S-0639
tags: []
after: [T-2730]
---
# T-2731 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes
