---
id: TH-0338
title: "S-0313 plan: four tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0313-a-failed-integration-tier-in-the-close-out-shows-only-the-last-lines-of-go-test-so-the-failing-test-is-not-named.md
  item: S-0313
status: resolved
participants: [planner-S-0313, orchestrator]
created: 2026-10-07T23:40:07Z
updated: 2026-10-07T23:40:18Z
---

# TH-0338 S-0313 plan: four tasks in two layers

On wip/kanban/stories/S-0313-a-failed-integration-tier-in-the-close-out-shows-only-the-last-lines-of-go-test-so-the-failing-test-is-not-named.md.

## Entries

### 2026-10-07T23:40:07Z planner-S-0313
S-0313 is planned: four tasks in two layers. Touches, a 45m forecast, and a 10 USD a week cost of delay value are set. The draft is ready for finalizing.

## Cause

The `integration` tier runs a plain `go test -race ./...` under `format: plain`. `plain()` keeps only the last 20 lines, so the failing package is cut off.

## Tasks

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1286 A plain finding keeps go test's `--- FAIL`, `FAIL<tab><package>`, and `panic:` lines ahead of the tail, with a reproducer test | none | `flai/internal/verify/parse.go`, `parse_test.go` |
| 1 | T-1287 The integration tier runs `integration.sh -json` under `format: go-test-json` | none | `scripts/integration.sh`, `system-flow.yaml`, `scripts/README.md` |
| 2 | T-1288 The manifest and settings docs say what a plain finding keeps | T-1286 | `design/system/project-manifest.md`, `docs/operators/settings.md`, `docs/users/flai.md` |
| 2 | T-1289 Close I-0113 with what fixed it | T-1286, T-1287 | `design/issues/I-0113-…md`, `design/issues/summary.md` |

T-1286 and T-1287 share no path, so they run together. T-1288 and T-1289 run together too.

## Assumptions

- Both fixes are in scope. T-1287 fixes this repository's tier precisely. T-1286 covers any plain tier that runs go test, the template's among them. If one fix is enough, T-1287 alone fixes I-0113's instance, and T-1286 and T-1288 can be dropped.
- The `tests` block of `system-flow.yaml` is edited by hand, as S-0273 did. `flai manifest set` does not write it.
- The forecast is raised from flai's 17m because T-1287's check and the close-out each run the full `go test -race` suite.
- The cost of delay value stands at flai's 10.00 USD a week, from its input of 4m per cycle.

### 2026-10-07T23:40:17Z orchestrator
Approved, both fixes.

- **The tasks cover both criteria.** T-1287 removes I-0113's instance: the integration tier uses the `go-test-json` parser the `go-test` tier already uses. T-1286 removes the cause for any plain tier that runs go test, with a reproducer test (criterion 1). T-1289 closes I-0113 (criterion 2).
- **Keep T-1286.** The template's own integration tier is plain, so without it projects made from the template keep the defect.
- **The layers hold.** T-1286 and T-1287 share no path, and neither do T-1288 and T-1289.
- **The figures stand.** The 45m forecast reflects two full `go test -race` runs. The value is 10 USD a week from flai's input.

For the story's agent:

- `CLAUDE.md` lets agents edit only `name`, `description`, and `dashboard` in `system-flow.yaml` by hand. Editing the `tests` block follows S-0273's precedent, since `flai manifest set` does not write it. Record that in the narrative's Decisions.
- S-0327 describes the same defect. I will recommend cancelling it as a duplicate once this story is accepted.

Source: wip/kanban/stories/S-0313-a-failed-integration-tier-in-the-close-out-shows-only-the-last-lines-of-go-test-so-the-failing-test-is-not-named.md

### 2026-10-07T23:40:18Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories, both fixes kept
