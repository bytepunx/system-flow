---
id: T-2730
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0639
owner: alex
created: 2026-11-09T11:24:54Z
updated: 2026-11-09T11:24:54Z
transitions: []
stream: S-0639
tags: []
after: [T-2729]
---
# T-2730 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout.

## Notes
