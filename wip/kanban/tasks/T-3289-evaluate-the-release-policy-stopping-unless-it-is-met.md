---
id: T-3289
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0768
owner: alex
created: 2026-11-13T11:43:32Z
updated: 2026-11-13T11:45:01Z
transitions:
  - to: ready
    at: 2026-11-13T11:44:51Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T11:44:51Z
    by: system-flow
  - to: done
    at: 2026-11-13T11:45:01Z
    by: agent-S-0768
stream: S-0768
tags: []
after: [T-3288]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 32
      cache_read: 336708
      cache_write: 1966
      cost: 0.1264
---
# T-3289 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0768 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: exited 0 in the main checkout. The threshold policy is met: 1 accepted story not yet released, at the threshold of 1, worth 2.5 USD/week. Pending: S-0761.
