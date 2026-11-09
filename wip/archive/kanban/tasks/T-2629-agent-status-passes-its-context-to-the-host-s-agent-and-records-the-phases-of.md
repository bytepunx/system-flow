---
id: T-2629
type: task
nature: improvement
title: agent.status passes its context to the host's Agent and records the phases of the state it builds
status: done
parent: S-0603
owner: alex
created: 2026-11-07T23:47:58Z
updated: 2026-11-08T01:47:36Z
transitions:
  - to: ready
    at: 2026-11-08T01:14:11Z
    by: system-flow
  - to: in-progress
    at: 2026-11-08T01:14:14Z
    by: system-flow
  - to: done
    at: 2026-11-08T01:47:36Z
    by: agent-S-0603
stream: S-0603
tags: [cli]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/writes_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go]
usage:
  source: log
  seconds: 2002
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 5300
      cache_read: 6142983
      cache_write: 94518
      cost: 1.9107
---
# T-2629 agent.status passes its context to the host's Agent and records the phases of the state it builds

## Work

Criterion 3, its phases. `hostapi.Host.Agent` is `func(root string) any` (`flai/internal/hostapi/writes.go`), so the `agent.status` method in `flai/internal/hostapi/hostapi.go` tracks one phase, `agent.state`, and nothing inside it.

- Give `Host.Agent` the request's context: `func(ctx context.Context, root string) any`, and pass it from `agent.status`. Keep the `agent.state` phase around the whole call.
- In `flai/cmd/serve_actions.go`, track a phase with `perf.Track` around each part of the state: reading the agent states, `serve.Activity`, and `serve.RoleStates`, named `agent.states`, `agent.activity`, and `agent.roles`.
- Update the callers and the stub in `flai/internal/hostapi/writes_test.go`.

Waits for nothing: it changes none of the files of T-2621 and T-2623, so it runs in the first layer beside T-2621.

## Done when

- A test shows `agent.status` records `agent.state` and the three phases inside it on the request's recorder.
- `flai test` on the changed paths passes.

## Notes
