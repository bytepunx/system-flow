---
id: T-3266
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0762
owner: alex
created: 2026-11-13T07:24:38Z
updated: 2026-11-13T07:26:22Z
transitions:
  - to: ready
    at: 2026-11-13T07:26:17Z
    by: agent-S-0762
  - to: in-progress
    at: 2026-11-13T07:26:17Z
    by: agent-S-0762
  - to: done
    at: 2026-11-13T07:26:22Z
    by: agent-S-0762
stream: S-0762
tags: []
after: [T-3265]
usage:
  source: log
  seconds: 5
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 359
      cache_read: 159081
      cache_write: 7586
      cost: 0.0997
---
# T-3266 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0762 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: `scripts/flai.sh release --evaluate --require-met --judged ""` ran in the main checkout and exited 0: threshold met, count 4 of 1, value 80 of 1000 USD/week; pending S-0683, S-0723, S-0726, S-0732.
