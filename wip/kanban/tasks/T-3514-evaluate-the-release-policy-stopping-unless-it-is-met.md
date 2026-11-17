---
id: T-3514
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0825
owner: alex
created: 2026-11-17T07:54:16Z
updated: 2026-11-17T08:15:15Z
transitions:
  - to: ready
    at: 2026-11-17T08:15:04Z
    by: agent-S-0825
  - to: in-progress
    at: 2026-11-17T08:15:05Z
    by: agent-S-0825
  - to: done
    at: 2026-11-17T08:15:15Z
    by: agent-S-0825
stream: S-0825
tags: []
after: [T-3513]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 653
      cache_read: 326264
      cache_write: 15210
      cost: 0.2
---
# T-3514 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0825 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-17: exited 0 in the main checkout. Policy `threshold` met: 4 accepted stories not yet released (count 4 of 1; value 842.34 of 1000 USD/week). Pending: S-0742, S-0753, S-0817, S-0819.
