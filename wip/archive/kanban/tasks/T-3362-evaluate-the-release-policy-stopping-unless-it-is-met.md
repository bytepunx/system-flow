---
id: T-3362
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0787
owner: alex
created: 2026-11-14T18:00:18Z
updated: 2026-11-14T18:01:55Z
transitions:
  - to: ready
    at: 2026-11-14T18:01:43Z
    by: system-flow
  - to: in-progress
    at: 2026-11-14T18:01:43Z
    by: system-flow
  - to: done
    at: 2026-11-14T18:01:55Z
    by: agent-S-0787
stream: S-0787
tags: []
after: [T-3361]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 914
      cache_read: 482398
      cache_write: 23949
      cost: 0.3064
---
# T-3362 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0787 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-14: `scripts/flai.sh release --evaluate --require-met --judged ""` in the main checkout exited 0: threshold met, 1 accepted story not yet released (S-0781, 2.5 USD/week), count 1 of 1.
