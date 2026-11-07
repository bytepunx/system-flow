---
id: T-2590
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0594
owner: alex
created: 2026-11-07T21:47:12Z
updated: 2026-11-07T21:47:12Z
transitions: []
stream: S-0594
tags: []
after: [T-2589]
---
# T-2590 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout.

## Notes
