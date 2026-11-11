---
id: T-3081
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0718
owner: alex
created: 2026-11-11T01:44:41Z
updated: 2026-11-11T22:15:10Z
transitions:
  - to: ready
    at: 2026-11-11T22:15:00Z
    by: agent-S-0718
  - to: in-progress
    at: 2026-11-11T22:15:00Z
    by: agent-S-0718
  - to: done
    at: 2026-11-11T22:15:10Z
    by: agent-S-0718
stream: S-0718
tags: []
after: [T-3080]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 48
      cache_read: 331992
      cache_write: 1787
      cost: 0.1239
---
# T-3081 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0718 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-11: ran in the main checkout with flai 1.86.17 and exited 0: `release policy threshold: met`, 4 accepted stories not yet released (S-0689, S-0707, S-0714, S-0715), count 4 of 1, value 97.5 USD/week of 1000 USD/week.
