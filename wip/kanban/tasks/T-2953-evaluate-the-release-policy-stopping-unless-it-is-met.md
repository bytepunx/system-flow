---
id: T-2953
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0704
owner: alex
created: 2026-11-10T08:56:47Z
updated: 2026-11-11T01:42:36Z
transitions:
  - to: ready
    at: 2026-11-11T01:42:27Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T01:42:27Z
    by: system-flow
  - to: done
    at: 2026-11-11T01:42:36Z
    by: agent-S-0704
stream: S-0704
tags: []
after: [T-2952]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 66
      cache_read: 330881
      cache_write: 1200
      cost: 0.1224
---
# T-2953 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0704 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-11: ran as `scripts/flai.sh release --evaluate --require-met --judged ""` in the main checkout; exit 0. Policy `threshold` met: 5 accepted stories not yet released (S-0665, S-0666, S-0675, S-0677, S-0679), worth 1087.58 USD/week against a threshold of 1000 USD/week; S-0666 and S-0679 carry no value.
