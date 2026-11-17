---
id: T-3501
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0820
owner: alex
created: 2026-11-17T02:48:43Z
updated: 2026-11-17T02:48:43Z
transitions: []
stream: S-0820
tags: []
after: [T-3500]
---
# T-3501 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0820 --step evaluate --exit 2`, giving it the step's output.

## Notes
