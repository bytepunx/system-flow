---
id: T-3082
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0718
owner: alex
created: 2026-11-11T01:44:41Z
updated: 2026-11-11T22:15:21Z
transitions:
  - to: ready
    at: 2026-11-11T22:15:21Z
    by: agent-S-0718
  - to: in-progress
    at: 2026-11-11T22:15:21Z
    by: agent-S-0718
stream: S-0718
tags: []
after: [T-3081]
---
# T-3082 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0718 --step publish --exit 2`, giving it the step's output.

## Notes
