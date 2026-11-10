---
id: T-2879
type: task
nature: feature
title: The kluster block sets how many story namespaces may live at once
status: backlog
parent: S-0660
owner: alex
created: 2026-11-10T08:27:08Z
updated: 2026-11-10T08:27:08Z
transitions: []
stream: S-0660
tags: [cli]
touches: [flai/internal/manifest/manifest_kluster.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, docs/operators/settings.md, design/system/project-manifest.md]
---
# T-2879 The kluster block sets how many story namespaces may live at once

## Work

Add the limit on story namespaces to the manifest's `kluster:` block, under the key and with the default S-0659's ADR sets. Parse and validate it in `flai/internal/manifest/manifest_kluster.go`, so that `flai check` reports `manifest.kluster` for a value that is not a positive whole number. Let `flai manifest set` change it through `flai/internal/manifest/settings.go`.

Give it its row in `docs/operators/settings.md` beside `kluster.namespace`, and its line in `design/system/project-manifest.md`'s `kluster:` block, in this task, since the settings test fails on a manifest key with no row.

It waits for nothing: the first layer.

## Done when

- A block without the key reads the default; a block with it reads its value; `flai manifest set` sets and unsets it; a bad value is a `manifest.kluster` finding.
- `docs/operators/settings.md` and `design/system/project-manifest.md` name the key, its default, and what a request over it gets.
- `go test ./internal/manifest/...` passes.

## Notes
