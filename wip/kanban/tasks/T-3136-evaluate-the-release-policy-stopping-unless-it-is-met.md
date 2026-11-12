---
id: T-3136
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0731
owner: alex
created: 2026-11-12T03:44:58Z
updated: 2026-11-12T03:47:02Z
transitions:
  - to: ready
    at: 2026-11-12T03:46:52Z
    by: agent-S-0731
  - to: in-progress
    at: 2026-11-12T03:46:52Z
    by: agent-S-0731
  - to: done
    at: 2026-11-12T03:47:02Z
    by: agent-S-0731
stream: S-0731
tags: []
after: [T-3135]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 81
      cache_read: 326134
      cache_write: 1068
      cost: 0.1216
---
# T-3136 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0731 --step evaluate --exit 2`, giving it the step's output.

## Notes

`flai release --evaluate --require-met --judged ""` exited 0 in the main checkout (2026-11-12), with flai 1.86.17: release policy threshold met, 3 accepted stories not yet released (S-0681, S-0701, S-0717), count 3 of 1, value 77.5 USD/week.
