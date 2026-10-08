---
id: T-1371
type: task
nature: feature
title: The manifest's agent takes provider, the manifest's providers map overrides api, base_url, and models, and one resolver merges it with the host's entry
status: backlog
parent: S-0349
owner: alex
created: 2026-10-08T08:42:58Z
updated: 2026-10-08T08:42:58Z
transitions: []
stream: S-0349
tags: [cli]
touches: [flai/internal/manifest/provider.go, flai/internal/manifest/provider_test.go, flai/internal/manifest/agent.go, flai/internal/manifest/agent_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, docs/operators/settings.md]
---
# T-1371 The manifest's agent takes provider, the manifest's providers map overrides api, base_url, and models, and one resolver merges it with the host's entry

## Work

- `manifest.Agent` gains `Provider string` (`yaml:"provider,omitempty"`), checked by the harness name pattern in `agent.go`; `With` carries it over as it carries `harness`. A story's front matter takes it through the same struct.
- A new `provider.go` holds `Provider` (`API`, `BaseURL`, `KeyEnv`, `Models map[string]string`), the three API names as constants, and its validation: an unknown `api`, a `base_url` that is not an absolute http or https URL, a `key_env` that is not `^[A-Za-z_][A-Za-z0-9_]*$`, an empty alias or model name.
- `Manifest` gains `Providers map[string]Provider`; its validation refuses `key_env` in an entry, saying that `key_env` is the host's and naming `flai serve agent provider <name> --key-env`.
- `Resolve(name string, host, project map[string]Provider) (*Provider, error)`: an empty name resolves to nil; a name the host lacks is refused naming `flai serve agent provider <name>`; otherwise the host's entry with the manifest's `api`, `base_url`, and `models` laid over it, `models` merged key by key.
- Rows in `docs/operators/settings.md` for `agent.provider` and `providers.<name>.*` in the manifest, so that `TestSettingsIndexIsComplete` passes in this task.

First layer: nothing in the story waits for less.

## Done when

- `provider_test.go` covers the merge, each refusal, and an empty name; `agent_test.go` and `manifest_test.go` cover the new field and the refused `key_env`.
- `flai test flai/internal/manifest/ docs/operators/settings.md` passes.

## Notes

Layer 1 of S-0349.
