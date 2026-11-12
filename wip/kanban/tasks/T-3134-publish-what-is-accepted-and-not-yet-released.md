---
id: T-3134
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0729
owner: alex
created: 2026-11-12T03:39:40Z
updated: 2026-11-12T03:41:47Z
transitions:
  - to: ready
    at: 2026-11-12T03:41:47Z
    by: agent-S-0729
  - to: in-progress
    at: 2026-11-12T03:41:47Z
    by: agent-S-0729
stream: S-0729
tags: []
after: [T-3133]
---
# T-3134 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0729 --step publish --exit 2`, giving it the step's output.

## Notes
