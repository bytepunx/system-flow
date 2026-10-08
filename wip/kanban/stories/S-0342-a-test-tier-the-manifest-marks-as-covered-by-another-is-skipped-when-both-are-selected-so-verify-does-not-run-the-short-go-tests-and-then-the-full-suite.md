---
id: S-0342
type: story
nature: improvement
title: A test tier the manifest marks as covered by another is skipped when both are selected, so verify does not run the short Go tests and then the full suite
status: ready
owner: alex
created: 2026-10-08T07:59:14Z
updated: 2026-10-08T08:37:06Z
transitions:
  - to: ready
    at: 2026-10-08T07:59:30Z
    by: system-flow
tags: [cli]
topics: [testing]
touches: [flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, system-flow.yaml, template/root/system-flow.yaml, design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/verify/verify.go, flai/internal/verify/run.go, flai/internal/verify/run_test.go, flai/internal/verify/text.go, flai/internal/verify/text_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/manifest.go, flai/internal/verify/manifest_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, docs/operators/settings.md, flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts]
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
    - kind: planner
      seconds: 44
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 16
          output: 4821
          cache_read: 422747
          cache_write: 51765
          cost: 0.5952
    - kind: orchestrator
      seconds: 373
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 52
          output: 883
          cache_read: 10160388
          cache_write: 24554
          cost: 2.5095
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h54m
    by: alex
    at: 2026-10-08T08:23:42Z
  value: 285
  by: planner-S-0342
  at: 2026-10-08T08:24:19Z
forecast:
  duration: 41m
  delivery: 2026-10-08T09:21:00Z
  basis: "Its own forecast of 41m; 1st in the pull order with an in-progress limit of 5, behind S-0232, S-0321, S-0322 and S-0340."
  by: flai
  at: 2026-10-08T08:37:06Z
---
# S-0342 A test tier the manifest marks as covered by another is skipped when both are selected, so verify does not run the short Go tests and then the full suite

## Goal

When a story's diff touches `flai/`, `flai verify` selects both the `go-test` tier, `go test -race -short` on the changed packages, and the `integration` tier, `go test -race -count=1 ./...` over everything. The second contains the first: every short test of every package runs again inside it. In S-0287's verify record go-test took 1 min 7 s and integration 3 min 6 s; the first minute bought nothing. The same holds for `flai test --all`.

The manifest can say which tier covers which: a tier that names `covers:` another is selected in its place when both would run, and verify reports the covered tier as skipped with the reason. Run alone, as `flai test` runs it on a few paths, the cheaper tier still runs.

## Acceptance criteria

- [ ] A tier in the manifest's `tests` may carry `covers: [<tier name>...]`; `flai manifest` validation refuses a name that is not a tier or a cycle, and `design/system/project-manifest.md` documents the key.
- [ ] When the tiers verify selects, or `flai test --all` runs, include a tier and one that covers it, the covered tier is not run and the result lists it as `skipped: covered by <tier>`, in the text and in `--json`; when only the covered tier is selected it runs as before.
- [ ] This repository's `system-flow.yaml` marks `integration` as covering `go-test`, and the template's manifest does the same for its `test` and `integration` tiers where they overlap, so a story's close-out that touches `flai/` runs the Go tests once.
- [ ] Tests cover the selection with and without the covering tier, the validation, and the reported skip; `docs/users/flai.md` and `design/system/flai-cli.md` describe it.

## Tasks

- T-1345 A test tier may carry covers, and the manifest's validation refuses an unknown name, the tier itself, and a cycle
- T-1349 A selected tier that another selected tier covers is skipped, not run, with state skipped and the covering tier named
- T-1353 This repository's manifest marks integration as covering go-test, and the template's manifest shows the key
- T-1355 flai test and flai verify list a covered tier as skipped: covered by its tier, in the text and in --json
- T-1358 The review page shows a covered tier as skipped, covered by its tier
- T-1359 The manifest, flai test, and flai verify documents describe covers and the skipped tier

## Notes

S-0287's record on 2026-10-08: go-test 1 min 7 s, integration 3 min 6 s, smoke 11 min 25 s. This story takes the first minute off every close-out that touches `flai/`; the smoke minutes are the subject of the local release server story filed with it.

### Planning

Planned by planner-S-0342 on 2026-10-08, and revisited the same day after the operator set the cost of delay input. Every touch is a file; no folder touch is kept.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/verify/select.go`, `select_test.go` | declared | Where `Select` and `SelectStory` choose the tiers; the cover step goes here (T-1349) |
| `flai/internal/manifest/manifest.go`, `manifest_test.go` | declared | `TestTier.Covers` and its validation in `TestErrors` (T-1345) |
| `system-flow.yaml` | declared | `integration` covers `go-test` (T-1353) |
| `template/root/system-flow.yaml` | declared | Kept as declared; the file does not exist. The template's manifest is `template/root/system-flow.yaml.tmpl` |
| `design/system/project-manifest.md`, `design/system/flai-cli.md`, `docs/users/flai.md` | declared | The key, the skip, and its report (T-1359) |
| `template/root/system-flow.yaml.tmpl` | layout | The template's real manifest (T-1353) |
| `template/CHANGELOG.md` | co-change (10%) | Every change under `template/root/` takes an entry (T-1353) |
| `flai/internal/manifest/settings.go`, `settings_test.go` | layout | `flai manifest set tests=` renders and compares tiers field by field (`testsLines`, `sameTiers`) (T-1345) |
| `docs/operators/settings.md` | co-change (19%) and design | `cmd/settings_doc_test.go` fails when a manifest key has no row, so `tests[].covers` lands with T-1345 |
| `flai/internal/verify/verify.go`, `manifest.go`, `manifest_test.go`, `run.go`, `run_test.go` | layout | `Tier.Covers`, `FromManifest`, the `skipped` state and `covered_by`, and `Run` not starting a covered tier (T-1349) |
| `flai/internal/verify/text.go`, `text_test.go`, `story.go`, `story_test.go` | layout | `flai test`'s text and verify's steps carry the skip (T-1355) |
| `flai/cmd/verify.go`, `verify_test.go` | layout | `verifyText` prints verify's steps (T-1355) |
| `flaiover/src/lib/review.ts`, `review.test.ts` | layout | The review page types a step's state as `passed`, `failed`, or `not-reached` (T-1358) |

`touches suggest` listed the most frequent co-changes (`docs/operators/index.md`, `design/system/flaiover-dashboard.md`, `docs/users/flai-reference.md`, `design/adrs/README.md`, and others); none is taken. `flai-reference.md` is generated from command help and no flag changes. No ADR is planned, since the key extends `tests` as `all_only` did.

Forecast: 41m, as `flai forecast` gives it: 78 s per unit over 51 done improvement stories on claude-opus-5-5 in the large band, times size 31 (4 criteria, 27 touches). It stands: six tasks in three layers, Go, one small flaiover change, and docs, with a close-out that runs integration, smoke, and the flaiover tier. flai plays the delivery out again as the board moves; the story is held meanwhile on overlap with S-0232 (`docs/operators/settings.md`) and S-0321 (`flai/internal/manifest/manifest.go`).

Cost of delay: 285 USD a week, as `flai cod` works it out from the operator's input of 1h54m lost per 168h cycle at 150 USD an hour. It stands as the inputs give it. The input matches the count behind it: about a minute lost on each of the 114 of 136 stories accepted in the 7 days to 2026-10-08 that changed `flai/`. The earlier value, 2.50 USD a week, came from the first input of 1m a cycle and is replaced.
