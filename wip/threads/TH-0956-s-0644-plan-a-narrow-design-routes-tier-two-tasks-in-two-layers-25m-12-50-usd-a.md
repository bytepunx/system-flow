---
id: TH-0956
title: "S-0644 plan: a narrow design-routes tier, two tasks in two layers, 25m, 12.50 USD a week"
anchor:
  path: wip/kanban/stories/S-0644-a-change-to-flaiover-dashboard-md-alone-selects-no-vitest-tier-so-design-routes.md
  item: S-0644
status: resolved
participants: [planner-S-0644, orchestrator]
created: 2026-11-09T22:09:02Z
updated: 2026-11-09T22:16:40Z
---

# TH-0956 S-0644 plan: a narrow design-routes tier, two tasks in two layers, 25m, 12.50 USD a week

On wip/kanban/stories/S-0644-a-change-to-flaiover-dashboard-md-alone-selects-no-vitest-tier-so-design-routes.md.

## Entries

### 2026-11-09T22:09:02Z planner-S-0644
Plan for S-0644. Nothing here needs an answer unless you would fix I-0322 differently.

## Proposed fix

Add a tier, `design-routes`, to `system-flow.yaml`'s `tests`. It runs `flaiover/src/lib/server/design-routes.test.ts` alone. `design/system/flaiover-dashboard.md` and `flaiover/src/routes/**` select it. The `flaiover` tier covers it under `--all`. This is the shape S-0592 used for `settings-index` and `flaiover-flai` (I-0304).

Adding the design to the vitest tier's paths would not work. That tier runs `vitest related` on the changed files, which follows the import graph, and the test imports neither the design nor the routes.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2809 Add the `design-routes` tier and a repository-manifest test that reproduces the gap | `system-flow.yaml`, `flai/internal/verify/select_test.go` | none |
| 2 | T-2810 Describe the tier in `devex.md` and `tooling.md`, close I-0322 | `design/system/devex.md`, `design/conventions/tooling.md`, the I-0322 file, `design/issues/summary.md` | T-2809 |

## Figures

- Forecast 25m, raised from `flai forecast`'s 14m. The close-out runs integration and smoke because `select_test.go` is under `flai/`, and the worktree has to install flaiover first.
- Cost of delay 12.50 USD a week, as `flai cod` gives it from the 5m input.

## Assumptions

- The tier also selects `flaiover/src/routes/**`. A deleted route whose row stays in the design is the same failure, and `vitest related` misses it for the same reason. This goes a little past the issue's words. Say so if you want only the design document to select it.
- No `after` on S-0634. That story is in progress and changes the same manifest block and the same three files for I-0308, but this story reads nothing it writes, so a trial merge judges the overlap (ADR-0230).
- Changing the manifest's `tests` by hand follows S-0592's commit 38e021d62.

### 2026-11-09T22:15:41Z orchestrator
Approved: the plan fits S-0644.

- A narrow `design-routes` tier is the fix that works: `vitest related` follows imports, and the test imports neither the design nor the routes. It follows S-0592's `settings-index` precedent.
- Selecting `flaiover/src/routes/**` too is within criterion 1: a deleted route whose row stays is the same cause, missed for the same reason.
- T-2809's repository-manifest test reproduces the gap; T-2810 closes I-0322 after it.
- No `after` on S-0634 is right: it reads nothing S-0634 writes, so the trial merge judges the overlap (ADR-0230).
- The 25m over flai's 14m gives its reason, and the value, 12.50 USD a week, is `flai cod`'s from flai's input.

Source: design/adrs/0230-a-ready-story-whose-claim-overlaps-a-story-in-progress-is-held-only-when-a.md

### 2026-11-09T22:16:40Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: it fits S-0644's criteria, routes selection included.
