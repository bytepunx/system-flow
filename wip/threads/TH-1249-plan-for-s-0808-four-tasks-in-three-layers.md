---
id: TH-1249
title: "Plan for S-0808: four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0808-the-acceptance-s-gates-the-merge-queue-s-eligibility-queue-tiers-and-verify.md
  item: S-0808
status: resolved
participants: [planner-S-0808, orchestrator]
created: 2026-11-16T23:48:02Z
updated: 2026-11-16T23:51:32Z
---

# TH-1249 Plan for S-0808: four tasks in three layers

On wip/kanban/stories/S-0808-the-acceptance-s-gates-the-merge-queue-s-eligibility-queue-tiers-and-verify.md.

## Entries

### 2026-11-16T23:48:02Z planner-S-0808
S-0808 is planned as four tasks in three layers. Nothing in it waits on an answer.

## Tasks and layers

| Layer | Task | Changes |
|-------|------|---------|
| 1 | T-3433 | `AcceptSuites` in `flai/internal/workitem/transition.go`. `StaleTiers`, `QueueTiersNotRun`, and a new suite-not-run check in `flai/internal/verify/story.go`. |
| 2 | T-3434, after T-3433 | The `queue_tiers` and `verify_behind_main` gates in `flai/internal/preview/accept.go`, the new refusal (criterion 3), and the `flai accept` help. |
| 2 | T-3435, after T-3433 | The merge queue's eligibility in `flai/internal/mergequeue/eligible.go`. It runs beside T-3434, since they share no file. |
| 3 | T-3436, after T-3434 and T-3435 | `design/system/workflow.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, and `docs/operators/index.md`. |

## Figures

- Forecast: 1h15m, raised from flai's 35m. Comparable one-gate stories, S-0683 and S-0792, took about 50m each.
- Delivery: 2026-11-17T06:30Z, after the S-0805 → S-0806 → S-0807 chain.
- Cost of delay: 185.64 USD a week. It is S-0808's share of E-0034's 1500 USD a week.

## Assumptions

- S-0806 exposes the accept transition's `verify:` as a parsed field. S-0805 gives each tier a suite. S-0807's report records each suite's tiers, and S-0807 offers a tiers-of-suites helper, which T-3433 reuses where it exists.
- The criterion-3 refusal is a blocker code of its own in the acceptance preview. On the baseline (`verify: all`) it fires nowhere a blocker did not fire before.
- New tests go in a new `flai/cmd/accept_suites_test.go`. The ADR-0291, ADR-0303, and ADR-0221 tests stay unchanged, as criterion 4 asks.
- The gates read `repo.Flow()` themselves. So `flai/internal/preview/orchestrator.go`, `flai/cmd/verify.go`, and `flai/internal/mergequeue/batch.go` are left out of the touches, though `flai touches suggest` listed them.
- The template has no copy of these documents, so the template does not change.

I would not split, merge, or drop anything.

### 2026-11-16T23:50:46Z orchestrator
Approved. The plan fits S-0808.

- T-3433 to T-3436 cover all four criteria. T-3434 and T-3435 share no file, so layer 2 runs together.
- Every task touch is among the story's touches. Leaving out `orchestrator.go`, `verify.go`, and `batch.go` is reasoned under `### Planning`.
- The forecast of 1h15m, raised from flai's 35m, rests on S-0683 and S-0792. The value of 185.64 USD a week is flai's share of E-0034's input.
- The assumptions only restate ADR-0312 §3 and criterion 4.

S-0808 is finalized. It stays in the backlog, held by `after` on S-0807, until that chain is accepted.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-16T23:51:32Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0808.
