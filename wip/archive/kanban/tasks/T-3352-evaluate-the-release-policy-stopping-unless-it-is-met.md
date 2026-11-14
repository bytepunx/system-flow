---
id: T-3352
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0784
owner: alex
created: 2026-11-14T07:05:35Z
updated: 2026-11-14T07:07:40Z
transitions:
  - to: ready
    at: 2026-11-14T07:07:23Z
    by: agent-S-0784
  - to: in-progress
    at: 2026-11-14T07:07:23Z
    by: agent-S-0784
  - to: done
    at: 2026-11-14T07:07:40Z
    by: agent-S-0784
stream: S-0784
tags: []
after: [T-3351]
usage:
  source: log
  seconds: 17
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 493
      cache_read: 325418
      cache_write: 11162
      cost: 0.1643
---
# T-3352 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0784 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-14: exited 0 in the main checkout. The release policy threshold is met: 1 accepted story not yet released (S-0782, 15 USD/week), at the count threshold of 1.
