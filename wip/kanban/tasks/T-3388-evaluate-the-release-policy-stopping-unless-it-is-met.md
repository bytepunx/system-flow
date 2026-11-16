---
id: T-3388
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0794
owner: alex
created: 2026-11-15T21:47:08Z
updated: 2026-11-16T00:38:35Z
transitions:
  - to: ready
    at: 2026-11-16T00:38:18Z
    by: agent-S-0794
  - to: in-progress
    at: 2026-11-16T00:38:18Z
    by: agent-S-0794
  - to: done
    at: 2026-11-16T00:38:35Z
    by: agent-S-0794
stream: S-0794
tags: []
after: [T-3387]
usage:
  source: log
  seconds: 17
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 6
      output: 77
      cache_read: 498186
      cache_write: 2300
      cost: 0.1895
---
# T-3388 Evaluate the release policy, stopping unless it is met

## Work

Run `scripts/flai.sh release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`scripts/flai.sh release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `scripts/flai.sh formula stop S-0794 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-16: the evaluation exited 0 in the main checkout. The threshold policy is met: 2 accepted stories are unreleased against a threshold of 1, worth 88.73 USD/week. The pending stories are S-0755 and S-0790.
