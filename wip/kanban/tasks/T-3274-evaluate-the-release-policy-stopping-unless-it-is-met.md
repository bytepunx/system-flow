---
id: T-3274
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0764
owner: alex
created: 2026-11-13T08:55:19Z
updated: 2026-11-13T08:57:16Z
transitions:
  - to: ready
    at: 2026-11-13T08:57:06Z
    by: agent-S-0764
  - to: in-progress
    at: 2026-11-13T08:57:06Z
    by: agent-S-0764
  - to: done
    at: 2026-11-13T08:57:16Z
    by: agent-S-0764
stream: S-0764
tags: []
after: [T-3273]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 325
      cache_read: 161665
      cache_write: 6585
      cost: 0.0915
---
# T-3274 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0764 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: ran in the main checkout and exited 0: the release policy threshold is met, with 1 accepted story not yet released (S-0660, 1250 USD/week, against 1000 USD/week and a count of 1).
