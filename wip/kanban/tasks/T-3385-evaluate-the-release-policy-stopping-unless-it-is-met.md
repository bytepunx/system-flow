---
id: T-3385
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0793
owner: alex
created: 2026-11-15T21:17:19Z
updated: 2026-11-15T21:17:19Z
transitions: []
stream: S-0793
tags: []
after: [T-3384]
---
# T-3385 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0793 --step evaluate --exit 2`, giving it the step's output.

## Notes
