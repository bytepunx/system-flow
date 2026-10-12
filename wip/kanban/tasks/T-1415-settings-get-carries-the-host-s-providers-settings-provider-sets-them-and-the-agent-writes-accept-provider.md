---
id: T-1415
type: task
nature: feature
title: settings.get carries the host's providers, settings.provider sets them, and the agent writes accept provider
status: done
parent: S-0359
owner: alex
created: 2026-10-08T08:53:09Z
updated: 2026-10-12T17:04:54Z
transitions:
  - to: ready
    at: 2026-10-12T16:56:00Z
    by: agent-S-0359
  - to: in-progress
    at: 2026-10-12T16:56:00Z
    by: agent-S-0359
  - to: done
    at: 2026-10-12T17:04:54Z
    by: agent-S-0359
stream: S-0359
tags: [cli]
touches: [flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/agent.ts]
usage:
  source: log
  seconds: 534
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 96
      output: 521
      cache_read: 4418570
      cache_write: 133988
      cost: 2.0546
---
# T-1415 settings.get carries the host's providers, settings.provider sets them, and the agent writes accept provider

## Work

- `hostSettings` in `serve_actions.go`: `agent.providers` with each entry's `api`, `base_url`, `key_env`, `models`, and `key_set`, whether the variable is set and not empty in this process's environment; never the value.
- `settings.provider` in `hostapi/settings.go`: `{ name, api, base_url, key_env, models }` or `{ name, reset: true }`, host-wide, running `flai serve agent provider` with flags built from it, never through a shell; the name checked by the harness name pattern.
- `settings.default_agent`, `settings.manifest`, and the item writes in `writes.go` that take an agent accept `provider` and pass `--provider`; their refusal text lists it.
- Tests in `serve_actions_test.go` for the read and in `writes_test.go` for the write, a reset, a refusal, and an agent with a provider.

First layer.

## Done when

- `flai test flai/cmd/ flai/internal/hostapi/` passes.

## Notes

Layer 1 of S-0359.
