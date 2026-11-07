---
id: T-2588
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0593
owner: alex
created: 2026-11-07T15:17:08Z
updated: 2026-11-07T15:19:04Z
transitions:
  - to: ready
    at: 2026-11-07T15:19:03Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T15:19:04Z
    by: system-flow
stream: S-0593
tags: []
after: [T-2587]
---
# T-2588 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes
