---
id: T-3264
type: task
nature: feature
title: Record each surface that did not follow the added state as an issue or a story under E-0034, and say in workflow.md that the outcome holds
status: done
parent: S-0756
owner: alex
created: 2026-11-12T16:10:40Z
updated: 2026-11-17T11:36:05Z
transitions:
  - to: ready
    at: 2026-11-17T11:35:14Z
    by: agent-S-0756
  - to: in-progress
    at: 2026-11-17T11:35:14Z
    by: agent-S-0756
  - to: done
    at: 2026-11-17T11:36:05Z
    by: agent-S-0756
stream: S-0756
tags: [cli]
touches: [design/system/workflow.md, docs/operators/runbooks/change-the-board.md]
after: [T-3261, T-3263]
usage:
  source: log
  seconds: 51
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 18
      output: 4853
      cache_read: 1032398
      cache_write: 37248
      cost: 0.5232
---
# T-3264 Record each surface that did not follow the added state as an issue or a story under E-0034, and say in workflow.md that the outcome holds

## Work

Waits for T-3261 and T-3263: it records what those two tests found. For each surface that needed a change to follow the added state, open an issue with `flai issue new`, or ask on a thread on the story for a story under E-0034. Say in `design/system/workflow.md` and `docs/operators/index.md` that a state, a rule, and a timer added in the workflow file alone are followed end to end, linking the test.

## Done when

Every gap the two tests found is an issue or a thread, and both documents state the outcome.

## Notes
