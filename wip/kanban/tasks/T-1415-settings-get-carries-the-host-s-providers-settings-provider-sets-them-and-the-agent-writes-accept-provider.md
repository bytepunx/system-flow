---
id: T-1415
type: task
nature: feature
title: settings.get carries the host's providers, settings.provider sets them, and the agent writes accept provider
status: backlog
parent: S-0359
owner: alex
created: 2026-10-08T08:53:09Z
updated: 2026-10-08T08:53:09Z
transitions: []
stream: S-0359
tags: [cli]
touches: [flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go]
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
