---
id: T-2549
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: cancelled
parent: S-0585
owner: alex
created: 2026-11-07T08:03:26Z
updated: 2026-11-07T08:16:06Z
transitions:
  - to: cancelled
    at: 2026-11-07T08:16:06Z
    by: system-flow
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
- 2026-11-07T08:16:06Z: moved to cancelled: step 2 stopped the run: nothing accepted is waiting for a release
