---
id: T-2596
type: task
nature: improvement
title: A change to the settings index selects a tier that runs TestSettingsIndex
status: done
parent: S-0592
owner: alex
created: 2026-11-07T22:03:04Z
updated: 2026-11-07T22:23:41Z
transitions:
  - to: ready
    at: 2026-11-07T22:19:44Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T22:19:45Z
    by: system-flow
  - to: done
    at: 2026-11-07T22:23:41Z
    by: agent-S-0592
stream: S-0592
tags: [flai, tests]
touches: [system-flow.yaml, flai/internal/verify/select_test.go]
after: [T-2594, T-2595]
usage:
  source: log
  seconds: 236
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 15
      output: 2777
      cache_read: 1674451
      cache_write: 33977
      cost: 0.6511
  strategic:
    - kind: orchestrator
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 19330
          cache_write: 573
          cost: 0.005
---
# T-2596 A change to the settings index selects a tier that runs TestSettingsIndex

## Work

A change to `docs/operators/settings.md` alone selects only the `markdown` tier, so no Go test checks the index against the configuration and the manifest until the close-out's integration tier (I-0304's second instance).

Add a tier `settings-index` to `system-flow.yaml` after `go-test`: `dir: flai`, command `[go, test, "-json", "-short", "-count=1", "-timeout", "5m", "-run", "^TestSettingsIndex", "./cmd"]`, `format: go-test-json`, paths `docs/operators/settings.md`. Add it to the `integration` tier's `covers`, which runs the test whole. Add a repository-manifest test to `flai/internal/verify/select_test.go`: a change to `docs/operators/settings.md` alone selects `settings-index`.

It waits for T-2594, since under `-short` the tier runs only the checks that task makes run there, and for T-2595, which changes the same two files.

## Done when

- The new test fails without the tier and passes with it.
- `flai test docs/operators/settings.md` runs `settings-index` and fails when a manifest key's row is removed from the page, checked by hand and restored.
- `flai check --strict` reports nothing on `system-flow.yaml`.

## Notes
