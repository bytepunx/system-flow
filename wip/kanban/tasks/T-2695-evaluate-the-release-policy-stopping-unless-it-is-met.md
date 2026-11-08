---
id: T-2695
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0635
owner: alex
created: 2026-11-08T01:14:12Z
updated: 2026-11-08T01:14:12Z
transitions: []
stream: S-0635
tags: []
after: [T-2694]
---
# T-2695 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout.

## Notes
