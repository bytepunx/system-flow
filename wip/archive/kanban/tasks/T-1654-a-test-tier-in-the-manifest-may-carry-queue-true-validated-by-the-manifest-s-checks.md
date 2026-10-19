---
id: T-1654
type: task
nature: improvement
title: "A test tier in the manifest may carry queue: true, validated by the manifest's checks"
status: done
parent: S-0405
owner: alex
created: 2026-10-16T05:01:38Z
updated: 2026-10-19T20:42:58Z
transitions:
  - to: ready
    at: 2026-10-19T20:39:51Z
    by: system-flow
  - to: in-progress
    at: 2026-10-19T20:39:51Z
    by: system-flow
  - to: done
    at: 2026-10-19T20:42:58Z
    by: system-flow
stream: S-0405
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/verify/manifest.go, flai/internal/verify/verify.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/verify/manifest_test.go]
usage:
  source: log
  seconds: 187
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 17
      output: 5684
      cache_read: 1086415
      cache_write: 44955
      cost: 0.5899
---
# T-1654 A test tier in the manifest may carry queue: true, validated by the manifest's checks

## Work

- Add `queue`, a boolean, to `TestTier` in `manifest.go`, and carry it to `verify.Tier` in `verify/manifest.go`.
- Check it in the manifest's checks: a value that is not a boolean is an error on its line.
- Test decoding, the default, and the error.
- First in its layer: it waits for no task.

## Done when

The key decodes and reaches the verify tiers, a bad value is reported, and the tests pass.

## Notes
