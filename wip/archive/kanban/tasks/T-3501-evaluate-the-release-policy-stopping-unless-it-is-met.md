---
id: T-3501
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0820
owner: alex
created: 2026-11-17T02:48:43Z
updated: 2026-11-17T03:26:08Z
transitions:
  - to: ready
    at: 2026-11-17T03:25:42Z
    by: agent-S-0820
  - to: in-progress
    at: 2026-11-17T03:25:43Z
    by: agent-S-0820
  - to: done
    at: 2026-11-17T03:26:08Z
    by: agent-S-0820
stream: S-0820
tags: []
after: [T-3500]
usage:
  source: log
  seconds: 25
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 896
      cache_read: 486808
      cache_write: 19027
      cost: 0.2675
---
# T-3501 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0820 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-17: exited 0 in the main checkout. The policy threshold is met: 2 accepted stories not yet released, S-0739 (80.17 USD/week) and S-0747 (67.7 USD/week), 147.87 USD/week of 1000, count 2 of 1.
