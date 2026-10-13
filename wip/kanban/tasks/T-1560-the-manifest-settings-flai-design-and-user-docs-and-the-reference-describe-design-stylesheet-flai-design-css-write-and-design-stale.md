---
id: T-1560
type: task
nature: feature
title: The manifest, settings, flai design, and user docs and the reference describe design.stylesheet, flai design css --write, and design.stale
status: backlog
parent: S-0385
owner: alex
created: 2026-10-13T17:48:59Z
updated: 2026-10-13T17:48:59Z
transitions: []
stream: S-0385
tags: [cli]
touches: [design/system/project-manifest.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-1558, T-1559]
---
# T-1560 The manifest, settings, flai design, and user docs and the reference describe design.stylesheet, flai design css --write, and design.stale

## Work

Add `design.stylesheet` to `design/system/project-manifest.md` and a row for it to `docs/operators/settings.md`, whose test fails on a manifest key without one; describe `flai design css --write` and `design.stale` in `design/system/flai-cli.md` and `docs/users/flai.md`; and regenerate `docs/users/flai-reference.md`. It waits for T-1558 and T-1559, whose behaviour it describes.

## Done when

- The five documents describe it, the settings test passes, and the markdown lint passes on them.

## Notes

- Layer 3.
