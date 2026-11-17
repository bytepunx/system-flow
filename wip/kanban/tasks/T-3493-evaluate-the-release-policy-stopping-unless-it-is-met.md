---
id: T-3493
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0818
owner: alex
created: 2026-11-17T02:27:29Z
updated: 2026-11-17T02:30:21Z
transitions:
  - to: ready
    at: 2026-11-17T02:29:31Z
    by: agent-S-0818
  - to: in-progress
    at: 2026-11-17T02:29:31Z
    by: agent-S-0818
  - to: done
    at: 2026-11-17T02:30:21Z
    by: agent-S-0818
stream: S-0818
tags: []
after: [T-3492]
usage:
  source: log
  seconds: 50
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 8
      output: 99
      cache_read: 673639
      cache_write: 1982
      cost: 0.2569
---
# T-3493 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0818 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-17: ran in the main checkout and exited 0. Policy `threshold` met: 1 accepted story not yet released (count 1 of 1; value 450 of 1000 USD/week). Pending: S-0797.
