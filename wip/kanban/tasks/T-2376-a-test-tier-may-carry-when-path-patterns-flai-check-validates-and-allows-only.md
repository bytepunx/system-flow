---
id: T-2376
type: task
nature: feature
title: "A test tier may carry when, path patterns flai check validates and allows only beside queue: true, read into verify's Tier and described in the settings index"
status: backlog
parent: S-0547
owner: alex
created: 2026-10-31T07:39:35Z
updated: 2026-10-31T07:39:35Z
transitions: []
stream: S-0547
tags: [flai]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/verify/verify.go, flai/internal/verify/manifest.go, flai/internal/verify/manifest_test.go, flai/internal/check/tests_test.go, docs/operators/settings.md]
---
# T-2376 A test tier may carry when, path patterns flai check validates and allows only beside queue: true, read into verify's Tier and described in the settings index

## Work

Criterion 2 needs the key before selection can use it. Add `When []string` (`yaml:"when,omitempty"`) to the manifest's test tier in `flai/internal/manifest/manifest.go`, beside `Queue`, and carry it into `verify.Tier` (`flai/internal/verify/verify.go`) through `flai/internal/verify/manifest.go`.

In `Manifest.TestErrors`, which `flai check` reports as `manifest.tests`, refuse:

- `when` on a tier without `queue: true`, saying that `when` only chooses when a queue tier also runs in the close-out;
- an empty pattern, or one the `paths` pattern language does not parse, by index, field, and tier name as the other tier errors read.

Give `when` its place wherever `settings.go` lists, compares, or writes a tier's fields (around `field("queue", ...)` and the tier equality), so that `flai manifest set tests=<JSON>` keeps it. Add the `tests[].when` row to `docs/operators/settings.md` after `tests[].queue`: the settings test fails on a manifest key without a row.

Waits for nothing: it touches no file another task of the first layer touches.

## Done when

- A manifest with `when` on a `queue: true` tier loads, and `verify.Tier.When` holds its patterns.
- `flai check` reports `manifest.tests` for `when` on a tier without `queue`, and for a bad pattern, each with its index and tier name, in tests in `flai/internal/manifest/manifest_test.go` and `flai/internal/check/tests_test.go`.
- `flai manifest set tests=<JSON>` round-trips a tier with `when`.
- `docs/operators/settings.md` has the `tests[].when` row, and `flai test` on the changed paths passes.

## Notes

Drafted by planner-S-0547. That `when` is refused without `queue: true` is an assumption on the plan's thread.
