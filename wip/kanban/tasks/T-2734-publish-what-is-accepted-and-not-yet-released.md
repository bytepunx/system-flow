---
id: T-2734
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: backlog
parent: S-0641
owner: alex
created: 2026-11-09T16:56:49Z
updated: 2026-11-09T16:56:49Z
transitions: []
stream: S-0641
tags: []
after: [T-2733]
---
# T-2734 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes
