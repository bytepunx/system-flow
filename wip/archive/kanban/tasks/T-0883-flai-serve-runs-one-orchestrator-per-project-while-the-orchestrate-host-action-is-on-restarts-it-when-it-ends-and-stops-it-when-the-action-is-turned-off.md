---
id: T-0883
type: task
nature: feature
title: flai serve runs one orchestrator per project while the orchestrate host action is on, restarts it when it ends, and stops it when the action is turned off
status: done
parent: S-0218
owner: alex
created: 2026-10-05T04:46:09Z
updated: 2026-10-05T07:58:59Z
transitions:
  - to: ready
    at: 2026-10-05T07:47:04Z
    by: agent-S-0218
  - to: in-progress
    at: 2026-10-05T07:47:05Z
    by: agent-S-0218
  - to: review
    at: 2026-10-05T07:58:59Z
    by: agent-S-0218
  - to: done
    at: 2026-10-05T07:58:59Z
    by: agent-S-0218
stream: S-0218
tags: [flai]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/agents.go, flai/internal/serve/serve.go, flai/internal/serve/projects.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go]
after: [T-0881, T-0882]
usage:
  source: log
  seconds: 714
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 132
      output: 54519
      cache_read: 7771400
      cache_write: 197490
      cost: 3.7504
---
# T-0883 flai serve runs one orchestrator per project while the orchestrate host action is on, restarts it when it ends, and stops it when the action is turned off

## Work

Add the host action and the run's lifecycle, mirroring the planner's (`flai/internal/serve/plan.go`, `replan.go`) but with no item and no end of its own:

- `flai/internal/hostapi/writes.go`: `ActionOrchestrate = "orchestrate"` in `Actions`, off by default, with a description naming `orchestration.agent`, and seen by the dashboard as `plan` is (`DashboardSees`). `flai serve enable orchestrate` and `disable` then work through `knownAction`.
- `flai/cmd/serve_actions.go`: `AgentConfig` gains `Orchestrate`, from `cfg.ActionEnabled(hostapi.ActionOrchestrate, root)`.
- `flai/internal/serve/orchestrate.go`: an orchestrator per served project that looks whenever the project's launcher does. With the action on and no run, it starts one through the harness adapter with a `harness.Request` whose `Role` is `orchestrate`, in the main checkout, as `orchestrator`, with `FLAI_ROLE=orchestrate`, `FLAI_SESSION`, and `FLAI_STARTED_BY=flai-serve`, and the agent `repo.Manifest.OrchestrationAgent()` gives. When the run ends with the action still on, it starts another; after a failed exit, no sooner than a minute later, so a run that fails at once does not spin. With the action off and a run going, it stops the run as `stop.go` stops a story's agent. Refuse to start, saying why, when the agent names no harness and no command is set on the host, or `orchestrator.md` is missing on `claude-code`.
- `flai/internal/serve/agents.go`: the state in `serve/agents.json` keeps the newest run under `orchestrator`, with the fields a planner's run has and no `item`; its log is `<key>-orchestrator-<start>.log`. The in-progress limit does not count it. `agent.status` carries it, and each start, failure, stop, and end is a journal entry with action `orchestrate` and an `agent` notification to the dashboard.
- When a run ends or is stopped, `serve.LogRunEnd` logs the time since its last `activity_log` in `wip/agents/orchestrator.md`, with the run's final reply, or `stopped: orchestrate turned off`, as the summary.

This task waits for T-0881, for `OrchestrationAgent()`, and for T-0882, for the request role it starts. It runs with T-0887, whose paths it does not share.

## Done when

- with `orchestrate` off nothing starts; turned on, one run starts for the project, and a test pins it on a fake harness
- a run that ends while the action is on is started again, after a minute when it failed, and a test pins both
- turning the action off stops the run and logs its end in `orchestrator.md`, and a test pins it
- `agents.json` and `agent.status` carry the run under `orchestrator`, and the journal has its entries
- `go test ./internal/serve/ ./internal/hostapi/ ./cmd/` passes

## Notes
