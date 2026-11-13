---
id: T-3267
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0762
owner: alex
created: 2026-11-13T07:24:38Z
updated: 2026-11-13T07:26:31Z
transitions:
  - to: ready
    at: 2026-11-13T07:26:31Z
    by: agent-S-0762
  - to: in-progress
    at: 2026-11-13T07:26:31Z
    by: agent-S-0762
stream: S-0762
tags: []
after: [T-3266]
---
# T-3267 Publish what is accepted and not yet released

## Work

Run `flai release --pending --require-pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --require-pending --json` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0762 --step publish --exit 2`, giving it the step's output.

## Notes
