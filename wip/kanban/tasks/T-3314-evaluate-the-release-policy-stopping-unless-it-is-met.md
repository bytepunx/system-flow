---
id: T-3314
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0775
owner: alex
created: 2026-11-13T16:13:36Z
updated: 2026-11-13T16:15:20Z
transitions:
  - to: ready
    at: 2026-11-13T16:15:09Z
    by: agent-S-0775
  - to: in-progress
    at: 2026-11-13T16:15:09Z
    by: agent-S-0775
  - to: done
    at: 2026-11-13T16:15:20Z
    by: agent-S-0775
stream: S-0775
tags: []
after: [T-3313]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 48
      cache_read: 330225
      cache_write: 1078
      cost: 0.1241
---
# T-3314 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0775 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: exited 0 in the main checkout. The threshold policy is met: 1 accepted story not yet released (S-0751), at the threshold of 1.
