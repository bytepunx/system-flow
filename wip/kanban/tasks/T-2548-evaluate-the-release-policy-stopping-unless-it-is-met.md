---
id: T-2548
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0585
owner: alex
created: 2026-11-07T08:03:26Z
updated: 2026-11-07T08:03:26Z
transitions: []
stream: S-0585
tags: []
after: [T-2546]
---
# T-2548 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout.

## Notes
