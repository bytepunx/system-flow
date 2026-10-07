---
id: T-0917
type: task
nature: feature
title: flai serves the orchestrator's and the analyzer's run streams, and stops and starts the orchestrator on the operator's word
status: done
parent: S-0228
owner: alex
created: 2026-10-05T05:44:45Z
updated: 2026-10-07T00:23:18Z
transitions:
  - to: ready
    at: 2026-10-07T00:03:28Z
    by: agent-S-0228
  - to: in-progress
    at: 2026-10-07T00:03:28Z
    by: agent-S-0228
  - to: done
    at: 2026-10-07T00:23:18Z
    by: agent-S-0228
stream: S-0228
tags: [dashboard]
touches: [flai/internal/hostapi, flai/internal/serve/stream.go, flai/internal/serve/stream_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/agents.go, flai/cmd/serve.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/orchestrate_hold_test.go, flai/internal/serve/agents_test.go]
usage:
  source: log
  seconds: 1190
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 268
      output: 119355
      cache_read: 18330223
      cache_write: 443706
      cost: 8.5438
---
# T-0917 flai serves the orchestrator's and the analyzer's run streams, and stops and starts the orchestrator on the operator's word

## Work

What the Orchestrator and Analyzer pages need from flai serve beyond what S-0259, S-0218, and S-0223 already give. `activity.document` answers `orchestrator` and `analyzer` already (S-0206), so this adds no read of the documents.

- `agent.stream {role}`: given `role`, `orchestrate` or `analyze`, in place of `story` or `plan`, the stream of the newest run of that kind that flai serve started (the `orchestrator` run S-0218 keeps in `serve/agents.json`, the analyzer runs S-0223 keeps), read as a planner run's is (`flai/internal/serve/stream.go`). Exactly one of `story`, `plan`, and `role` is given. No run of the kind is not found.
- `orchestrate.stop` and `orchestrate.start`, writes gated on the `orchestrate` host action, each running a `flai serve orchestrate stop|start` subcommand. Stop ends the current run as `stop.go` ends a story's agent, logs its end in `orchestrator.md` with `stopped: stopped from the dashboard`, and holds the orchestrator stopped in `agents.json`, so `flai serve` does not start it again while the action stays on. Start lifts the hold, and the orchestrator starts as it does when the action is turned on. Turning the action off and on also lifts the hold. `agent.status` carries `held` on the orchestrator. Each is journalled, and the dashboard hears an `agent` notification.

It waits for no task: nothing of the dashboard is needed to build or test it. It shares no path with the API task that runs beside it.

## Done when

- [ ] `agent.stream` with `role` reads the newest run of the kind, refuses more or fewer than one of `story`, `plan`, and `role`, and answers not found with no run, with tests in `flai/internal/serve` and `flai/internal/hostapi`
- [ ] `orchestrate.stop` stops the run and holds it stopped while the action is on, `orchestrate.start` lifts the hold, both are refused while the action is off, and tests pin each on a fake harness
- [ ] `go test -race -short` passes for the packages it changed

## Notes
