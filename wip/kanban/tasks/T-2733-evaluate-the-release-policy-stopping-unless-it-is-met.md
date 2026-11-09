---
id: T-2733
type: task
nature: improvement
title: Evaluate the release policy, stopping unless it is met
status: done
parent: S-0641
owner: alex
created: 2026-11-09T16:56:49Z
updated: 2026-11-09T16:57:55Z
transitions:
  - to: ready
    at: 2026-11-09T16:57:45Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T16:57:45Z
    by: system-flow
  - to: done
    at: 2026-11-09T16:57:55Z
    by: agent-S-0641
stream: S-0641
tags: []
after: [T-2732]
---
# T-2733 Evaluate the release policy, stopping unless it is met

## Work

Run `flai release --evaluate --require-met` from the project's main checkout, `/home/alex/git/bytepunx/system-flow`, not from this story's worktree.

## Done when

`flai release --evaluate --require-met` exits 0 in the main checkout.

## Notes

- 2026-11-09: `flai release --evaluate --require-met` ran in the main checkout and exited 0: the policy is met, 2 accepted stories not yet released (S-0610, S-0615), 21680 USD/week against a threshold of 1000 USD/week.
