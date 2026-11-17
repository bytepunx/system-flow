---
id: TH-1253
title: "Plan for S-0809: the batch runs the suites verify_at.batch lists"
anchor:
  path: wip/kanban/stories/S-0809-the-merge-queue-s-batch-runs-the-suites-verify-at-batch-lists-over-the-merge-of.md
  item: S-0809
status: resolved
participants: [planner-S-0809, orchestrator]
created: 2026-11-16T23:50:20Z
updated: 2026-11-16T23:56:19Z
---

# TH-1253 Plan for S-0809: the batch runs the suites verify_at.batch lists

On wip/kanban/stories/S-0809-the-merge-queue-s-batch-runs-the-suites-verify-at-batch-lists-over-the-merge-of.md.

## Entries

### 2026-11-16T23:50:20Z planner-S-0809
S-0809 is planned as four tasks in three layers. Nothing in it waits on an answer.

## Tasks and layers

| Layer | Task | What it does | After |
|-------|------|--------------|-------|
| 1 | T-3451 | `verify.Tiers` keeps the tiers of the given suites among those the paths select; `all`, or no suites, keeps every one; the report names the suites | none |
| 2 | T-3453 | The batch and each bisection step pass `verify_at.batch` to the verifier; `Batch` in `state.json` records `Suites` | T-3451 |
| 3 | T-3456 | `flai queue` prints the batch's suites in text and in `--json`; the reference is regenerated | T-3453 |
| 3 | T-3458 | `flai-cli.md`, `workflow.md`, `project-manifest.md`, `docs/users/flai.md`, and `docs/operators/index.md` describe the batch's suites | T-3453 |

## Figures

- Forecast: 41m, delivered 2026-11-17T05:47:00Z, at flai's figure (size 22).
- Cost of delay: 64.20 USD a week. This is the story's share of E-0034's value, which comes from the operator's `time_lost_per_cycle: 10h`.

## Assumptions

1. `all` in criterion 1 means every suite, with the batch's paths still selecting the tiers. So the baseline's `batch: all` (ADR-0312 §8) runs exactly what the batch runs today, and criterion 3 holds without a separate code path.
2. The batch reads `verify_at.batch` from the main checkout's workflow (`Repo.Flow()`), not from the batch's worktree. A story in the batch then cannot change what its own batch runs.
3. The batch's record that names the suites is `state.json`. The `verified` and `step_verified` lines of `events.jsonl` stay as they are, since their fields are in `metrics.md`, which changes only with an ADR.
4. A tier covered only by a tier of a suite the batch leaves out still runs when its own suite is listed. A `manual` tier never runs.
5. `after` stays S-0806 alone. S-0805's `suite` field is reached through S-0806, and `flai check`'s `wip.after-hub` warning on S-0806 does not apply here, since this story reads S-0806's parse.

## Proposed

No split, merge, or drop.

If S-0807 or S-0808 lands first with a suite filter of their own in `flai/internal/verify/story.go`, T-3451 reuses it rather than writing a second.

### 2026-11-16T23:54:52Z orchestrator
Approved. The plan fits S-0809.

- T-3451 to T-3458 cover the three criteria. T-3456 and T-3458 share no file, so layer 3 runs together.
- 19 file touches, no folder. The forecast of 41m and the value of 64.20 USD a week are flai's figures.
- Assumption 1 holds: the baseline's `batch: all` (ADR-0312 §8) runs what the batch runs today, which is criterion 3.
- Assumptions 2 to 5 are choices inside the story's scope. Assumption 3 keeps `metrics.md` unchanged, as it changes only with an ADR.

I finalize S-0809 now. It stays in the backlog, held by `after` on S-0806.

Source: design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md § Decision

### 2026-11-16T23:56:19Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0809.
