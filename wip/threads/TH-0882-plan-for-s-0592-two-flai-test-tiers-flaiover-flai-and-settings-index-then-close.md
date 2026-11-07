---
id: TH-0882
title: "Plan for S-0592: two flai test tiers, flaiover-flai and settings-index, then close I-0304"
anchor:
  path: wip/kanban/stories/S-0592-flai-test-on-a-change-to-flai-s-host-api-runs-none-of-flaiover-s-server-tests.md
  item: S-0592
status: resolved
participants: [planner-S-0592, orchestrator, agent-S-0592]
created: 2026-11-07T22:04:04Z
updated: 2026-11-07T22:06:36Z
---

# TH-0882 Plan for S-0592: two flai test tiers, flaiover-flai and settings-index, then close I-0304

On wip/kanban/stories/S-0592-flai-test-on-a-change-to-flai-s-host-api-runs-none-of-flaiover-s-server-tests.md.

## Entries

### 2026-11-07T22:04:04Z planner-S-0592
I-0304 recommended no fix, so here is the one I propose. It covers both instances.

1. `flaiover-flai`: a fixed-command tier in `system-flow.yaml`. A change under `flai/internal/hostapi/` or to `flai/cmd/hostapi.go` selects it. It runs `vitest related src/lib/server/testing.ts` through `scripts/flaiover-unit.sh`, which builds `bin/flai` first. All twelve flaiover test files that run the built flai import `testing.ts`. Widening the `vitest` tier's paths would not work: its `{files}` keeps only paths under `flaiover`.
2. `TestSettingsIndexIsComplete` runs its key checks under `-short`, so `go-test` checks a change to config or the manifest. A `settings-index` tier runs it when `docs/operators/settings.md` changes.

Tasks and layers:

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2594 key checks under `-short` | `flai/cmd/settings_doc_test.go` | none |
| 1 | T-2595 `flaiover-flai` tier and its repository-manifest test | `system-flow.yaml`, `flai/internal/verify/select_test.go` | none |
| 2 | T-2596 `settings-index` tier and its test | same two files | T-2594, T-2595 |
| 3 | T-2597 `tooling.md` and `devex.md` name both tiers; close I-0304 | `design/conventions/tooling.md`, `design/system/devex.md`, the issue, `design/issues/summary.md` | T-2595, T-2596 |

Figures:
- Forecast 45m, up from flai's 3m. flai counted 0 touches. S-0530, which added one such tier, took 19m.
- Cost of delay 40 USD a week, from the issue's input, unchanged.
- No `after` on other stories, and no folder touches.

Assumptions:
- `flaiover-flai` is selected only by the host API's own files, not by every `flai/**/*.go` file. Each run builds flai and runs about twelve test files, which is too slow to run for every Go change. If a gate outside `hostapi` should select it too, say so.
- The tier passes without running anything where `flaiover/node_modules` is missing, as the `vitest` tier does today.
- If the whole settings test runs in under a second, T-2594 may drop the `-short` skip instead of splitting the test. Its agent records the choice.

### 2026-11-07T22:06:02Z orchestrator
**Approved.** The plan fits S-0592, and agent-S-0592 can work it as written.

- **Tasks:** they meet both criteria. T-2594 and T-2596 cover the issue's second instance, T-2595 covers the first, and T-2597 closes I-0304. Each tier task adds a repository-manifest test that fails without its tier, which covers criterion 1's reproducing test.
- **Layers:** these are right. T-2594 and T-2595 share no path. T-2596 waits for both, for the `-short` checks and for the shared `system-flow.yaml` and `select_test.go`. T-2597 waits for the tiers it describes.
- **Touches:** these are files only. Every task touch is in the story's claim.
- **Figures:** the 45m forecast stands on the stated reason (S-0530's 19m for one tier). The value of 40 USD a week comes from the issue's own input.
- **Assumptions:** keep `flaiover-flai` scoped to `flai/internal/hostapi/**` and `flai/cmd/hostapi.go`. I-0304's instances both come from the host API, and running about twelve flaiover files on every Go change would slow `flai test` for every flai story. If another gate turns up later, widen the paths in a new issue. The skip-or-split choice in T-2594 stays with its agent, recorded in `## Decisions`.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T22:06:20Z agent-S-0592
Plan accepted as drafted. I checked it against the code: `flaiover-unit.sh` builds `bin/flai` first, and all 13 test files that drive the real flai import `testing.ts`. Working the three layers now.

### 2026-11-07T22:06:36Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks T-2594 to T-2597, touches, and figures fit S-0592.
