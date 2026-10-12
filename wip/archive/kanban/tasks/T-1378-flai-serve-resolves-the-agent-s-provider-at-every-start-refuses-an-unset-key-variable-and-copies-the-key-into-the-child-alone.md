---
id: T-1378
type: task
nature: feature
title: flai serve resolves the agent's provider at every start, refuses an unset key variable, and copies the key into the child alone
status: done
parent: S-0350
owner: alex
created: 2026-10-08T08:44:04Z
updated: 2026-10-12T17:04:00Z
transitions:
  - to: ready
    at: 2026-10-12T16:57:05Z
    by: agent-S-0350
  - to: in-progress
    at: 2026-10-12T16:57:05Z
    by: agent-S-0350
  - to: done
    at: 2026-10-12T17:04:00Z
    by: agent-S-0350
stream: S-0350
tags: [cli]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/plan.go, flai/internal/serve/orchestrate.go, flai/internal/serve/analyze.go, flai/cmd/serve_actions.go]
after: [T-1376]
usage:
  source: log
  seconds: 415
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 54
      output: 18237
      cache_read: 2903547
      cache_write: 126038
      cost: 1.7408
---
# T-1378 flai serve resolves the agent's provider at every start, refuses an unset key variable, and copies the key into the child alone

## Work

- `serve.AgentConfig` in `agents.go` gains `Providers map[string]manifest.Provider`, filled from the host's `agent.providers` by `agentConfig` in `flai/cmd/serve_actions.go` at every look, as `Harnesses` is.
- Before each `adapter.Start` in `agents.go` (a story's agent), `plan.go`, `orchestrate.go`, and `analyze.go`, resolve the agent's provider with `manifest.Resolve` from `AgentConfig.Providers` and the project's manifest `providers`, and pass it in the `Request`. A refusal fails the start as an unknown harness fails it today, with the same reporting.
- In `spawn`, for each `EnvFrom` entry, read the source variable from `flai serve`'s environment; when it is unset or empty, refuse the start naming the variable and `flai serve agent provider <name> --key-env`. Otherwise append the value to `cmd.Env` only, never to the run's log, the start line, or the start record.
- The run's start record keeps the provider's name and the key's variable name, so that a later measurement prices the run by the gateway it ran through.
- A test starts a fake agent with a provider and a key in the environment, checks the child received it, and reads the run log and the start record to find the provider's name and no key; another starts with the variable unset and gets the refusal.

It waits for T-1376, whose `Request.Provider` and `Start.EnvFrom` it fills and reads.

## Done when

- The tests above pass, and a start with no provider is unchanged.
- `flai test flai/internal/serve/ flai/cmd/` passes.

## Notes

Layer 2 of S-0350.
