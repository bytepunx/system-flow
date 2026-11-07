---
id: T-2578
type: task
nature: remediation
title: flai's built-in default agent, claude-code with no model, is what an agent falls back to when neither it, the project's default, nor the host's command names a harness
status: done
parent: S-0587
owner: alex
created: 2026-11-07T10:11:59Z
updated: 2026-11-07T10:51:29Z
transitions:
  - to: ready
    at: 2026-11-07T10:31:50Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:31:52Z
    by: system-flow
  - to: done
    at: 2026-11-07T10:51:29Z
    by: agent-S-0587
stream: S-0587
tags: [cli]
touches: [flai/internal/manifest/agent.go, flai/internal/manifest/agent_test.go, flai/internal/workitem/agent.go, flai/internal/workitem/agent_test.go]
after: [T-2577]
usage:
  source: log
  seconds: 1177
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 22
      output: 123
      cache_read: 2987510
      cache_write: 25310
      cost: 1.0807
---
# T-2578 flai's built-in default agent, claude-code with no model, is what an agent falls back to when neither it, the project's default, nor the host's command names a harness

## Work

It waits for T-2577: ADR-0261 § 1 fixes the built-in agent and the order an agent resolves in.

- In `flai/internal/manifest/agent.go`, add the built-in default agent, harness `claude-code` (the adapter's name) with no model and no config, and a method on `*Agent`, such as `OrBuiltin`, that answers the agent itself when it names a harness and otherwise the built-in with the agent's own words laid over it (`Agent.With`). Doc comments name ADR-0261.
- In `flai/internal/workitem/agent.go`, keep `Repo.AgentOf` as it is (story's own, then the project's default, ADR-0251), and add `Repo.StartAgentOf(it, commandSet bool)`: `AgentOf`, then, when that names no harness and `commandSet` is false, the built-in under it. With `commandSet` true it is `AgentOf`, so the host's command still starts a story whose agent and project name none.
- `flai story new` (`flai/internal/workitem/create.go`) is unchanged: it copies only the project's default, never the built-in.
- A story whose agent names a harness, `none` included (ADR-0242), keeps its agent whole.

## Done when

- Tests in `flai/internal/manifest/agent_test.go` cover `OrBuiltin` for nil, an agent with only a model, and an agent naming a harness.
- Tests in `flai/internal/workitem/agent_test.go` cover `StartAgentOf` for a project with no `agent:` and no command (built-in harness, the story's model kept), with a command set (no harness, as before), and for a story naming a harness.
- `flai test` on the changed paths passes.

## Notes
