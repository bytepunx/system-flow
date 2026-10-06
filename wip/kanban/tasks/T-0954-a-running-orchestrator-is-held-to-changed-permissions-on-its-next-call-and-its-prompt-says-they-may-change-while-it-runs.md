---
id: T-0954
type: task
nature: feature
title: A running orchestrator is held to changed permissions on its next call, and its prompt says they may change while it runs
status: done
parent: S-0229
owner: alex
created: 2026-10-05T05:46:25Z
updated: 2026-10-06T21:41:10Z
transitions:
  - to: ready
    at: 2026-10-06T21:30:48Z
    by: agent-S-0229
  - to: in-progress
    at: 2026-10-06T21:30:48Z
    by: agent-S-0229
  - to: done
    at: 2026-10-06T21:41:10Z
    by: agent-S-0229
stream: S-0229
tags: [flai]
touches: [flai/cmd/guard_permissions_change_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
usage:
  source: log
  seconds: 622
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 37
      output: 229
      cache_read: 1001890
      cache_write: 73583
      cost: 0.4707
---
# T-0954 A running orchestrator is held to changed permissions on its next call, and its prompt says they may change while it runs

## Work

The second criterion says a running orchestrator picks up changed permissions on its next decision. S-0218's T-0887 has `flai guard` read `orchestration.permissions` from the manifest on each hook call, and a hook call is a process of its own, so a change should reach the next call already. This task pins that and tells the orchestrator.

- `flai/cmd/guard_permissions_change_test.go`: run the guard as a hook with `FLAI_ROLE=orchestrate` against a fixture manifest, change a permission in the file between two calls, and show the second call's verdict follows the change, both from off to on and from on to off. If the guard turns out to cache the permissions, make it read them per call.
- `flai/internal/harness/harness.go`: S-0218's `orchestratePrompt` says that the operator may change its permissions, policy, and release policy while it runs, that it reads them again before each decision rather than keeping what it read at the start, and that the guard's verdict on a call is what holds.

This task waits for nothing in this story; it needs S-0218 done, which the story's `after` already holds. It runs with the task that adds `flai manifest set`, whose paths it does not share.

## Done when

- the test shows a permission turned on and one turned off between two guard calls change the second verdict
- a harness test finds the sentence about changing permissions in the orchestrator's prompt
- `go test ./cmd/ ./internal/harness/` passes

## Notes
