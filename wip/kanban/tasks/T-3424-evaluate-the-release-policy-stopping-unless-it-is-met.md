---
id: T-3424
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0802
owner: alex
created: 2026-11-16T04:12:37Z
updated: 2026-11-16T04:14:01Z
transitions:
  - to: ready
    at: 2026-11-16T04:13:50Z
    by: agent-S-0802
  - to: in-progress
    at: 2026-11-16T04:13:50Z
    by: agent-S-0802
  - to: done
    at: 2026-11-16T04:14:01Z
    by: agent-S-0802
stream: S-0802
tags: []
after: [T-3423]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 30
      cache_read: 335647
      cache_write: 966
      cost: 0.1277
---
# T-3424 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0802 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-16: exited 0 in the main checkout. Release policy `threshold` met: 1 accepted story not yet released (S-0799, 50 USD/week), count 1 of 1.
