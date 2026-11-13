---
id: T-3301
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: backlog
parent: S-0771
owner: alex
created: 2026-11-13T13:41:40Z
updated: 2026-11-13T13:41:40Z
transitions: []
stream: S-0771
tags: []
after: [T-3300]
---
# T-3301 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0771 --step evaluate --exit 2`, giving it the step's output.

## Notes
