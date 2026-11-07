---
id: T-2598
type: task
nature: feature
title: The kluster block takes components, checked by flai check and written by flai manifest set as a components setting
status: done
parent: S-0595
owner: alex
created: 2026-11-07T22:06:20Z
updated: 2026-11-07T22:56:36Z
transitions:
  - to: ready
    at: 2026-11-07T22:09:54Z
    by: agent-S-0595
  - to: in-progress
    at: 2026-11-07T22:10:26Z
    by: agent-S-0595
  - to: done
    at: 2026-11-07T22:24:18Z
    by: agent-S-0595
stream: S-0595
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/kluster_test.go, docs/operators/settings.md, flai/internal/hostapi/settings.go, flai/internal/hostapi/settings_test.go, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 832
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 73
      output: 25775
      cache_read: 12599088
      cache_write: 99409
      cost: 3.6836
    - model: claude-sonnet-5-5
      input: 18
      output: 4846
      cache_read: 261666
      cache_write: 52011
      cost: 0.2047
---
# T-2598 The kluster block takes components, checked by flai check and written by flai manifest set as a components setting

## Work

ADR-0263 §1. Add `Components []KlusterComponent` (`name`, `build` argv) to `manifest.Kluster` (`flai/internal/manifest/manifest.go`), parsed by its `UnmarshalYAML` without failing the load, `components` added to `klusterKeys`. `KlusterErrors` reports, as one sentence each beginning with the key: both `build` and `components` set, a component that is not a block, a missing or repeated `name`, a name that is not a DNS label, a component with no `build` or an empty program, an unknown key in a component. Add the setting `kluster.components` to the catalog in `settings.go` with a new kind `components` (JSON list of `{"name","build"}` objects, parsed strictly, `[]` unsets it), rendered into `system-flow.yaml` as a YAML block list that reads back equal; its `meaning` says what it is. Add its row to `docs/operators/settings.md` beside the other `kluster.*` rows (`TestSettingsIndexIsComplete` reads the manifest's key paths). Layer 1: waits for nothing.

## Done when

- `flai manifest set kluster.components '[{"name":"api","build":["scripts/build.sh","api"]}]'` writes the block and reads back; `[]` removes it; a bad value is refused with the reason.
- `flai check` reports each fault above as `manifest.kluster` on the block's line.
- Tests in `manifest_test.go`, `settings_test.go`, and `check/kluster_test.go` cover them and fail without the change; `flai test` on the changed paths passes.

## Notes
