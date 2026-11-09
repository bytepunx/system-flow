---
id: T-2733
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0641
owner: alex
created: 2026-11-09T16:56:49Z
updated: 2026-11-09T16:56:49Z
transitions: []
stream: S-0641
tags: []
after: [T-2732]
---
# T-2733 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout.

## Notes
