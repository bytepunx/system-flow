---
id: T-1455
type: task
nature: improvement
title: A test tier takes fix_command, and verify.Fix runs it over the files the tier selects
status: done
parent: S-0367
owner: alex
created: 2026-10-12T16:52:33Z
updated: 2026-10-12T16:57:25Z
transitions:
  - to: ready
    at: 2026-10-12T16:53:17Z
    by: agent-S-0367
  - to: in-progress
    at: 2026-10-12T16:53:17Z
    by: agent-S-0367
  - to: done
    at: 2026-10-12T16:57:25Z
    by: agent-S-0367
stream: S-0367
tags: []
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/settings.go, flai/internal/manifest/manifest_test.go, flai/internal/verify/fix.go, flai/internal/verify/fix_test.go, flai/internal/verify/manifest.go, flai/internal/verify/verify.go]
usage:
  source: log
  seconds: 248
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 60
      output: 372
      cache_read: 2555281
      cache_write: 94160
      cost: 1.1967
---
# T-1455 A test tier takes fix_command, and verify.Fix runs it over the files the tier selects

## Work

- `manifest.TestTier` gets `FixCommand []string` (`fix_command`), validated as `command` is, placeholders allowed; `settings.go` writes and compares it.
- `verify.Tier` carries it (`FromManifest`).
- `verify.Fix` runs, for each tier with a fix command, that command over the existing files among the given paths the tier's paths select, in the tier's dir, and answers the paths whose content it changed. A command that fails or leaves findings is logged, never an error.
- Waits for nothing: it is the shape every later task uses.

## Done when

- [ ] A manifest with `fix_command` loads, writes back, and a bad one is refused.
- [ ] `verify.Fix` changes and names a file a fix command rewrites, leaves alone a path no tier selects, and survives a failing command, under tests.

## Notes
