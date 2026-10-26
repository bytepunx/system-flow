---
id: T-2142
type: task
nature: improvement
title: The manifest's orchestration.permissions.confirm_after is declared as a duration, off by default, validated by flai check, settable with flai manifest set, shown in the settings catalog, and indexed in the operator docs
status: backlog
parent: S-0502
owner: alex
created: 2026-10-26T04:43:18Z
updated: 2026-10-26T04:45:54Z
transitions: []
stream: S-0502
tags: [flai]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/orchestration_test.go, docs/operators/settings.md, design/system/project-manifest.md]
after: [T-2141]
---
# T-2142 The manifest's orchestration.permissions.confirm_after is declared as a duration, off by default, validated by flai check, settable with flai manifest set, shown in the settings catalog, and indexed in the operator docs

## Work

The permission T-2141's ADR names, point 1, declared as `answer_threads` is, the one permission that is not a boolean today:

- `flai/internal/manifest/manifest.go`: `ConfirmAfter string` on `Permissions` (`yaml:"confirm_after,omitempty"`), the constant and its entry in `PermissionNames`, `UnmarshalYAML` keeping it, `ConfirmAfterDuration()` that parses it as `Digest.CompactAfterDuration()` parses `compact_after` and gives zero when unset, `errors()` reporting a value that is not a positive Go duration as `manifest.orchestration`, and `Allows()` treating it as on when the duration is positive. Tests in `manifest_test.go` beside `TestOrchestrationPermissions`: unset, `48h`, `0`, `two days`.
- `flai/internal/manifest/settings.go`: a `permissionText` entry with its meaning (flai serve confirms an orchestrator recommendation the operator has not answered in this long, save one marked for the operator alone and those the ADR keeps for the operator) and its risk (a decision the orchestrator recommended becomes the answer without the operator reading it); the catalog entry as `KindDuration`, special-cased beside `answer_threads`, so `flai manifest set orchestration.permissions.confirm_after 48h` sets it and the dashboard's `StrategicSettings` renders its duration input with no change of its own. `settings_test.go`'s `TestSettingsCatalogAgreesWithTheValidation` covers the new name and its risk.
- `flai/internal/check/orchestration_test.go`: a bad value is the finding `manifest.orchestration`.
- `docs/operators/settings.md`: the row under `## Project manifest`, so `TestSettingsIndexIsComplete` passes.
- `design/system/project-manifest.md`: the key in the `orchestration` block's example and its row in the permissions table.

## Done when

- `flai manifest set orchestration.permissions.confirm_after 48h` writes the key, `flai check` reports `two days` as `manifest.orchestration`, and an unset or `0` value reads as off.
- `flai settings get --json` lists it as a duration with its meaning and risk.
- `flai test` on the changed paths passes, `TestSettingsIndexIsComplete` and `TestSettingsCatalogAgreesWithTheValidation` among them.

## Notes

Layer 2, waits for T-2141 (the key's name and semantics). The serve look task waits for it. S-0498's T-2128 adds the `threads` block to the same files and rows to the same docs; this story waits for S-0498, so the two never meet on a branch.
