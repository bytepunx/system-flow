---
id: TH-0402
title: "S-0364 plan: an {importers} placeholder for the vet tier, four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0364-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md
  item: S-0364
status: resolved
participants: [planner-S-0364, orchestrator]
created: 2026-10-12T14:41:46Z
updated: 2026-10-12T14:54:15Z
---

# TH-0402 S-0364 plan: an {importers} placeholder for the vet tier, four tasks in three layers

On wip/kanban/stories/S-0364-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md.

## Entries

### 2026-10-12T14:41:46Z planner-S-0364
S-0364 is planned: four tasks in three layers, 14 file touches, forecast 21m, value 25 USD a week. Each figure's reasoning is under `### Planning` in its Notes.

## Proposed remedy

The goal asks for a proposal before building. I propose a new tier placeholder, `{importers}`: the changed packages and every package of the module that imports one of them, directly or transitively, test files included. This repository's `vet` tier runs on it. `go vet` type-checks test files, so S-0338's break in `flai/cmd/story_start_test.go` would have failed `vet` on `flai test flai/internal/workitem`.

**Recommendation: build `{importers}` as planned.** The alternative is to change `vet` to `go vet ./...` in `system-flow.yaml`, with no code change. It fixes this repository only and vets all 63 packages on every Go change. On this host that took 2.0s warm and 11.8s with an empty build cache. If you prefer it, T-1443 to T-1445 shrink to one manifest-and-docs task, and the reproducing test becomes a manifest test.

## Tasks

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1443 `{importers}` in the manifest's validation | none | `flai/internal/manifest/manifest.go`, `manifest_test.go` |
| 1 | T-1445 Docs for `{importers}` | none | `design/system/project-manifest.md`, `design/system/flai-cli.md`, `docs/operators/settings.md`, `docs/users/flai.md` |
| 2 | T-1444 Fill `{importers}` from the reverse-import graph, with the S-0338 reproduction | T-1443 | `flai/internal/verify/importers.go`, `importers_test.go`, `select.go`, `select_test.go`, `verify.go` |
| 3 | T-1446 `vet` runs on `{importers}`; close I-0128 | T-1444 | `system-flow.yaml`, I-0128's file, `design/issues/summary.md` |

## Assumptions

1. `go-test` and `golangci-lint` keep `{packages}`. Running every importer's short tests on each change would cost most of the integration tier, and vet catches the build break I-0128 records.
2. The reverse graph is read from the checkout with `go/parser`, imports only, without running `go list`. That keeps `Select` pure over the checkout, as `{packages}` is. Build constraints are ignored, so the set is never smaller than the go tool's.
3. The template is not changed: its manifest has no Go tier.
4. No ADR: the placeholder extends `tests` as `covers` (S-0342) did.

## Proposals

No split, merge, or drop.

### 2026-10-12T14:54:14Z orchestrator
Approved: build `{importers}` as planned.

- The four tasks in three layers fit both criteria. T-1443 and T-1445 run together, then T-1444 adds the reverse-import closure with the S-0338 reproduction, then T-1446 switches `vet` to it and closes I-0128.
- All 14 touches are files.
- The figures are flai's: 21m, and 25 USD a week from flai's input of 10m lost per cycle.
- The placeholder extends `tests` the way S-0342's `covers` did, recorded in the living design (project-manifest.md) and the story's Notes, so no ADR is needed. Keeping `go-test` and `golangci-lint` on `{packages}` keeps a leaf change cheap.
- The `go vet ./...` alternative fixes only this repository, so I do not take it.

Source: design/conventions/decisions.md § Rules

### 2026-10-12T14:54:15Z orchestrator
Resolved: Plan approved: {importers} for the vet tier, four tasks in three layers.
