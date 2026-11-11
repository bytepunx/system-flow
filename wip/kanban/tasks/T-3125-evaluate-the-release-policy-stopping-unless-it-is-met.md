---
id: T-3125
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0727
owner: alex
created: 2026-11-11T23:22:51Z
updated: 2026-11-11T23:24:31Z
transitions:
  - to: ready
    at: 2026-11-11T23:24:18Z
    by: agent-S-0727
  - to: in-progress
    at: 2026-11-11T23:24:18Z
    by: agent-S-0727
  - to: done
    at: 2026-11-11T23:24:31Z
    by: agent-S-0727
stream: S-0727
tags: []
after: [T-3124]
usage:
  source: log
  seconds: 13
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 40
      cache_read: 320434
      cache_write: 1449
      cost: 0.1195
---
# T-3125 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0727 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-11: exited 0 in the main checkout. The release policy threshold is met: 1 accepted story not yet released (S-0720, 375 USD/week), count 1 of 1.
