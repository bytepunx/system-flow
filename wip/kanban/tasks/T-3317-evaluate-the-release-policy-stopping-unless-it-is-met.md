---
id: T-3317
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0776
owner: alex
created: 2026-11-13T17:20:42Z
updated: 2026-11-13T17:20:42Z
transitions: []
stream: S-0776
tags: []
after: [T-3316]
---
# T-3317 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0776 --step evaluate --exit 2`, giving it the step's output.

## Notes
