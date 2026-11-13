---
id: T-3317
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0776
owner: alex
created: 2026-11-13T17:20:42Z
updated: 2026-11-13T17:22:14Z
transitions:
  - to: ready
    at: 2026-11-13T17:22:05Z
    by: agent-S-0776
  - to: in-progress
    at: 2026-11-13T17:22:05Z
    by: agent-S-0776
  - to: done
    at: 2026-11-13T17:22:14Z
    by: agent-S-0776
stream: S-0776
tags: []
after: [T-3316]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 48
      cache_read: 325508
      cache_write: 1221
      cost: 0.1224
---
# T-3317 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0776 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: exited 0 in the main checkout (flai 1.86.17): threshold met, 1 accepted story not yet released (S-0757, no value), count 1 of 1.
