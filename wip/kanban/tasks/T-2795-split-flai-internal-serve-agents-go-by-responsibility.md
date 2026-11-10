---
id: T-2795
type: task
nature: improvement
title: Split flai/internal/serve/agents.go by responsibility
status: done
parent: S-0631
owner: alex
created: 2026-11-09T22:04:32Z
updated: 2026-11-10T06:38:40Z
transitions:
  - to: ready
    at: 2026-11-10T06:16:43Z
    by: agent-S-0631
  - to: in-progress
    at: 2026-11-10T06:16:44Z
    by: agent-S-0631
  - to: done
    at: 2026-11-10T06:38:40Z
    by: agent-S-0631
stream: S-0631
tags: [cli]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_state.go, flai/internal/serve/agents_launch.go, flai/internal/serve/agents_restart.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents_activity.go, flai/internal/mcpserver/server.go]
usage:
  source: log
  seconds: 1285
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 25
      output: 171
      cache_read: 3046885
      cache_write: 298545
      cost: 1.2303
---
# T-2795 Split flai/internal/serve/agents.go by responsibility

## Work

`agents.go` is 3,298 lines. Move its declarations into new files of package `serve` by responsibility. What any of them does stays the same. The line ranges below are approximate, from the plan's read on 2026-11-09:

| Lines | Declarations | File |
|-------|--------------|------|
| 74–443 | `AgentConfig` and its methods, `HarnessLifts`, `startAllowed`, the constants at 226–236, `AgentRun` and its predicates | `agents.go`, which keeps them |
| 444–854 | `AgentState`, `RoleState`, `RoleStates`, `put`, `keptAgents`, `AgentStates`, `parseAgents`, `updateAgent`, the clone helpers | `agents_state.go` |
| 855–1728 | `launcher`, `told`, `ended`, `readyStories`, `judgeRun`, `endRun`, `blockAsked`, the conversation helpers, `look`, `askHolds`, `settleOrphans` | `agents_launch.go` |
| 1729–2156 | `resume`, `startSentBack`, `startAcceptSentBack`, `sentBackFor`, `restartEnded`, `autoRestart`, `stopRestarting`, `countedRuns`, `awaits` | `agents_restart.go` |
| 2157–2823 | `agentEnv`, `mcpEnv`, `start`, `holdStart`, `startWhy`, `startOf`, `provide`, `sharedWith`, `openLog`, `spawnWith`, `await` | `agents_start.go` |
| 2824–3298 | `StoryActivity`, `Elsewhere`, `Activity`, `Unfinished`, `runActivity`, `held`, `asking`, `answered`, `messaged`, `newSession` | `agents_activity.go` |

The package already splits `launcher`'s methods by concern: `role.go`, `container.go`, `restart.go`, `stop.go`, `agent_events.go`, `orchestrate.go`, `plan.go`, and others. Keep to that. Do not put agent status into `activity.go`, which holds the activity logs.

Fix the comment at `flai/internal/mcpserver/server.go:1041` that names `internal/serve/agents.go` if what it points to moved.

Change no test file. `agents_test.go` stays as it is.

This task waits for nothing and shares no path with the other splits, so all of them run together.

## Done when

- No non-test file in `flai/internal/serve` is over 2,000 lines, and none of the new files is over 1,200.
- `git diff --stat` shows no `_test.go` file changed.
- `flai test flai/internal/serve flai/internal/mcpserver/server.go` passes, gofmt, goimports, and golangci-lint included.

## Notes
