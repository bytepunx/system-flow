---
id: TH-0371
title: "S-0342 plan: six tasks in three layers; the cost of delay input looks one close-out a week"
anchor:
  path: wip/kanban/stories/S-0342-a-test-tier-the-manifest-marks-as-covered-by-another-is-skipped-when-both-are-selected-so-verify-does-not-run-the-short-go-tests-and-then-the-full-suite.md
  item: S-0342
status: resolved
participants: [planner-S-0342, orchestrator]
created: 2026-10-08T08:09:20Z
updated: 2026-10-08T08:24:34Z
---

# TH-0371 S-0342 plan: six tasks in three layers; the cost of delay input looks one close-out a week

On wip/kanban/stories/S-0342-a-test-tier-the-manifest-marks-as-covered-by-another-is-skipped-when-both-are-selected-so-verify-does-not-run-the-short-go-tests-and-then-the-full-suite.md.

## Entries

### 2026-10-08T08:09:20Z planner-S-0342
**Plan.** Six tasks in three layers. Every touch is a file.

| Layer | Task | What | After |
|-------|------|------|-------|
| 1 | T-1345 | `TestTier.Covers`, its validation in `TestErrors`, `flai manifest set` rendering, and the `tests[].covers` row in `docs/operators/settings.md` | none |
| 2 | T-1349 | `verify`: the cover step after `Select`/`SelectStory`, the `skipped` state, `covered_by`, and `Run` not starting a covered tier | T-1345 |
| 2 | T-1353 | `system-flow.yaml`: `integration` covers `go-test`; a comment in the template's manifest; `template/CHANGELOG.md` | T-1345 |
| 3 | T-1355 | `flai test` and `flai verify` text and `--json` show `skipped <name>: covered by <tier>` | T-1349 |
| 3 | T-1358 | The review page labels a `skipped` step | T-1349 |
| 3 | T-1359 | `project-manifest.md`, `flai-cli.md`, `docs/users/flai.md` | T-1349 |

Forecast: 41m. Cost of delay value: 2.50 USD a week, from your inputs.

**Question: the cost of delay input.** Recommended: change `time_lost_per_cycle` from 1m to 1h54m. The value goes from 2.50 to 285 USD a week.

The input says 1m lost per 168h cycle. The story's Notes say the minute is lost on every close-out that touches `flai/`. In the 7 days to 2026-10-08, 114 of the 136 accepted stories changed `flai/`. So the real loss is about 1h54m a week.

Alternative: keep 1m if you meant one minute a week. The input is yours, so I left it unchanged. The plan does not wait on this answer.

**Assumptions.**

1. The template marks nothing as covered. Its `test.sh` and `integration.sh` are placeholders that run different things, so they do not overlap. T-1353 adds a comment showing `covers: [test]` for projects whose integration script runs the behaviour tests too. If you want criterion 3 met by marking the template's `integration` as covering `test`, say so.
2. The declared touch `template/root/system-flow.yaml` does not exist. The template's manifest is `template/root/system-flow.yaml.tmpl`. I kept the declared touch and added the real file.
3. Covering is transitive and needs the covering tier selected. Validation refuses an unknown name, the tier itself, a duplicate, and a cycle. A tier covered only by a tier that was not selected runs as before. So `flai test` on a few paths still runs `go-test`.
4. A covered tier is `skipped`, never `not-reached`, even when an earlier step failed.
5. Fail-fast is lost on these runs. A short-test failure now shows up after integration's ~3 minutes, not go-test's ~1 minute. I accepted this because a run that passes is the common case.
6. T-1358, the dashboard change, goes beyond the criteria. It keeps the review page's reading of the verify record correct. It can be dropped.

**Proposals.**

- Add (not drafted): in this repository, the `flaiover` tier runs `pnpm lint` and `pnpm test:unit`, so it could cover `flaiover-lint` and `vitest`. I would add that to T-1353 if you agree. The criteria name only `integration` and `go-test`.
- Overlap with S-0341 (ready): both stories change `flai/internal/verify/story.go`, `story_test.go`, `run.go`, `run_test.go`, `flai/cmd/verify.go`, `verify_test.go`, `docs/users/flai.md`, and `design/system/flai-cli.md`. Whichever starts second is held until the first reaches review. No split is proposed.
- Nothing to split, merge, or drop.

### 2026-10-08T08:09:30Z orchestrator (recommendation)
Plan approved: six tasks in three layers, all file touches, and a 41m forecast. Covering applies only when the covering tier is selected (assumption 3), so `flai test` on a few paths keeps its fast `go-test`. Keep T-1358: without it, the review page would misread the new `skipped` state. The overlap with S-0341 is handled by flai's hold, and S-0342 (41m) is first in the pull order.

Recommendations for the operator:

1. **Cost of delay input:** set `time_lost_per_cycle` to `1h54m` on S-0342, as the planner recommends. The story's Notes count the minute per close-out that touches `flai/`, not per week. I could not set it myself: flai refuses me inputs on a story in ready.
2. **Adding `flaiover` covers `flaiover-lint` and `vitest`: not in this story.** It goes beyond the criteria, which name only `integration` and `go-test`, and `flaiover-lint` went in deliberately with S-0319 to catch formatting before the close-out. If you want it, a small follow-up story can add it once `covers` exists.
3. **The template (assumption 1):** take the planner's comment-only change. The template's placeholder scripts do not overlap.

Source: wip/kanban/stories/S-0342-a-test-tier-the-manifest-marks-as-covered-by-another-is-skipped-when-both-are-selected-so-verify-does-not-run-the-short-go-tests-and-then-the-full-suite.md

### 2026-10-08T08:24:28Z planner-S-0342
The cost of delay is settled. You set `time_lost_per_cycle` to `1h54m`. `flai cod` now gives 285 USD a week, and I wrote that value in place of 2.50 USD.

## What changed

- S-0342's cost of delay value: 2.50 → 285 USD a week.
- `### Planning` in S-0342's Notes: the new figure and its reason. The forecast paragraph also names the current holds: S-0232 on `docs/operators/settings.md` and S-0321 on `flai/internal/manifest/manifest.go`.

## What stands

- The forecast: 41m, unchanged.
- The six tasks and their touches. I revisited each one against the story. Nothing to split, merge, or drop.

| Layer | Tasks |
|-------|-------|
| 1 | T-1345 |
| 2 | T-1349, T-1353 |
| 3 | T-1355, T-1358, T-1359 |

## Your call

The orchestrator's recommendations 2 and 3 still await your confirmation:

- 2: no flaiover `covers` in this story.
- 3: the template gets a comment only.

The plan needs neither to go ahead.

### 2026-10-08T08:24:34Z orchestrator
Resolved: Settled: plan approved (6 tasks, 41m). The operator set time_lost_per_cycle 1h54m, and the planner recorded the value of 285 USD/week. The plan stands as drafted, with no flaiover covers added and the comment-only template change.
