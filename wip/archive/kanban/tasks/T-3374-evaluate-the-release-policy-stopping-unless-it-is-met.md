---
id: T-3374
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0791
owner: alex
created: 2026-11-15T07:18:49Z
updated: 2026-11-15T21:00:53Z
transitions:
  - to: ready
    at: 2026-11-15T21:00:41Z
    by: agent-S-0791
  - to: in-progress
    at: 2026-11-15T21:00:41Z
    by: agent-S-0791
  - to: done
    at: 2026-11-15T21:00:53Z
    by: agent-S-0791
stream: S-0791
tags: []
after: [T-3373]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1027
      cache_read: 479883
      cache_write: 19871
      cost: 0.2755
---
# T-3374 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0791 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-15: exit 0. The threshold policy is met: 3 accepted stories are not yet released (S-0743, S-0744, S-0749), worth 261.88 USD/week, and 3 is at or over the threshold of 1.
