---
id: T-3514
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0825
owner: alex
created: 2026-11-17T07:54:16Z
updated: 2026-11-17T07:54:16Z
transitions: []
stream: S-0825
tags: []
after: [T-3513]
---
# T-3514 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0825 --step evaluate --exit 2`, giving it the step's output.

## Notes
