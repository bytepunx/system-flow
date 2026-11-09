---
id: T-2823
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0646
owner: alex
created: 2026-11-09T22:22:55Z
updated: 2026-11-09T22:22:55Z
transitions: []
stream: S-0646
tags: []
after: [T-2822]
---
# T-2823 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0646 --step evaluate --exit 2`, giving it the step's output.

## Notes
