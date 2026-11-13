---
id: T-3282
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0765
owner: alex
created: 2026-11-13T09:17:34Z
updated: 2026-11-13T09:18:54Z
transitions:
  - to: ready
    at: 2026-11-13T09:18:46Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T09:18:46Z
    by: system-flow
  - to: done
    at: 2026-11-13T09:18:54Z
    by: agent-S-0765
stream: S-0765
tags: []
after: [T-3281]
usage:
  source: log
  seconds: 8
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 644
      cache_read: 309993
      cache_write: 16768
      cost: 0.209
---
# T-3282 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0765 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: ran in the main checkout and exited 0: the release policy threshold is met, with 1 accepted story not yet released (S-0735, no value, count 1 of 1).
