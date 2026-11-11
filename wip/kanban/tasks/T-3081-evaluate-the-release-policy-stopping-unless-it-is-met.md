---
id: T-3081
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0718
owner: alex
created: 2026-11-11T01:44:41Z
updated: 2026-11-11T01:44:41Z
transitions: []
stream: S-0718
tags: []
after: [T-3080]
---
# T-3081 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0718 --step evaluate --exit 2`, giving it the step's output.

## Notes
