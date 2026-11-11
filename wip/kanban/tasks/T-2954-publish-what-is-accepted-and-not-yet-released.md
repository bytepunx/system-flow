---
id: T-2954
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0704
owner: alex
created: 2026-11-10T08:56:47Z
updated: 2026-11-11T01:42:53Z
transitions:
  - to: ready
    at: 2026-11-11T01:42:53Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T01:42:53Z
    by: system-flow
stream: S-0704
tags: []
after: [T-2953]
---
# T-2954 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0704 --step publish --exit 2`, giving it the step's output.

## Notes
