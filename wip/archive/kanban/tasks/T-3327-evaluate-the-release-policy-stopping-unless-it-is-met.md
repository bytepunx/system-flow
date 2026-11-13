---
id: T-3327
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0778
owner: alex
created: 2026-11-13T20:23:36Z
updated: 2026-11-13T20:25:43Z
transitions:
  - to: ready
    at: 2026-11-13T20:25:30Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T20:25:30Z
    by: system-flow
  - to: done
    at: 2026-11-13T20:25:43Z
    by: agent-S-0778
stream: S-0778
tags: []
after: [T-3326]
usage:
  source: log
  seconds: 13
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 781
      cache_read: 477152
      cache_write: 16688
      cost: 0.2446
---
# T-3327 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0778 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: `scripts/flai.sh release --evaluate --require-met --judged ""` in the main checkout exited 0: threshold met, 1 accepted story not yet released (S-0760, no value), count 1 of 1.
