---
id: T-3128
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0728
owner: alex
created: 2026-11-11T23:40:03Z
updated: 2026-11-11T23:41:43Z
transitions:
  - to: ready
    at: 2026-11-11T23:41:30Z
    by: agent-S-0728
  - to: in-progress
    at: 2026-11-11T23:41:30Z
    by: agent-S-0728
  - to: done
    at: 2026-11-11T23:41:43Z
    by: agent-S-0728
stream: S-0728
tags: []
after: [T-3127]
usage:
  source: log
  seconds: 13
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 280
      cache_read: 155838
      cache_write: 6283
      cost: 0.087
---
# T-3128 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0728 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-11: exited 0 in the main checkout. The release policy threshold is met: 1 accepted story not yet released (S-0703, 22.5 USD/week), count 1 of 1.
