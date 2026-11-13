---
id: T-3053
type: task
nature: remediation
title: TestADefinedRoleIsHeldToItsFile reads the conditions' Allows from internal/roles instead of copying them
status: done
parent: S-0683
owner: alex
created: 2026-11-10T10:35:09Z
updated: 2026-11-12T15:54:15Z
transitions:
  - to: ready
    at: 2026-11-12T15:52:57Z
    by: agent-S-0683
  - to: in-progress
    at: 2026-11-12T15:52:57Z
    by: agent-S-0683
  - to: done
    at: 2026-11-12T15:54:15Z
    by: agent-S-0683
stream: S-0683
tags: [flai, tests]
touches: [flai/internal/guard/roles_test.go]
usage:
  source: log
  seconds: 78
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 2331
      cache_read: 776081
      cache_write: 14640
      cost: 0.3065
---
# T-3053 TestADefinedRoleIsHeldToItsFile reads the conditions' Allows from internal/roles instead of copying them

## Work

`flai/internal/guard/roles_test.go` (around line 130) asserts a literal copy of the to-backlog and to-ready conditions' `Allows`, as `flai/internal/roles/conditions.go` words them. The guard builds its refusal reason from `Allows` in `flai/internal/guard/roles.go`, so any rewording breaks the test. That is how I-0229 turned main red. Build the expected reason from the `roles` package's conditions, looked up by name, so the test checks how the reason is put together and not the wording. Search `flai/` for any other test that copies a condition's `Allows` and fix it the same way; the planner's explorer found none. This task waits for no task: it shares no file with the others.

## Done when

- The test holds no copy of any condition's `Allows` text.
- Rewording an `Allows` in `flai/internal/roles/conditions.go`, tried locally and not committed, leaves `flai test flai/internal/guard` passing.
- `flai test flai/internal/guard` passes.

## Notes
