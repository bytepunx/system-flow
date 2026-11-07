---
id: T-2591
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0594
owner: alex
created: 2026-11-07T21:47:12Z
updated: 2026-11-07T21:58:17Z
transitions:
  - to: ready
    at: 2026-11-07T21:58:10Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T21:58:17Z
    by: system-flow
stream: S-0594
tags: []
after: [T-2590]
---
# T-2591 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes
