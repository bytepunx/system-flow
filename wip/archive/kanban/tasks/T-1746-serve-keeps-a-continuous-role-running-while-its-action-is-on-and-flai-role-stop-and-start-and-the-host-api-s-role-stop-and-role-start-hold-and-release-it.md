---
id: T-1746
type: task
nature: feature
title: serve keeps a continuous role running while its action is on, and flai role stop and start and the host API's role.stop and role.start hold and release it
status: cancelled
parent: S-0418
owner: alex
created: 2026-10-16T05:45:36Z
updated: 2026-10-29T16:36:18Z
transitions:
  - to: cancelled
    at: 2026-10-29T16:36:18Z
    by: agent-S-0418
stream: S-0418
tags: [cli]
touches: [flai/internal/serve/role_continuous.go, flai/internal/serve/role_continuous_test.go, flai/internal/serve/agents.go, flai/cmd/role.go, flai/cmd/role_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/agent.ts, docs/users/flai-reference.md]
after: [T-1741]
---
# T-1746 serve keeps a continuous role running while its action is on, and flai role stop and start and the host API's role.stop and role.start hold and release it

## Work

Generalize the orchestrator's launcher (`orchestrate.go`, `orchestrate_hold.go`) for a continuous role, in a file of its own. Leave the orchestrator's own code to S-0419.

- Start the role while its action is on, and start it again when its run ends.
- Wait at least a minute after a failed run before starting it again, and record a refused start once.
- Stop it when the action is turned off, logging `stopped: <action> turned off`.
- Hold it with `flai role stop <name>` and `role.stop`, and release it with `flai role start <name>` and `role.start`.
- Regenerate `docs/users/flai-reference.md`.

Waits for T-1741. Runs beside T-1743 and T-1744.

## Done when

`go test ./internal/serve/... ./cmd/... ./internal/hostapi/...` passes, with tests for the start, the restart, the minute's wait, the stop, and the hold.

## Notes
- 2026-10-29T16:36:18Z: moved to cancelled: TH-0749 option A: ADR-0209 §1 rejects continuous, so no role is kept running
