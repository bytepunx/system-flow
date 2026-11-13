---
id: T-3289
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0768
owner: alex
created: 2026-11-13T11:43:32Z
updated: 2026-11-13T11:43:32Z
transitions: []
stream: S-0768
tags: []
after: [T-3288]
---
# T-3289 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0768 --step evaluate --exit 2`, giving it the step's output.

## Notes
