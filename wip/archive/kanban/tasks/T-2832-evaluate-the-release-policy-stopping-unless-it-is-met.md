---
id: T-2832
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0656
owner: alex
created: 2026-11-10T06:45:10Z
updated: 2026-11-10T07:01:02Z
transitions:
  - to: ready
    at: 2026-11-10T07:00:54Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T07:00:54Z
    by: system-flow
  - to: done
    at: 2026-11-10T07:01:02Z
    by: agent-S-0656
stream: S-0656
tags: []
after: [T-2831]
usage:
  source: log
  seconds: 8
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 539
      cache_read: 308111
      cache_write: 11737
      cost: 0.1663
---
# T-2832 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0656 --step evaluate --exit 2`, giving it the step's output.

## Notes

- `flai release --evaluate --require-met --judged ""` exited 0 in the main checkout: release policy threshold met, 3 accepted stories not yet released (S-0625, S-0630, S-0631), 100 USD/week, count 3 of 1.
