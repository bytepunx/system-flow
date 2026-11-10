---
id: T-2930
type: task
nature: improvement
title: ready-empty counts only the ready stories not held on after, so held stories in ready no longer keep the orchestrator from promoting
status: done
parent: S-0673
owner: alex
created: 2026-11-10T08:49:13Z
updated: 2026-11-10T10:19:25Z
transitions:
  - to: ready
    at: 2026-11-10T09:53:41Z
    by: agent-S-0673
  - to: in-progress
    at: 2026-11-10T09:53:43Z
    by: agent-S-0673
  - to: done
    at: 2026-11-10T10:19:25Z
    by: agent-S-0673
stream: S-0673
tags: [flai, serve]
touches: [flai/internal/serve/readyempty.go, flai/internal/serve/readyempty_test.go, flai/internal/serve/agents_launch.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/inbox/events.go, flai/internal/mcpserver/server.go, design/system/agent-narrative.md, design/system/flai-cli.md, design/system/strategic-agents.md]
after: [T-2928]
usage:
  source: log
  seconds: 1542
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 5729
      cache_read: 4267931
      cache_write: 209839
      cost: 2.167
---
# T-2930 ready-empty counts only the ready stories not held on after, so held stories in ready no longer keep the orchestrator from promoting

## Work

In `look` (`flai/internal/serve/agents_launch.go`), pass `noteReadyEmpty` the count of ready stories that `holds` does not hold on `after` (`workitem.HoldAfter`), in place of `len(stories)`. A story held on overlap alone still counts. Rewrite the comments of `readyDrained` and `noteReadyEmpty` in `flai/internal/serve/readyempty.go` to say so, and the reason's words in `orchestrate_reasons.go` ("the ready column holds n") to say what n counts.

Add a test in `readyempty_test.go` that reproduces I-0285: ready holds as many stories held on `after` as the in-progress room, and one look tells `ready-empty`; with the same stories not held, it tells none. Keep the existing cases passing.

Say the same in the event's and the tools' descriptions, `flai/internal/inbox/events.go` and the `inbox` and `wait_for_events` descriptions in `flai/internal/mcpserver/server.go`, and in `design/system/agent-narrative.md`, `design/system/flai-cli.md`, and `design/system/strategic-agents.md`, where each reads "held stories counted".

It waits for T-2928, whose ADR decides which holds the count leaves out.

## Done when

- The new test fails before the change and passes after it.
- `flai test` on the changed paths passes.
- No description of `ready-empty` in the touched files says held stories are counted.

## Notes
