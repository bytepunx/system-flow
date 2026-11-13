---
id: T-3290
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0768
owner: alex
created: 2026-11-13T11:43:32Z
updated: 2026-11-13T11:45:13Z
transitions:
  - to: ready
    at: 2026-11-13T11:45:13Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T11:45:13Z
    by: system-flow
stream: S-0768
tags: []
after: [T-3289]
---
# T-3290 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0768 --step publish --exit 2`, giving it the step's output.

## Notes
