---
id: T-2455
type: task
nature: feature
title: flai serve starts, restarts, resumes, and commits a story's agent with the agent Repo.AgentOf resolves, so a story made before the default was set starts
status: done
parent: S-0573
owner: alex
created: 2026-11-03T17:28:57Z
updated: 2026-11-03T17:45:26Z
transitions:
  - to: ready
    at: 2026-11-03T17:36:56Z
    by: system-flow
  - to: in-progress
    at: 2026-11-03T17:36:56Z
    by: system-flow
  - to: done
    at: 2026-11-03T17:45:26Z
    by: agent-S-0573
stream: S-0573
tags: []
touches: [flai/internal/serve/start.go, flai/internal/serve/restart.go, flai/internal/serve/commit.go, flai/internal/serve/agents.go, flai/internal/serve/usage.go, flai/internal/serve/default_agent_test.go]
after: [T-2454]
usage:
  source: log
  seconds: 510
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 41
      output: 9874
      cache_read: 2194995
      cache_write: 67458
      cost: 1.089
---
# T-2455 flai serve starts, restarts, resumes, and commits a story's agent with the agent Repo.AgentOf resolves, so a story made before the default was set starts

## Work

Read a story's agent through `repo.AgentOf(it)` wherever `flai serve` reads it to start, restart, resume, send back, commit, or show a story's agent: `start.go` (the harness check, `heldStart`), `restart.go` (`startNow`, the restart refusal, `heldStart`), `commit.go`, `agents.go` (the ready list, the answered, sent-back, and auto-restart starts, the stand-in shown for a story), and `usage.go` (`logReader`). Add `default_agent_test.go` reproducing S-0573's case: a ready story with no `agent` in a project whose `system-flow.yaml` has a default naming a harness, and no host command, is started by `Start` and by the launcher, not refused with "names no harness". Waits for T-2454, whose `Repo.AgentOf` it calls.

## Done when

- The new test fails against main's serve code and passes with the change.
- A story whose agent names a harness still starts with its own agent.
- A story with no agent and no default is still refused, or run by the host's command.
- `flai test` on the changed serve files passes.

## Notes
