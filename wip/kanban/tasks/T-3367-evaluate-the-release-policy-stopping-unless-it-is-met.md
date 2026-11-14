---
id: T-3367
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0789
owner: alex
created: 2026-11-14T23:49:29Z
updated: 2026-11-14T23:49:29Z
transitions: []
stream: S-0789
tags: []
after: [T-3366]
---
# T-3367 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0789 --step evaluate --exit 2`, giving it the step's output.

## Notes
