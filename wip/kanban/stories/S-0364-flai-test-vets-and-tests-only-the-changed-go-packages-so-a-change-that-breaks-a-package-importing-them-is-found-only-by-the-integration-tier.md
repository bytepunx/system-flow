---
id: S-0364
type: story
nature: improvement
title: flai test vets and tests only the changed Go packages, so a change that breaks a package importing them is found only by the integration tier
status: backlog
owner: alex
created: 2026-10-12T14:12:31Z
updated: 2026-10-12T14:54:15Z
transitions: []
tags: [cli]
topics: [testing]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/verify/importers.go, flai/internal/verify/importers_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/verify/verify.go, system-flow.yaml, design/system/project-manifest.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md, design/issues/I-0128-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 46
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 14
          output: 87
          cache_read: 1772395
          cache_write: 24321
          cost: 0.453
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-12T14:12:31Z
  value: 25
  by: planner-S-0364
  at: 2026-10-12T14:41:33Z
forecast:
  duration: 21m
  delivery: 2026-10-13T02:32:00Z
  basis: "Median 78 s per unit over 57 done improvement stories on claude-opus-5-5 in the large band, times size 16 (2 criteria, 14 touches); 26th in the pull order with an in-progress limit of 5."
  by: planner-S-0364
  at: 2026-10-12T14:41:33Z
finalized:
  by: orchestrator
  at: 2026-10-12T14:54:15Z
---
# S-0364 flai test vets and tests only the changed Go packages, so a change that breaks a package importing them is found only by the integration tier

## Goal

This story remediates [I-0128](../../../design/issues/I-0128-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md), "flai test vets and tests only the changed Go packages, so a change that breaks a package importing them is found only by the integration tier". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0128 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0128 is closed with `flai issue close I-0128 --reason` saying what fixed it

## Tasks
- T-1443 A test tier's command may hold {importers}, and the manifest's validation accepts it
- T-1444 flai test and flai verify fill {importers} with the changed packages and every package of the module that imports them
- T-1445 The manifest, settings, and flai test documents describe {importers} and which tier uses it
- T-1446 This repository's vet tier runs on {importers}, and I-0128 is closed

## Notes

Cost of delay inputs set by flai from I-0128. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T10:50:48Z, 4.1 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0364 on 2026-10-12.

Proposed remedy, from I-0128's one instance (S-0338): a new tier placeholder, `{importers}`. It stands for the changed packages and every package of the module that imports one of them, directly or transitively, test files included. This repository's `vet` tier runs on it. `go vet` type-checks test files, so S-0338's break in `flai/cmd/story_start_test.go` would have failed `vet` on `flai test flai/internal/workitem`. `go-test` and `golangci-lint` keep `{packages}`: running every importer's tests on each change would cost most of the integration tier.

The alternative was to change `vet` to `go vet ./...` in `system-flow.yaml` alone. It is simpler, but it fixes only this repository and vets every package for every Go change. The placeholder keeps a change to a leaf package cheap and is available to any project's manifest. The plan's thread asks which to build.

Every touch is a file; no folder touch is kept.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/manifest/manifest.go`, `manifest_test.go` | layout and co-change (41%) | `PlaceholderImporters` and `commandProblems` (T-1443) |
| `flai/internal/verify/importers.go`, `importers_test.go` | layout | New: the reverse-import closure from the checkout, and the S-0338 reproduction (T-1444) |
| `flai/internal/verify/select.go`, `select_test.go` | layout and co-change (11%) | `fill` and `allArgv` fill the placeholder (T-1444) |
| `flai/internal/verify/verify.go` | layout | Its comments describe the placeholders (T-1444) |
| `system-flow.yaml` | layout | The `vet` tier's command (T-1446) |
| `design/system/project-manifest.md` | design and co-change (30%) | Documents `tests[].command`'s placeholders (T-1445) |
| `docs/operators/settings.md` | co-change (30%) | The `tests[].command` row (T-1445) |
| `design/system/flai-cli.md`, `docs/users/flai.md` | co-change (22%) | `flai test`'s tier selection (T-1445) |
| I-0128's file, `design/issues/summary.md` | design | `flai issue close` writes both (T-1446) |

`touches suggest` listed other frequent co-changes; none is taken. `flai/cmd/root.go`, `flai/internal/check/check.go`, and `flai/internal/manifest/settings.go` need no change: the check calls the manifest's validation, and the settings render tiers field by field. `template/root/system-flow.yaml.tmpl` and `template/CHANGELOG.md` are not taken because the template's manifest has no Go tier. `docs/users/flai-reference.md` is generated from command help, and no flag changes.

Forecast: 21m, as `flai forecast` gives it: 78 s per unit over 57 done improvement stories on claude-opus-5-5 in the large band, times size 16 (2 criteria, 14 touches). It stands. There are four tasks in three layers, and only T-1444 holds new logic. The close-out runs the integration and smoke tiers as every `flai/` story's does, and the median already includes that. Delivery 2026-10-13T02:32Z plays out the pull order, where the story is 26th.

Cost of delay: 25 USD a week, as `flai cod` works it out from flai's input of 10m lost per 168h cycle at 150 USD an hour. It stands as the input gives it. The input counts I-0128's single instance: 3m24s of integration run plus one more filtered full run to find the package.
