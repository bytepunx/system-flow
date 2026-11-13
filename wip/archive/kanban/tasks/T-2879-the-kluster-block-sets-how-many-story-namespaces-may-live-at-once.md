---
id: T-2879
type: task
nature: feature
title: The kluster block sets how many story namespaces may live at once
status: done
parent: S-0660
owner: alex
created: 2026-11-10T08:27:08Z
updated: 2026-11-13T07:55:17Z
transitions:
  - to: ready
    at: 2026-11-13T07:33:53Z
    by: agent-S-0660
  - to: in-progress
    at: 2026-11-13T07:33:53Z
    by: agent-S-0660
  - to: done
    at: 2026-11-13T07:55:17Z
    by: agent-S-0660
stream: S-0660
tags: [cli]
touches: [flai/internal/manifest/manifest_kluster.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, docs/operators/settings.md, design/system/project-manifest.md, flai/internal/hostapi/kluster_test.go, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 1284
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 48
      output: 23010
      cache_read: 3757557
      cache_write: 117870
      cost: 1.8929
---
# T-2879 The kluster block sets how many story namespaces may live at once

## Work

ADR-0292 §4: add `kluster.agent_namespaces` to the manifest's `kluster:` block, an optional whole number of 0 or more, default 2, the most stories in progress that hold a namespace at once; 0 turns deploys before review off. Parse and validate it in `flai/internal/manifest/manifest_kluster.go`, so that `flai check` reports `manifest.kluster` for a value that is not a whole number of 0 or more, with a getter that gives the default when absent. Let `flai manifest set kluster.agent_namespaces` set and unset it through `flai/internal/manifest/settings.go`, beside the block's other keys, so the dashboard's setup form lists it under `settings`.

Give it its row in `docs/operators/settings.md` beside `kluster.namespace`, and its line and row in `design/system/project-manifest.md`'s `kluster:` block, in this task, since the settings test fails on a manifest key with no row.

It waits for nothing: the first layer.

## Done when

- A block without the key reads 2; a block with it reads its value, 0 included; `flai manifest set` sets and unsets it; a negative or non-integer value is a `manifest.kluster` finding.
- `docs/operators/settings.md` and `design/system/project-manifest.md` name the key, its default, what 0 does, and what a request over it gets.
- `flai test` on the changed paths passes.

## Notes
