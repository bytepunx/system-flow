---
id: T-2824
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: cancelled
parent: S-0646
owner: alex
created: 2026-11-09T22:22:55Z
updated: 2026-11-10T06:18:26Z
transitions:
  - to: cancelled
    at: 2026-11-10T06:18:26Z
    by: system-flow
stream: S-0646
tags: []
after: [T-2823]
---
# T-2824 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes
- 2026-11-10T06:18:26Z: moved to cancelled: step evaluate stopped the run with exit 2
