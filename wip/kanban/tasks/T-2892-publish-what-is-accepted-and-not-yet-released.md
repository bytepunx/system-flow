---
id: T-2892
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0664
owner: alex
created: 2026-11-10T08:33:20Z
updated: 2026-11-10T08:35:11Z
transitions:
  - to: ready
    at: 2026-11-10T08:35:11Z
    by: agent-S-0664
  - to: in-progress
    at: 2026-11-10T08:35:11Z
    by: agent-S-0664
stream: S-0664
tags: []
after: [T-2891]
---
# T-2892 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0664 --step publish --exit 2`, giving it the step's output.

## Notes
