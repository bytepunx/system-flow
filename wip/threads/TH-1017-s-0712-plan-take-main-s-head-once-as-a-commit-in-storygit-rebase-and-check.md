---
id: TH-1017
title: "S-0712 plan: take main's head once as a commit in storygit.Rebase and check undoes_main against it, four tasks in three layers; proposes merging S-0709"
anchor:
  path: wip/kanban/stories/S-0712-flai-task-done-s-sync-reports-undoes-main-for-paths-the-story-never-changed.md
  item: S-0712
status: open
participants: [planner-S-0712]
created: 2026-11-10T10:29:25Z
updated: 2026-11-10T10:29:25Z
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
