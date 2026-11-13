---
id: T-3318
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0776
owner: alex
created: 2026-11-13T17:20:42Z
updated: 2026-11-13T17:22:22Z
transitions:
  - to: ready
    at: 2026-11-13T17:22:22Z
    by: agent-S-0776
  - to: in-progress
    at: 2026-11-13T17:22:22Z
    by: agent-S-0776
stream: S-0776
tags: []
after: [T-3317]
---
# T-3318 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0776 --step publish --exit 2`, giving it the step's output.

## Notes
