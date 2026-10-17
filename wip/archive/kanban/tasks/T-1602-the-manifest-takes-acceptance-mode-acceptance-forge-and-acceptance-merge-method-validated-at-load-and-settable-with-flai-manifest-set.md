---
id: T-1602
type: task
nature: feature
title: The manifest takes acceptance.mode, acceptance.forge, and acceptance.merge_method, validated at load and settable with flai manifest set
status: done
parent: S-0392
owner: alex
created: 2026-10-16T04:33:20Z
updated: 2026-10-17T07:02:13Z
transitions:
  - to: ready
    at: 2026-10-17T06:42:31Z
    by: agent-S-0392
  - to: in-progress
    at: 2026-10-17T06:42:32Z
    by: agent-S-0392
  - to: done
    at: 2026-10-17T07:02:12Z
    by: agent-S-0392
stream: S-0392
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/manifest/provider.go]
usage:
  source: log
  seconds: 1180
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 47
      output: 21356
      cache_read: 2631056
      cache_write: 89766
      cost: 1.4662
---
# T-1602 The manifest takes acceptance.mode, acceptance.forge, and acceptance.merge_method, validated at load and settable with flai manifest set

## Work

- Add an `Acceptance` block to `Manifest` in `flai/internal/manifest/manifest.go`, with these fields:
  - `mode`: `local` or `pull-request`, `local` when absent.
  - `forge`: `github`, required when the mode is `pull-request`.
  - `merge_method`: `merge`, `squash`, or `rebase`, `merge` when absent.
- Validate the block in `Load`, naming the key and the values allowed.
- Add the three keys to `buildCatalog` in `settings.go`, so that `flai manifest set` takes them and refuses a bad value.
- Test the defaults, each refusal, and a round trip through `EditSettings`. Also test that a manifest with no `acceptance` block reads as before.
- First in its layer: it waits for no task.

## Done when

`flai manifest set acceptance.mode pull-request` writes the key, a bad value is refused with the values allowed, and the manifest tests pass.

## Notes
