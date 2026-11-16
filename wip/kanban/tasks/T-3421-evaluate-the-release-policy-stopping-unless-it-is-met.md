---
id: T-3421
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0801
owner: alex
created: 2026-11-16T01:54:24Z
updated: 2026-11-16T03:52:55Z
transitions:
  - to: ready
    at: 2026-11-16T03:52:45Z
    by: agent-S-0801
  - to: in-progress
    at: 2026-11-16T03:52:46Z
    by: agent-S-0801
  - to: done
    at: 2026-11-16T03:52:55Z
    by: agent-S-0801
stream: S-0801
tags: []
after: [T-3420]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 48
      cache_read: 334557
      cache_write: 1429
      cost: 0.1274
---
# T-3421 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0801 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-16: the evaluate ran in the main checkout and exited 0. The policy is met: 4 accepted stories are not yet released (S-0795, S-0796, S-0798, S-0800), at or over the threshold of 1, worth 737.5 USD/week.
