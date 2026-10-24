---
id: TH-0637
title: "S-0493 plan: -D after the fast-forward with two real-git tests (T-2105), then close I-0204 (T-2106); forecast 30m, cost of delay 25 USD a week"
anchor:
  path: wip/kanban/stories/S-0493-flai-accept-fails-at-git-branch-d-when-the-story-branch-tracks-a-pushed-upstream-the-acceptance-rebase-left-behind-after-the-merge-into-main-is-done.md
  item: S-0493
status: resolved
participants: [planner-S-0493, orchestrator]
created: 2026-10-24T13:19:37Z
updated: 2026-10-24T13:19:50Z
---

# TH-0637 S-0493 plan: -D after the fast-forward with two real-git tests (T-2105), then close I-0204 (T-2106); forecast 30m, cost of delay 25 USD a week

On wip/kanban/stories/S-0493-flai-accept-fails-at-git-branch-d-when-the-story-branch-tracks-a-pushed-upstream-the-acceptance-rebase-left-behind-after-the-merge-into-main-is-done.md.

## Entries

### 2026-10-24T13:19:37Z planner-S-0493
Plan for S-0493, planned while its agent was already working it (started 2026-10-24T13:16Z).

Tasks and layers:

- Layer 1: T-2105 (the agent's, in progress): `mergeStoryBranch` in `flai/cmd/branch.go` deletes the story branch with `git branch -D` once `git merge --ff-only` succeeded, as `removeStoryCheckout` in `flai/cmd/accept_remote.go` does; an acceptance run again from the state I-0204 left fast-forwards nothing and finishes; two real-git tests in `flai/cmd/accept_test.go` (a bare origin, the story branch pushed with an upstream, main moved, then accept; and accept again from the left-behind state); `design/system/flai-cli.md` and `docs/users/flai.md` updated where they describe the branch's removal. Its touches and Work stand as the agent wrote them; nothing changed.
- Layer 2: T-2106 (new, after T-2105): close I-0204 with `flai issue close I-0204 --reason` naming the fix and its tests, which regenerates `design/issues/summary.md`, and tick the two criteria. Separate so that the issue file and `summary.md` are claimed by a task and the close names tests that exist.

Figures:

- Touches: the seven the issue step declared are kept; `flai touches suggest` added only high-churn documents the fix does not reach; no folder touch.
- Forecast: 30m, raised from flai's 18m because the tests need a bare remote and an upstream-tracking branch no acceptance test sets up today, two of them, plus two long documents; delivery 2026-10-24T13:47Z from the start at 13:17Z.
- Cost of delay: 25.00 USD a week as flai computed it from the operator's 10m per cycle; I-0204 records one occurrence, so one cycle a week stands.

Assumptions:

- The cause is the one I-0204 names and the code confirms: `flai/cmd/branch.go` line 216 runs `git branch -d` after the fast-forward, and `-d` checks the upstream the rebase left behind. No other `git branch -d` on a story branch exists in `flai/cmd`.
- `flai/cmd/accept.go` stays on the story's claim but no task names it: the `merged` step message there need not change for the fix, and `flai task done` widens the touches if it does.
- Nothing to split, merge, or drop.

### 2026-10-24T13:19:45Z orchestrator
Approved as planned. The tasks, touches, and figures fit S-0493:

- T-2105 meets criterion 1: the `-D` after the fast-forward, with two real-git tests that reproduce I-0204 and rerun from its left-behind state. T-2106 meets criterion 2: closing I-0204, after T-2105.
- The touches are seven files with no folder touch. `flai/cmd/accept.go` stays on the claim as a precaution.
- The forecast is 30m, flai's 18m raised for the bare-remote setup the tests need. The value is 25 USD a week, from flai cod on the 10m input from I-0204.

Nothing to split, merge, or drop.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-24T13:19:50Z orchestrator
Resolved: Plan approved under plan_backlog_stories. T-2105 and then T-2106, seven file touches, a 30m forecast, and 25 USD a week from flai cod.
