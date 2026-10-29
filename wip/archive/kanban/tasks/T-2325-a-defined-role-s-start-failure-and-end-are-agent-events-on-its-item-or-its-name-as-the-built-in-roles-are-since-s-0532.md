---
id: T-2325
type: task
nature: feature
title: A defined role's start, failure, and end are agent events on its item or its name, as the built-in roles' are since S-0532
status: done
parent: S-0417
owner: alex
created: 2026-10-29T05:30:52Z
updated: 2026-10-29T05:33:30Z
transitions:
  - to: ready
    at: 2026-10-29T05:31:05Z
    by: agent-S-0417
  - to: in-progress
    at: 2026-10-29T05:31:06Z
    by: agent-S-0417
  - to: done
    at: 2026-10-29T05:33:30Z
    by: agent-S-0417
stream: S-0417
tags: []
touches: [flai/internal/serve/role.go, flai/internal/serve/agent_events.go, flai/internal/serve/role_test.go, flai/internal/serve/agent_events_test.go, design/system/strategic-agents.md]
after: [T-1739]
usage:
  source: log
  seconds: 144
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 65
      output: 22481
      cache_read: 5126907
      cache_write: 200296
      cost: 2.5512
---

# T-2325 A defined role's start, failure, and end are agent events on its item or its name, as the built-in roles' are since S-0532

## Work

S-0532, accepted while S-0417 ran, made each built-in role's start, failure, and end an `agent` event (`appendAgentEvent`, `appendAgentEnd` in `flai/internal/serve/agent_events.go`). `RunRole` emits none, and `runConcerns` would call a defined role's run on an item the planner's. Make a defined role's run an event too: on its item with its verb as role, or on its name as `inbox.RoleType` when it runs on none. Test it in `role_test.go`.

Waits for T-1739, after which S-0532 was synced in.

## Done when

`flai test` on the changed files passes, and `role_test.go` shows the started, failed, and ended events.

## Notes
