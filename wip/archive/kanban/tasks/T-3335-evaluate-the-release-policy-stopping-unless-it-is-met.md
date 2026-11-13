---
id: T-3335
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0780
owner: alex
created: 2026-11-13T20:51:49Z
updated: 2026-11-13T20:53:53Z
transitions:
  - to: ready
    at: 2026-11-13T20:53:21Z
    by: agent-S-0780
  - to: in-progress
    at: 2026-11-13T20:53:21Z
    by: agent-S-0780
  - to: done
    at: 2026-11-13T20:53:53Z
    by: agent-S-0780
stream: S-0780
tags: []
after: [T-3334]
usage:
  source: log
  seconds: 32
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 1297
      cache_read: 656397
      cache_write: 24273
      cost: 0.3514
---
# T-3335 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0780 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: ran in the main checkout and exited 0. The policy, `threshold`, is met: 1 accepted story not yet released, at the threshold of 1; value 75 USD/week of 1000 USD/week. Pending: S-0777.
