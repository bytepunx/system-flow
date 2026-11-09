---
id: T-2730
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0639
owner: alex
created: 2026-11-09T11:24:54Z
updated: 2026-11-09T11:25:59Z
transitions:
  - to: ready
    at: 2026-11-09T11:25:50Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T11:25:50Z
    by: system-flow
  - to: done
    at: 2026-11-09T11:25:59Z
    by: agent-S-0639
stream: S-0639
tags: []
after: [T-2729]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 6
      output: 90
      cache_read: 457002
      cache_write: 2025
      cost: 0.1674
---
# T-2730 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout.

## Notes

- 2026-11-09: `flai release --evaluate --require-met` (as `scripts/flai.sh`) ran in the main checkout and exited 0: threshold met, 1 accepted story not yet released (S-0618, 25 USD/week), count 1 of 1.
