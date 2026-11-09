---
id: T-2761
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0642
owner: alex
created: 2026-11-09T17:28:50Z
updated: 2026-11-09T21:54:00Z
transitions:
  - to: ready
    at: 2026-11-09T21:53:52Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T21:53:53Z
    by: system-flow
  - to: done
    at: 2026-11-09T21:54:00Z
    by: agent-S-0642
stream: S-0642
tags: []
after: [T-2760]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 2
      output: 8
      cache_read: 154054
      cache_write: 1223
      cost: 0.0567
---
# T-2761 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout.

## Notes

- 2026-11-09: `flai release --evaluate --require-met` ran in the main checkout and exited 0: threshold met, 8 accepted stories not yet released (S-0611, S-0619, S-0620, S-0621, S-0624, S-0626, S-0627, S-0629), worth 310 USD/week.
