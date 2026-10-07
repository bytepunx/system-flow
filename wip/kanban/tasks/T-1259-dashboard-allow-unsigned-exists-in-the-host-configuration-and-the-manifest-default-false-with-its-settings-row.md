---
id: T-1259
type: task
nature: feature
title: dashboard.allow_unsigned exists in the host configuration and the manifest, default false, with its settings row
status: backlog
parent: S-0239
owner: alex
created: 2026-10-07T23:09:01Z
updated: 2026-10-07T23:09:01Z
transitions: []
stream: S-0239
tags: [cli]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, docs/operators/settings.md]
---
# T-1259 dashboard.allow_unsigned exists in the host configuration and the manifest, default false, with its settings row

## Work

The first half of criterion 1. It waits for nothing: every other flai task reads the setting this one adds.

- `flai/internal/config/config.go`: `Dashboard.AllowUnsigned bool` as `allow_unsigned`, beside `no_restart`, and `dashboard.allow_unsigned` in the keys `flai config set` and `get` accept. Default false.
- `flai/internal/manifest/manifest.go`: the manifest's `dashboard` block gains `allow_unsigned`, optional, so that a project may set it on or off over the host's value. One function gives the value in force for a project from the two: the manifest's when it sets one, else the host configuration's.
- `docs/operators/settings.md`: the row for `dashboard.allow_unsigned` in the configuration and in the manifest, saying it is off by default, what it allows (ADR-0070), and that it is shown wherever it is in force. `flai/cmd/settings_doc_test.go` fails without the row, so it lands here, not with the docs task.
- Tests: `config_test.go` for the key's default and `config set`; `manifest_test.go` for the manifest's value over the host's, both ways, and the host's when the manifest is silent.

## Done when

- `flai config get dashboard.allow_unsigned` answers false on a fresh configuration and true after `flai config set dashboard.allow_unsigned true`.
- A manifest's `dashboard.allow_unsigned` overrides the host's value for that project, and its absence leaves the host's.
- `docs/operators/settings.md` has the row, and the tests pass with `flai test` on the changed paths.

## Notes

Drafted by the planner.
