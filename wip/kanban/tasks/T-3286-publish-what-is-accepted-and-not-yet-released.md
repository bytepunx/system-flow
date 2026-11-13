---
id: T-3286
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0766
owner: alex
created: 2026-11-13T11:40:39Z
updated: 2026-11-13T11:42:17Z
transitions:
  - to: ready
    at: 2026-11-13T11:42:17Z
    by: agent-S-0766
  - to: in-progress
    at: 2026-11-13T11:42:17Z
    by: agent-S-0766
stream: S-0766
tags: []
after: [T-3285]
---
# T-3286 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0766 --step publish --exit 2`, giving it the step's output.

## Notes
