---
id: T-1625
type: task
nature: feature
title: flai accept refuses on the forge blockers before merging, for the operator, the dashboard's acceptance, and the orchestrator's alike
status: backlog
parent: S-0395
owner: alex
created: 2026-10-16T04:38:15Z
updated: 2026-10-16T04:38:15Z
transitions: []
stream: S-0395
tags: [cli]
touches: [flai/cmd/accept.go, flai/cmd/accept_forge_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go]
after: [T-1623]
---
# T-1625 flai accept refuses on the forge blockers before merging, for the operator, the dashboard's acceptance, and the orchestrator's alike

## Work

- In `flai/cmd/accept.go`, make sure the real run refuses on T-1623's blockers before anything is merged, with exit and message as for any blocker, and that `--json` carries them.
- In `flai/internal/hostapi/writes.go`, make sure `accept.run` and `orchestratorAccepts` refuse on them too, naming each.
- Test, in `flai/cmd/accept_forge_test.go`, the operator's, the dashboard's, and the orchestrator's acceptance of a story:
  - with a failing check
  - with a pending check
  - with changes requested
  - with too few approvals
  - with a stale state
  - one that passes, which then merges locally as today until S-0397
- Waits for T-1623: it refuses on that task's blockers.

## Done when

Each of the three acceptances refuses on each forge blocker with nothing merged, accepts a story that passes, and the tests pass.

## Notes
