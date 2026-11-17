---
id: T-3489
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0816
owner: alex
created: 2026-11-17T01:33:37Z
updated: 2026-11-17T01:36:33Z
transitions:
  - to: ready
    at: 2026-11-17T01:36:24Z
    by: agent-S-0816
  - to: in-progress
    at: 2026-11-17T01:36:25Z
    by: agent-S-0816
  - to: done
    at: 2026-11-17T01:36:33Z
    by: agent-S-0816
stream: S-0816
tags: []
after: [T-3488]
usage:
  source: log
  seconds: 8
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 620
      cache_read: 324389
      cache_write: 12184
      cost: 0.1748
---
# T-3489 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0816 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-17: exited 0 in the main checkout. The policy is threshold, and it is met: 1 accepted story is not yet released, S-0804, at or over the threshold of 1. Its value is 0 USD/week of 1000 USD/week, and S-0804 has no value.
