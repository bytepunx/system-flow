---
id: T-1636
type: task
nature: feature
title: The dashboard's acceptance under pull-request acceptance merges through the remote in flai serve, and refuses while the forge host action is off
status: backlog
parent: S-0397
owner: alex
created: 2026-10-16T04:40:11Z
updated: 2026-10-16T04:40:18Z
transitions: []
stream: S-0397
tags: [cli, dashboard]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go]
after: [T-1634, T-1635]
---
# T-1636 The dashboard's acceptance under pull-request acceptance merges through the remote in flai serve, and refuses while the forge host action is off

## Work

- In `flai/internal/hostapi/writes.go`, have `accept.run` for a project under pull-request acceptance merge through the remote in flai serve's own process, where the token's variable is at hand, and then finish the acceptance by T-1635's path. If `accept.run` already runs `flai accept` in flai serve's environment, this is T-1635's direct merge, and the task is the refusal and the tests.
- Refuse while the `forge` host action is off, naming it.
- Keep `orchestratorAccepts` as S-0395 left it: the gate first.
- Test a dashboard acceptance that merges, one already merged, and one refused for the host action.
- Waits for T-1634 (the merge) and T-1635 (the acceptance path it finishes with).

## Done when

A dashboard acceptance under pull-request acceptance merges through the fake forge, refuses with the host action off, and the hostapi tests pass.

## Notes
