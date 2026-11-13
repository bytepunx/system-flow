---
id: T-3295
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0769
owner: alex
created: 2026-11-13T12:40:07Z
updated: 2026-11-13T12:54:25Z
transitions:
  - to: ready
    at: 2026-11-13T12:54:19Z
    by: agent-S-0769
  - to: in-progress
    at: 2026-11-13T12:54:19Z
    by: agent-S-0769
  - to: done
    at: 2026-11-13T12:54:25Z
    by: agent-S-0769
stream: S-0769
tags: []
after: [T-3294]
usage:
  source: log
  seconds: 6
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 2
      output: 16
      cache_read: 175193
      cache_write: 1151
      cost: 0.0658
---
# T-3295 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0769 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: `flai release --evaluate --require-met --judged ""` ran in the main checkout and exited 0: the threshold policy is met, 2 accepted stories not yet released (S-0763, S-0767), their count 2 at or over the threshold of 1, worth 580.56 USD/week.
