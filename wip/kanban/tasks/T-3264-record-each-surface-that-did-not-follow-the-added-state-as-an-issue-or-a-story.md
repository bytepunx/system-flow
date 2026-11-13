---
id: T-3264
type: task
nature: feature
title: Record each surface that did not follow the added state as an issue or a story under E-0034, and say in workflow.md that the outcome holds
status: backlog
parent: S-0756
owner: alex
created: 2026-11-12T16:10:40Z
updated: 2026-11-12T16:10:40Z
transitions: []
stream: S-0756
tags: [cli]
touches: [design/system/workflow.md, docs/operators/index.md]
after: [T-3261, T-3263]
---
# T-3264 Record each surface that did not follow the added state as an issue or a story under E-0034, and say in workflow.md that the outcome holds

## Work

Waits for T-3261 and T-3263: it records what those two tests found. For each surface that needed a change to follow the added state, open an issue with `flai issue new`, or ask on a thread on the story for a story under E-0034. Say in `design/system/workflow.md` and `docs/operators/index.md` that a state, a rule, and a timer added in the workflow file alone are followed end to end, linking the test.

## Done when

Every gap the two tests found is an issue or a thread, and both documents state the outcome.

## Notes
