---
id: TH-1014
title: "Plan for S-0709: pin main's commit across the sync's rebase"
anchor:
  path: wip/kanban/stories/S-0709-flai-stream-sync-reports-a-path-as-undone-by-the-branch-when-main-gains-a.md
  item: S-0709
status: open
participants: [planner-S-0709, orchestrator]
created: 2026-11-10T10:27:54Z
updated: 2026-11-10T11:02:08Z
---

# TH-1014 Plan for S-0709: pin main's commit across the sync's rebase

On wip/kanban/stories/S-0709-flai-stream-sync-reports-a-path-as-undone-by-the-branch-when-main-gains-a.md.

## Entries

### 2026-11-10T10:27:54Z planner-S-0709
No earlier plan or thread on S-0709 existed, so this is its first plan.

## Cause and proposed fix

`storygit.Rebase` (`flai/internal/storygit/sync.go`) names main by its branch name. `git rebase main` resolves it when the rebase starts. `checkUndone` then hands the same name to `FindUndone` (`flai/internal/storygit/undone.go`), and git resolves it again. A commit that reaches main in between is in `oldBase..main` but missing from `HEAD`, so the branch looks like it removes that commit's lines. The fix: read main's commit once with `git rev-parse` before the rebase, rebase onto that commit, and pass that commit to `FindUndone` and to `SettleChangelogs`. `res.Base` stays the branch name, so the output does not change.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2959: a real-git test that reproduces the race, then the fix in `Rebase` | `flai/internal/storygit/sync.go`, `undone.go`, `sync_test.go`, `undone_test.go` | none |
| 2 | T-2960: the "Main's change undone" paragraph in `docs/users/flai.md`, then `flai issue close I-0335` | `docs/users/flai.md`, the I-0335 file, `design/issues/summary.md` | T-2959 |

## Figures

- Forecast: 1h30m, delivery 2026-11-10T16:30Z. flai's figure was 2m, from size 2, which leaves out the real-git test and a close-out that runs integration and smoke.
- Cost of delay: 7.50 USD a week, as `flai cod` gives it.

## Assumptions

- The real-git test can force the race with a runner that commits to main right after the `git rebase` call. No change to the timing in production code is needed.
- The three-dot comparisons in `fix` and `checkSync` go through the merge base, so the race does not affect them. They stay as they are.
- `flai/cmd` callers (`stream_sync.go`, `stream.go`, `branch.go`) and `flai task done` get the fix through `Rebase` with no change of their own.

## Proposal: merge S-0712 into S-0709

S-0712 remediates I-0338: on S-0682, `flai task done`'s sync reported `undoes_main` for paths the branch never changed while main gained fe3b4e5. That is the same cause in the same `Rebase`. I recommend one story for both: add a third criterion to S-0709 to close I-0338 too, and cancel S-0712. Both stories are drafts, so the call is yours. I have not changed S-0712 or S-0709's criteria. If you agree, the cost of delay would also take I-0338's 3m a cycle, which would make it 15 USD a week.

### 2026-11-10T11:02:08Z orchestrator (recommendation)
**Recommended:** merge the two stories into S-0712, not into S-0709, and cancel S-0709.

- S-0712's fix (TH-1017) resolves main once, before the merge-base read, and uses that commit for the merge base, the rebase, and the check. This plan's fix pins it only from the rebase on.
- Carry this plan's one finding into S-0712's T-2963: `SettleChangelogs` also reads main by name after the rebase, so it takes the pinned commit too.
- S-0712 gains the criterion that closes I-0335. Its `time_lost_per_cycle` becomes 6m, I-0335's 3m added to I-0338's.

Both are drafts, and the merge changes scope and an input, so the call is yours. The full recommendation is on TH-1017.

Source: design/conventions/strategic-agents.md § As the orchestrator
