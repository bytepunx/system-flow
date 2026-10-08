---
id: T-1372
type: task
nature: feature
title: The host's configuration holds agent.providers, set and reset with flai serve agent provider
status: backlog
parent: S-0349
owner: alex
created: 2026-10-08T08:43:05Z
updated: 2026-10-08T08:54:14Z
transitions: []
stream: S-0349
tags: [cli]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, docs/operators/settings.md, docs/users/flai-reference.md]
after: [T-1371]
---
# T-1372 The host's configuration holds agent.providers, set and reset with flai serve agent provider

## Work

- `config.AgentStart` gains `Providers map[string]manifest.Provider` (`json:"providers,omitempty"`), validated with T-1371's checks when the configuration is read and written.
- `flai serve agent provider <name> --api <api> --base-url <url> --key-env <VAR> [--model <alias>=<name>]...` sets an entry, `--reset` removes it, and with no flags it prints the entry. No flag takes a key; the help says the key stays in the variable `--key-env` names.
- Rows in `docs/operators/settings.md` for `agent.providers.<name>.api`, `.base_url`, `.key_env`, and `.models`, and the command's flags, so that `TestSettingsIndexIsComplete` passes.

It waits for T-1371, whose `Provider` and validation it stores, and touches `docs/operators/settings.md` after it.

## Done when

- `config_test.go` covers a round trip and each refusal; `serve_actions_test.go` covers set, print, and reset.
- `flai test flai/internal/config/ flai/cmd/ docs/operators/settings.md` passes.

## Notes

Layer 2 of S-0349.
