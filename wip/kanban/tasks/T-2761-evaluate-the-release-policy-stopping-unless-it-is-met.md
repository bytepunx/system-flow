---
id: T-2761
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0642
owner: alex
created: 2026-11-09T17:28:50Z
updated: 2026-11-09T17:28:50Z
transitions: []
stream: S-0642
tags: []
after: [T-2760]
---
# T-2761 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout.

## Notes
