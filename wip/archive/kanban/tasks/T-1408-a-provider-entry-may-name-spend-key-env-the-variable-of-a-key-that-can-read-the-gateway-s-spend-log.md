---
id: T-1408
type: task
nature: feature
title: A provider entry may name spend_key_env, the variable of a key that can read the gateway's spend log
status: done
parent: S-0357
owner: alex
created: 2026-10-08T08:50:50Z
updated: 2026-10-13T19:25:43Z
transitions:
  - to: ready
    at: 2026-10-13T19:17:09Z
    by: agent-S-0357
  - to: in-progress
    at: 2026-10-13T19:17:09Z
    by: agent-S-0357
  - to: done
    at: 2026-10-13T19:25:43Z
    by: agent-S-0357
stream: S-0357
tags: [cli]
touches: [flai/internal/manifest/provider.go, flai/internal/manifest/provider_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, docs/operators/settings.md, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 514
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 26
      output: 13925
      cache_read: 1709657
      cache_write: 51993
      cost: 0.9331
---
# T-1408 A provider entry may name spend_key_env, the variable of a key that can read the gateway's spend log

## Work

- `manifest.Provider` gains `SpendKeyEnv string` (`spend_key_env`, omitempty), checked as `key_env` is by `Problems`; like `key_env` it is the host's alone and refused in a manifest entry (`ProviderErrors`); `ResolveProvider` keeps the host's; it is optional on a host entry.
- `flai serve agent provider <name> --spend-key-env <VAR>`, beside `--key-env` in `serve_actions.go`; an empty value clears it; the listing names it where set, and the JSON listing gives it with whether the variable is set, never its value.
- Regenerate `docs/operators/settings.md` and `docs/users/flai-reference.md` with `make flai-reference` (never hand-edit their generated tables: S-0378 shares `docs/operators/settings.md` and regenerates it too); add the row for `agent.providers.<name>.spend_key_env` where the settings index keeps hand-written rows for configuration keys, so that `TestSettingsIndexIsComplete` passes.

First layer.

## Done when

- `provider_test.go` covers the new key and its refusal in a manifest entry; `serve_actions_test.go` sets and clears it.
- `flai test flai/internal/manifest/ flai/cmd/ docs/operators/settings.md` passes.

## Notes

Layer 1 of S-0357. ADR-0132's consequences: the entry gains the name where the agent's key cannot read the spend log. Meets criterion 5.
