---
id: TH-1250
title: "Plan for S-0811: four tasks in three layers, the publish runs verify_at.publish's suites before it tags"
anchor:
  path: wip/kanban/stories/S-0811-the-release-s-publish-runs-the-suites-verify-at-publish-lists-at-the-main.md
  item: S-0811
status: resolved
participants: [planner-S-0811, orchestrator, agent-S-0811, flai serve]
created: 2026-11-16T23:48:42Z
updated: 2026-11-20T08:18:35Z
---

# TH-1250 Plan for S-0811: four tasks in three layers, the publish runs verify_at.publish's suites before it tags

On wip/kanban/stories/S-0811-the-release-s-publish-runs-the-suites-verify-at-publish-lists-at-the-main.md.

## Entries

### 2026-11-16T23:48:42Z planner-S-0811
S-0811 is planned: 14 file touches, a forecast of 34m, a cost of delay of 90.11 USD a week, and four tasks in three layers. Nothing here waits on an answer. The assumptions below are open to change.

## Tasks and layers

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-3441 `flai release --pending` runs `verify_at.publish`'s suites at main's head before it applies or tags, and exits with a new code naming the suite and tier | none | `flai/cmd/release.go`, `flai/cmd/release_publish_suites.go`, `flai/cmd/release_publish_suites_test.go`, `docs/users/flai-reference.md` |
| 2 | T-3442 the dashboard's Publish and `release_publish` map the new exit (`publishExits`) | T-3441 | `flai/internal/hostapi/writes_outcome.go`, `flai/internal/hostapi/writes_test.go` |
| 2 | T-3443 the publish formula's publish step lists the new exit in `stop`, and its story records the suite, tier, and log | T-3441 | both copies of `publish.md`, `flai/cmd/publish_formula_suites_test.go`, `template/CHANGELOG.md` |
| 3 | T-3445 the release design, the workflow design, the board runbook, and the user guide | T-3441, T-3443 | `design/system/pushing-from-the-board.md`, `design/system/workflow.md`, `docs/operators/runbooks/change-the-board.md`, `docs/users/flai.md` |

T-3442 and T-3443 share no path, so they run together.

## Assumptions

1. **Tier selection.** At publish, a suite runs every one of its tiers, as `flai test --suite <name> --all` does, `manual` tiers excepted. Paths do not narrow it, since a publish has no story's change set.
2. **Order.** The suites run after the `push` gate and the remote check, and only when something is pending. With nothing pending they do not run. `--dry-run` names the suites and runs none.
3. **The exit.** The new code is the next unused one after 2, 3, and 4: 5 unless another is taken. Its message names the suite, the tier, and the log under `.flai-cache/`.
4. **The formula stops cleanly.** "Stops on that code" is read as ADR-0271's `stop`. The instance records the failed suite in its Notes and reaches review without a thread, as the step's exit 2 does today. The alternative is the failure path: block the task and ask the operator on a thread. That makes a broken main harder to miss. Say if you want it, and T-3443 changes.
5. **The release design** is the opening current-state paragraph of `design/system/pushing-from-the-board.md`. **The runbook** is `docs/operators/runbooks/change-the-board.md`, where the workflow file is changed.
6. **After.** Only S-0806 stays in `after`, since S-0811 reads its `verify_at` parse. S-0805's suite selection is reached through S-0806. `flai check` advises `wip.after-hub` on S-0806, which five open stories name. S-0811 does read its output, so the entry stays.
7. **Figures.** The forecast and the cost of delay are flai's as computed. The reasons are under `### Planning` in the story's Notes.

## Proposed changes

None. No task to split, merge, or drop.

### 2026-11-16T23:51:39Z orchestrator (recommendation)
**Recommended:** take the alternative in assumption 4. When a suite fails at publish, the formula's publish step fails: the agent blocks the task and asks on a thread. Do not list the new exit in the step's `stop`. T-3443 changes to that. Everything else in the plan stands, and I approve it.

**Why:**

- ADR-0271 keeps `stop` for exits that are the formula's answer, not a failure. Exit 2 from `evaluate` (policy not met) is an answer. A suite that fails at main's head means main is broken, which is a failure.
- If the step stops cleanly, the instance reaches review and I accept it with no one asked. Nothing then tells anyone that main is broken, except one line in the Notes.
- The failure path still meets criterion 2. The step stops on the code, and its blocked task's Notes and its thread name the suite and the tier.

**Why I am asking you:** ADR-0312 §6 says the publish refuses when a suite fails. It does not say what the formula does then, so this is a decision no one has recorded yet.

**Alternative:** keep the planner's reading, a clean stop with the suite recorded in the Notes. Say so, and T-3443 stays as written.

The rest of the plan fits S-0811: four tasks over all four criteria, 14 file touches, and flai's own 34m forecast and value of 90.11 USD a week. S-0811 is finalized and held by `after` on S-0806, so this answer delays nothing.

Source: design/adrs/0271-a-formula-step-may-name-the-exits-that-stop-its-instance-cleanly-flai-formula.md § Decision

### 2026-11-18T08:16:33Z agent-S-0811
S-0811 is started and follows the orchestrator's recommendation: a suite that fails at publish fails the formula's publish step. The new exit is not added to its `stop`, so the instance blocks the task and asks, with the suite, tier, and log in the task's Notes. T-3443 is rewritten to say so.

Two other changes to the plan:

- T-3443 changes only the template's copy of the formula. This repository has had no copy since main's 4052fc603, and it is not restored (D-0662). ADR-0316 and S-0837 retire both later.
- T-3442 covers the publish path that ADR-0316 keeps: `release_publish` and the dashboard's Publish answer the refusal with flai's message, which names the suite and the tier.

If the operator would rather have the clean stop, say so here. T-3443 is small and changes easily.

Source: design/digest/D-0662-since-main-s-4052fc603-this-repository-has-no-design-workflows-formulas-publish.md

### 2026-11-18T08:17:37Z orchestrator (recommendation)
**Recommended:** confirm the failure path that agent-S-0811 is building. When a suite fails at publish, the formula's publish step fails. The new exit is not in its `stop`. The instance blocks the task and asks on a thread, with the suite, the tier, and the log in the task's Notes.

**Why I am posting this again:** agent-S-0811's entry cleared my recommendation of 2026-11-16T23:51:39Z. Posting it again lets you confirm it or answer otherwise. The recommendation itself has not changed.

**Why:**

- ADR-0271 §1 and §3 keep `stop` for an exit that is the formula's answer, after which the agent "asks nothing". A suite that fails at main's head means main is broken. That is a failure, and someone should be asked about it.
- The failure path still meets criterion 2. The step stops on the code, and the blocked task's Notes name the suite.

**The agent's two other changes fit the story:**

- T-3443 changes only the template's copy of the formula. This repository has had no copy since 4052fc603 (D-0662).
- T-3442 covers `release_publish` and the dashboard's Publish. Proposed ADR-0316 keeps that route.

**Why I am asking you:** ADR-0312 §6 says the publish refuses when a suite fails. It does not say what the formula does then. No one has recorded that decision yet.

**Alternative:** a clean stop, with the suite recorded in the story's Notes. Say so here, and T-3443 adds the exit to `stop`.

Source: design/adrs/0271-a-formula-step-may-name-the-exits-that-stop-its-instance-cleanly-flai-formula.md § Decision

### 2026-11-20T08:17:39Z flai serve
Confirmed the recommendation of 2026-11-18T08:17:37Z orchestrator after 48h with no answer from the operator (orchestration.permissions.confirm_after).

Source: design/adrs/0271-a-formula-step-may-name-the-exits-that-stop-its-instance-cleanly-flai-formula.md § Decision

### 2026-11-20T08:18:35Z orchestrator
Resolved: The plan is approved and its open decision is settled: flai serve confirmed the failure path after 48h (confirm_after). A suite that fails at publish fails the formula's publish step, and the instance blocks and asks. S-0811 was built that way.
