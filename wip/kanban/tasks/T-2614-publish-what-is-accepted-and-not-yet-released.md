---
id: T-2614
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: backlog
parent: S-0597
owner: alex
created: 2026-11-07T22:31:23Z
updated: 2026-11-07T22:31:23Z
transitions: []
stream: S-0597
tags: []
after: [T-2613]
---
# T-2614 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes
