---
id: T-1408
type: task
nature: feature
title: A provider entry may name spend_key_env, the variable of a key that can read the gateway's spend log
status: backlog
parent: S-0357
owner: alex
created: 2026-10-08T08:50:50Z
updated: 2026-10-08T08:54:20Z
transitions: []
stream: S-0357
tags: [cli]
touches: [flai/internal/manifest/provider.go, flai/internal/manifest/provider_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, docs/operators/settings.md, docs/users/flai-reference.md]
---
# T-1408 A provider entry may name spend_key_env, the variable of a key that can read the gateway's spend log

## Work

- `manifest.Provider` gains `SpendKeyEnv string` (`spend_key_env`), checked as `key_env` is; like `key_env` it is the host's alone and refused in a manifest entry.
- `flai serve agent provider <name> --spend-key-env <VAR>`, beside S-0349's `--key-env` in `serve_actions.go`; empty clears it.
- Rows in `docs/operators/settings.md` for `agent.providers.<name>.spend_key_env` and the flag, so that `TestSettingsIndexIsComplete` passes.

First layer.

## Done when

- `provider_test.go` covers the new key and its refusal in a manifest entry; `serve_actions_test.go` sets and clears it.
- `flai test flai/internal/manifest/ flai/cmd/ docs/operators/settings.md` passes.

## Notes

Layer 1 of S-0357. ADR-0132's consequences: the entry gains the name where the agent's key cannot read the spend log.
