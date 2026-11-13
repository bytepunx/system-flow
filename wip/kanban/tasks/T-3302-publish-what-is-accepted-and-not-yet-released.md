---
id: T-3302
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0771
owner: alex
created: 2026-11-13T13:41:40Z
updated: 2026-11-13T14:37:10Z
transitions:
  - to: ready
    at: 2026-11-13T14:37:10Z
    by: agent-S-0771
  - to: in-progress
    at: 2026-11-13T14:37:10Z
    by: agent-S-0771
stream: S-0771
tags: []
after: [T-3301]
---
# T-3302 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0771 --step publish --exit 2`, giving it the step's output.

## Notes
