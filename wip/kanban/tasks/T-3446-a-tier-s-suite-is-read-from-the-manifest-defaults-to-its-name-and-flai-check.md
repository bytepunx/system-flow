---
id: T-3446
type: task
nature: feature
title: A tier's suite is read from the manifest, defaults to its name, and flai check refuses one that is not kebab-case
status: done
parent: S-0805
owner: alex
created: 2026-11-16T23:48:35Z
updated: 2026-11-17T00:13:04Z
transitions:
  - to: ready
    at: 2026-11-16T23:54:06Z
    by: agent-S-0805
  - to: in-progress
    at: 2026-11-16T23:54:06Z
    by: agent-S-0805
  - to: done
    at: 2026-11-17T00:13:04Z
    by: agent-S-0805
stream: S-0805
tags: [flai]
touches: [flai/internal/manifest/manifest_tests.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/tests_test.go, docs/operators/settings.md]
after: [T-3444]
usage:
  source: log
  seconds: 1138
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 6653
      cache_read: 1916803
      cache_write: 70652
      cost: 0.955
---
# T-3446 A tier's suite is read from the manifest, defaults to its name, and flai check refuses one that is not kebab-case

## Work

Add `Suite string` (`yaml:"suite,omitempty" json:"suite,omitempty"`) to `TestTier` in `flai/internal/manifest/manifest_tests.go`, kept raw as `Format` is, with `SuiteOrDefault()` beside `FormatOrDefault` giving the tier's name when it is empty. In `TestErrors`, refuse a `suite` that is not kebab-case, with the pattern the role and formula files use (`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`); `flai check` reports it as `manifest.tests` through `check_layout.go` with no new rule. Make `flai manifest set tests=` in `settings.go` parse, render (`testsLines`), and compare (`sameTiers`) the key, and give it its Meaning in the settings catalog. Add the `tests[].suite` row to `docs/operators/settings.md`, which `TestSettingsIndexHasEveryKey` requires once the key exists. `manifest.Load` is not strict, so an older flai ignores the key. Waits for T-3444, whose ADR fixes the key and its default.

## Done when

- [ ] A tier with `suite: static` reads back `static`, one with none reads back its own name from `SuiteOrDefault()`, and the YAML and JSON round trips keep it.
- [ ] `TestErrors`, and so `flai check`, refuses `suite: Static` and `suite: a_b` with `manifest.tests`, naming the tier, and the check test covers it.
- [ ] `flai manifest set tests=` keeps a tier's `suite`, and the settings index test passes with the new row.
- [ ] `flai test` on the changed paths passes.

## Notes
