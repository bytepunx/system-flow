---
id: T-3518
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0826
owner: alex
created: 2026-11-17T08:16:46Z
updated: 2026-11-17T08:16:46Z
transitions: []
stream: S-0826
tags: []
after: [T-3517]
---
# T-3518 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0826 --step evaluate --exit 2`, giving it the step's output.

## Notes
