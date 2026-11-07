---
id: T-2613
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0597
owner: alex
created: 2026-11-07T22:31:23Z
updated: 2026-11-07T22:31:23Z
transitions: []
stream: S-0597
tags: []
after: [T-2612]
---
# T-2613 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout.

## Notes
