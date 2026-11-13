---
id: T-3309
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0773
owner: alex
created: 2026-11-13T15:20:30Z
updated: 2026-11-13T15:23:18Z
transitions:
  - to: ready
    at: 2026-11-13T15:23:18Z
    by: agent-S-0773
  - to: in-progress
    at: 2026-11-13T15:23:18Z
    by: agent-S-0773
stream: S-0773
tags: []
after: [T-3308]
---
# T-3309 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0773 --step publish --exit 2`, giving it the step's output.

## Notes
