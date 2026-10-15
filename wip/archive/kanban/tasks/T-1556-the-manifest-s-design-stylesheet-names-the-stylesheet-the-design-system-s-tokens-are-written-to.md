---
id: T-1556
type: task
nature: feature
title: The manifest's design.stylesheet names the stylesheet the design system's tokens are written to
status: done
parent: S-0385
owner: alex
created: 2026-10-13T17:48:42Z
updated: 2026-10-15T02:30:28Z
transitions:
  - to: ready
    at: 2026-10-15T02:23:00Z
    by: agent-S-0385
  - to: in-progress
    at: 2026-10-15T02:23:00Z
    by: agent-S-0385
  - to: done
    at: 2026-10-15T02:30:27Z
    by: agent-S-0385
stream: S-0385
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, docs/operators/settings.md, design/system/project-manifest.md, flai/internal/designsystem/testdata/full.tailwind.css.golden, flai/internal/designsystem/testdata/show-problems.json.golden, flai/internal/designsystem/testdata/show.json.golden, flai/internal/designsystem/testdata/starter.stylesheet.css.golden, flai/internal/designsystem/testdata/themes.tailwind.css.golden]
usage:
  source: log
  seconds: 447
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 29
      output: 9293
      cache_read: 1434034
      cache_write: 53882
      cost: 0.8102
---
# T-1556 The manifest's design.stylesheet names the stylesheet the design system's tokens are written to

## Work

Add `design.stylesheet` to the manifest in `flai/internal/manifest/manifest.go`: a repository path, optional, which `Load` refuses as it refuses a bad `design.tokens` (absolute, leaving the repository, not clean), and a helper that answers the sub-project (`projects` entry) whose path holds it, or none. Tests cover it set, unset, refused, and outside every sub-project.

Add its row to `docs/operators/settings.md`, whose test (`flai/cmd/settings_doc_test.go`) fails on a manifest key without one, and describe it in `design/system/project-manifest.md` beside `design.tokens`, in the same change, so that the key and its documentation land together. It waits for nothing.

## Done when

- `flai test` on the manifest package and `flai/cmd/settings_doc_test.go`'s package passes with the new cases.

## Notes

- Layer 1, beside T-1557. The settings row moved here from T-1560 so that no task leaves the settings test failing.
