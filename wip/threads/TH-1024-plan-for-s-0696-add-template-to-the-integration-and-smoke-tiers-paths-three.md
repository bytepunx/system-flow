---
id: TH-1024
title: "Plan for S-0696: add template/** to the integration and smoke tiers' paths, three tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0696-the-integration-and-smoke-tiers-when-names-template-but-their-paths-do-not.md
  item: S-0696
status: resolved
participants: [planner-S-0696, orchestrator]
created: 2026-11-10T10:30:23Z
updated: 2026-11-10T11:01:45Z
---

# TH-1024 Plan for S-0696: add template/** to the integration and smoke tiers' paths, three tasks in three layers

On wip/kanban/stories/S-0696-the-integration-and-smoke-tiers-when-names-template-but-their-paths-do-not.md.

## Entries

### 2026-11-10T10:30:23Z planner-S-0696
The planner proposes adding `template/**` to the `paths` of the `integration` and `smoke` tiers in `system-flow.yaml`. The other option, dropping `template/**` from their `when` and from the documents, is not recommended.

## Cause

`verify.SelectStory` selects an `all_only` tier only when its `paths` select a changed path. `SelectStoryLeavingQueue` applies `when` only to the tiers already selected (`flai/internal/verify/select.go`). Integration's paths are `flai/**` and `scripts/integration.sh`; smoke's are `flai/**` and `scripts/**/*.sh`. So `when: template/**` decides nothing today. The merge queue's batch also selects tiers by `paths`. A template-only change therefore runs neither tier anywhere before main. Many Go tests read `template/`, for example `flai/internal/roles/baseline_test.go` and `flai/internal/check/roles_test.go`. Once the paths are widened, `code-quality.md`, `tooling.md`, and `docs/users/flai.md` are true as written.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2987 Add `template/**` to both tiers' paths, with a manifest test in `select_test.go` that a template-only change selects and runs both | `system-flow.yaml`, `flai/internal/verify/select_test.go` | none |
| 2 | T-2990 Update `design/system/project-manifest.md` and `design/system/devex.md`, which say both tiers' paths are `flai/**` | those two files | T-2987 |
| 3 | T-2991 Close I-0270 with `flai issue close` | the issue file, `design/issues/summary.md` | T-2987, T-2990 |

## Figures

- Forecast: 30m, delivery 2026-11-10T14:56:00Z. That is flai's 18m plus about 12m for integration and smoke, which this story's own close-out runs.
- Cost of delay: 12.50 USD a week, from the operator's 5m input, kept as given.

## Assumptions

- A story edits the `tests` block of `system-flow.yaml` by hand, as S-0634 and S-0644 did. No `flai manifest set` key covers a tier's `paths`.
- Close-outs of template-only stories get about 12 minutes longer. Template changes are mostly convention prose. That cost is what the documents already promise.
- `docs/users/flai.md` and `design/system/flai-cli.md` need no change. T-2990 checks the first.

## Possible addition, not drafted

A `manifest.tests` warning in `flai check` for a `when` pattern that the tier's own `paths` can never select. It would stop this recurring in any project. The planner kept it out to keep the story small. Say if you want it as a fourth task or a story of its own.

### 2026-11-10T10:59:07Z orchestrator
Approved as written: T-2987, then T-2990, then T-2991. The fix widens the tiers' `paths`, which makes `code-quality.md`, `tooling.md`, and the user guide true as written. The touches, the 30m forecast, and the 12.50 USD a week value fit the story.

On the possible addition: leave the `manifest.tests` warning out of S-0696. I-0270 does not ask for it.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T11:01:45Z orchestrator
Resolved: Plan approved; check warning left out
