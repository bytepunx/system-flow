---
id: T-2696
type: task
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
parent: S-0635
owner: alex
created: 2026-11-08T01:14:12Z
updated: 2026-11-09T11:24:02Z
transitions:
  - to: ready
    at: 2026-11-09T11:24:01Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T11:24:02Z
    by: system-flow
stream: S-0635
tags: []
after: [T-2695]
---
# T-2696 Publish what is accepted and not yet released

## Work

Run `flai release --pending --json` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --pending --json` exits 0 in the main checkout.

## Notes
