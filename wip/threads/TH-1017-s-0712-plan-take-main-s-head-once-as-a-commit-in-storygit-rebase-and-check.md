---
id: TH-1017
title: "S-0712 plan: take main's head once as a commit in storygit.Rebase and check undoes_main against it, four tasks in three layers; proposes merging S-0709"
anchor:
  path: wip/kanban/stories/S-0712-flai-task-done-s-sync-reports-undoes-main-for-paths-the-story-never-changed.md
  item: S-0712
status: resolved
participants: [planner-S-0712, orchestrator, alex]
created: 2026-11-10T10:29:25Z
updated: 2026-11-11T09:09:23Z
---

# TH-1017 S-0712 plan: take main's head once as a commit in storygit.Rebase and check undoes_main against it, four tasks in three layers; proposes merging S-0709

On wip/kanban/stories/S-0712-flai-task-done-s-sync-reports-undoes-main-for-paths-the-story-never-changed.md.

## Entries

### 2026-11-10T10:29:25Z planner-S-0712
The plan for S-0712 is written: four tasks in three layers, touches, a forecast, and a cost of delay value.

Cause and fix: `storygit.Rebase` (`flai/internal/storygit/sync.go`) reads the merge base with main and rebases onto `main` by the branch name. Then `FindUndone` (`flai/internal/storygit/undone.go`) runs `git rev-parse main` again and diffs against that newer tip. A commit main gains after the rebase is missing from the branch, so the check counts its lines as removed by the branch. In I-0338 that commit was fe3b4e5, which created S-0706. The fix resolves main to one commit before the merge-base read and uses that commit for the merge base, the rebase, and the check. `flai accept` rebases through the same function and takes the fix with it.

Tasks and layers:

| Layer | Task | Touches | after |
|-------|------|---------|-------|
| 1 | T-2963 take main's head once as a commit in `Rebase` and pass it to `FindUndone`, with a scripted case | `sync.go`, `undone.go`, `undone_test.go` | none |
| 2 | T-2965 real-git test: a runner that commits on main right after the rebase; the sync must answer `undoes_main` empty and record nothing | `sync_test.go` | T-2963 |
| 2 | T-2967 one sentence each in `design/system/workflow.md` (S-0371 paragraph) and `docs/users/flai.md` ("Main's change undone.") | those two files | T-2963 |
| 3 | T-2970 `flai issue close I-0338 --reason`, naming the fix and the test | I-0338 file, `design/issues/summary.md` | T-2965, T-2967 |

Figures:

- Forecast: 45m, raised from flai's 23m for the race reproduction and the integration and smoke tiers at close-out. Delivery is 2026-11-10T16:27:00Z.
- Cost of delay value: 7.50 USD a week, from flai's input of 3m per 168h cycle.

Assumptions:

- `res.Base` keeps the branch's name, so the JSON answer, the help text, and `docs/users/flai-reference.md` do not change.
- ADR-0155 stays as it is. Its check text still fits a base taken as a commit, so no new ADR is needed.
- The scope is the undone check. `fold` and `fix` in the same function also pass main by name. They compare with `main...HEAD`, whose merge base the rebase fixes, so I left them alone.

Proposal, merge S-0709 into S-0712: I-0335 (S-0709, a draft) is the same race. S-0680's sync reported ADR-0123 as undone because main gained two commits during it. T-2963 fixes both. I recommend:

1. Cancel S-0709.
2. Add a criterion to S-0712: "I-0335 is closed with `flai issue close I-0335 --reason` saying what fixed it".
3. Widen T-2970 to close both issues.
4. Add I-0335's 3m to S-0712's cost of delay inputs, which makes the value 15 USD a week.

I changed neither story. Cancelling S-0709 and setting inputs are your decisions.

### 2026-11-10T11:01:57Z orchestrator (recommendation)
**Recommended:** merge the two stories into S-0712, and cancel S-0709.

- Add to S-0712 a third criterion: "I-0335 is closed with `flai issue close I-0335 --reason` saying what fixed it".
- Widen T-2970 to close I-0335 and I-0338 together.
- Have T-2963 also pass the pinned commit to `SettleChangelogs`, which S-0709's plan (TH-1014) found reads main by name after the rebase too.
- Set S-0712's `time_lost_per_cycle` to 6m, I-0335's 3m added to I-0338's 3m. That makes the value about 15 USD a week.

**Why S-0712 and not S-0709:** S-0712's fix resolves main once, before the merge-base read, and uses that one commit for the merge base, the rebase, and the check. S-0709's fix pins it only from the rebase on. S-0712 also already plans the separate real-git reproduction and the `workflow.md` paragraph.

**Why not both:** both stories change the same lines of `storygit.Rebase` for the same race. Two stories would rebuild one fix and conflict on `sync.go`.

This changes scope and an input on S-0712, so it is yours. The same recommendation is on TH-1014.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-11T09:03:11Z alex
Confirmed the recommendation of 2026-11-10T11:01:57Z orchestrator.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-11T09:05:03Z orchestrator
I set S-0712's `time_lost_per_cycle` to 6m, as you confirmed. Its value is left for flai and its planner to work out.

Three parts of the confirmed merge are not mine to do, and are still open:

1. **Cancel S-0709.** No permission of mine covers cancelling a story. TH-1014 carries the same recommendation for you to confirm.
2. **Add S-0712's third criterion**, "I-0335 is closed with `flai issue close I-0335 --reason` saying what fixed it". You finalized S-0712 at 09:03:17Z, so its words are yours.
3. **Widen T-2970 to close I-0335 too, and have T-2963 pass the pinned commit to `SettleChangelogs`.** These are edits to S-0712's tasks. `flai plan --candidates` does not list S-0712, so I cannot start its planner for them.

What you need to do: cancel S-0709, and make edits 2 and 3, or tell S-0712's agent to make them when it starts the story.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-11T09:05:37Z orchestrator
Correction to my last entry: planner-S-0712 was already running. It replanned S-0712 for your confirmation at 09:04–09:05Z, so edits 2 and 3 are done:

- The third criterion, which closes I-0335, is in.
- T-2963 also passes the pinned commit to `settleChangelogs`.
- T-2965 reproduces both issues, and T-2970 closes I-0338 and I-0335.
- The forecast is 1h.

The input is now 6m. The value still reads 7.50 USD a week, because it was worked out from the old 3m. `flai cod` gives 15 USD a week on the next plan.

One thing is still yours: cancel S-0709. TH-1014 holds the same recommendation.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-11T09:06:29Z planner-S-0712
S-0712 is replanned for the confirmed merge with S-0709. Four tasks in three layers, as before, now covering I-0335 too.

## What changed

- S-0712 has a second Goal paragraph naming I-0335, and a third criterion: "I-0335 is closed with `flai issue close I-0335 --reason` saying what fixed it".
- The story's touches gain the I-0335 issue file.
- T-2963 also passes the pinned commit to `SettleChangelogs`, through `settleChangelogs` in `sync.go`, which passed `res.Base`, the name.
- T-2965 has two real-git cases: an undone case with two commits main gains after the rebase (I-0338, I-0335), and a changelog case.
- T-2967 also mentions the changelog settle and links I-0335.
- T-2970 closes I-0338 and I-0335, and touches both issue files.

## Tasks and layers

| Layer | Task | Touches | after |
|-------|------|---------|-------|
| 1 | T-2963 pin main's commit in `Rebase`; pass it to the merge base, the rebase, `FindUndone`, and `SettleChangelogs`; scripted case | `sync.go`, `undone.go`, `undone_test.go` | none |
| 2 | T-2965 real-git undone and changelog cases with a runner that commits on main after the rebase | `sync_test.go` | T-2963 |
| 2 | T-2967 one sentence each in `design/system/workflow.md` and `docs/users/flai.md` | those two files | T-2963 |
| 3 | T-2970 close I-0338 and I-0335 | both issue files, `design/issues/summary.md` | T-2965, T-2967 |

## Figures

- Forecast: 1h, raised from flai's 23m and from the earlier 45m, for the second test case, the changelog pin, and the integration and smoke tiers. Delivery 2026-11-11T22:56:00Z.
- Cost of delay value: 15 USD a week, from the 6m input the orchestrator set.

## Assumptions

- `SettleChangelogs` keeps its signature. The merge queue's `build.settleChangelogs` (`flai/internal/mergequeue/batch.go`) already passes the batch head it rebased onto, so it is left alone.
- `res.Base` keeps the branch's name, so the JSON answer, the help text, and `docs/users/flai-reference.md` do not change.
- ADR-0155 stays as it is.

## For the operator

S-0709 is still open, with T-2959 and T-2960. Cancelling it is yours or the orchestrator's. This plan did not touch it.

### 2026-11-11T09:08:00Z planner-S-0712
S-0712's plan holds after the orchestrator set the 6m input. The tasks and layers are as in my last entry. Only the forecast's numbers changed.

- Cost of delay value: 15 USD a week, as `flai cod` gives it from 6m per 168h cycle at 150 USD an hour. It stands.
- Forecast: 1h stands. flai now gives 27m at size 12 (3 criteria, 9 touches). The 33m added is for the two race reproductions, the changelog pin, and the integration and smoke tiers. Delivery is 2026-11-11T22:58:00Z.
- Touches: unchanged. `flai touches suggest` adds only documents that change with most commits, and none of them changes here.
- Tasks: T-2963, T-2965, T-2967, and T-2970 were revisited and need no change.

Still yours: cancel S-0709 (T-2959, T-2960). TH-1014 carries the same recommendation.

### 2026-11-11T09:08:58Z orchestrator
Plan approved as replanned. T-2963, then T-2965 and T-2967, then T-2970, in three layers. The touches cover every task's files and both issues. The 1h forecast and the 15 USD a week value fit the story.

S-0709's cancellation is not part of this plan. It stays open on TH-1014, where the recommendation awaits the operator. I am resolving this thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-11T09:09:23Z orchestrator
Resolved: Replanned plan approved; S-0709's cancellation remains with the operator on TH-1014.
