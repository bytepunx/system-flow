---
id: T-3311
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0774
owner: alex
created: 2026-11-13T15:43:49Z
updated: 2026-11-13T15:45:09Z
transitions:
  - to: ready
    at: 2026-11-13T15:44:58Z
    by: agent-S-0774
  - to: in-progress
    at: 2026-11-13T15:44:59Z
    by: agent-S-0774
  - to: done
    at: 2026-11-13T15:45:09Z
    by: agent-S-0774
stream: S-0774
tags: []
after: [T-3310]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 6
      output: 64
      cache_read: 490432
      cache_write: 1759
      cost: 0.1843
---
# T-3311 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met --judged ""` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met --judged ""` exits 0 in the main checkout, or exits 2, which stops the run: record it with `flai formula stop S-0774 --step evaluate --exit 2`, giving it the step's output.

## Notes

- 2026-11-13: exited 0 in the main checkout (flai 1.86.17): threshold met, 1 accepted story not yet released (S-0663, 5 USD/week), count 1 of 1.
