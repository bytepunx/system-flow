---
id: T-3294
type: task
nature: improvement
title: Fetch the remote's branch and tags
status: in-progress
parent: S-0769
owner: alex
created: 2026-11-13T12:40:07Z
updated: 2026-11-13T12:53:36Z
transitions:
  - to: ready
    at: 2026-11-13T12:53:35Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T12:53:36Z
    by: system-flow
stream: S-0769
tags: []
---
# T-3294 Fetch the remote's branch and tags

## Work

Run `scripts/fetch.sh` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/fetch.sh` exits 0 in the main checkout.

## Notes
