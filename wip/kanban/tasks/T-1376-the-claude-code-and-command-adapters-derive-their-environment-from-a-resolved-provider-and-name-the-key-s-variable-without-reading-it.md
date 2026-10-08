---
id: T-1376
type: task
nature: feature
title: The claude-code and command adapters derive their environment from a resolved provider and name the key's variable without reading it
status: backlog
parent: S-0350
owner: alex
created: 2026-10-08T08:43:59Z
updated: 2026-10-08T08:43:59Z
transitions: []
stream: S-0350
tags: [cli]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go]
---
# T-1376 The claude-code and command adapters derive their environment from a resolved provider and name the key's variable without reading it

## Work

- `harness.Request` gains `Provider *manifest.Provider` and `ProviderName string`, resolved by the caller.
- `harness.Start` gains `EnvFrom map[string]string`: the child's variable to set, mapped to the variable of `flai serve`'s environment to copy it from. The adapter never reads a key; `Env` and `Argv` carry no secret.
- `claude-code`: with a provider whose `api` is not `anthropic-messages`, `Start` refuses, naming the provider and `flai serve agent provider <name> --api`. Otherwise `Env` gets `ANTHROPIC_BASE_URL`, `ANTHROPIC_API_KEY=` (empty), and `ANTHROPIC_DEFAULT_OPUS_MODEL`, `ANTHROPIC_DEFAULT_SONNET_MODEL`, and `ANTHROPIC_DEFAULT_HAIKU_MODEL` from `models`; `EnvFrom` maps `ANTHROPIC_AUTH_TOKEN` to `key_env`. The strategic starts, which use `claudeCodeStrategic`, get the same.
- `command`: `FLAI_PROVIDER`, `FLAI_PROVIDER_API`, `FLAI_PROVIDER_BASE_URL`, `FLAI_PROVIDER_KEY_ENV`, and `{provider}` replaced in the operator's arguments; no `EnvFrom`.

First layer.

## Done when

- `harness_test.go` covers each adapter with and without a provider, the refused `api`, and an alias missing from `models`, which sets no variable.
- `flai test flai/internal/harness/` passes.

## Notes

Layer 1 of S-0350.
