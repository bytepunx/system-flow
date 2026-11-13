---
id: T-3298
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0770
owner: alex
created: 2026-11-13T13:16:45Z
updated: 2026-11-13T13:16:45Z
transitions: []
stream: S-0770
tags: []
after: [T-3297]
---
# T-3298 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0770 --step evaluate --exit 2`, giving it the step's output.

## Notes
