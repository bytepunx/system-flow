---
id: T-3506
type: task
nature: remediation
title: flai guard prunes the session's record from its transcript before it holds a story agent's wait_for_events, with a test that reproduces I-0376
status: backlog
parent: S-0819
owner: alex
created: 2026-11-17T03:34:53Z
updated: 2026-11-17T03:34:53Z
transitions: []
stream: S-0819
tags: [guard, go]
touches: [flai/cmd/guard.go, flai/cmd/guard_test.go]
after: [T-3505]
---
# T-3506 flai guard prunes the session's record from its transcript before it holds a story agent's wait_for_events, with a test that reproduces I-0376

## Work

In `flai/cmd/guard.go`, have `holdWaits` prune the session's record with T-3505's functions, from the transcript the hook input names, before it reads the running sub-agents; it waits for T-3505 because it calls what that task builds. A transcript it cannot read is warned of and leaves the record as it is; a record or project it cannot read still lets the call through, as ADR-0092 §4 says.

In `flai/cmd/guard_test.go`, reproduce I-0376 end to end beside `TestGuardRefusesTheStorysAgentsWaitOnItsSubAgents`: a `SubagentStart` with no `SubagentStop`, a session transcript holding that sub-agent's foreground result marked an error, then the story agent's `wait_for_events`, which passes and leaves no record. Keep the case where the transcript lacks the result, or the launch was in the background, refused as before.

## Done when

- The new test fails on the code before this task and passes after it.
- `flai test` on the two paths passes, and the existing sub-agent wait tests pass unchanged.

## Notes
