---
id: TH-0945
title: "S-0623 plan: the guard learns the read forms of git tag, worktree, and branch; three tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0623-flai-guard-refuses-the-orchestrator-git-tag-list-and-git-worktree-list-which.md
  item: S-0623
status: resolved
participants: [planner-S-0623, orchestrator]
created: 2026-11-09T21:56:17Z
updated: 2026-11-09T21:57:17Z
---

# TH-0945 S-0623 plan: the guard learns the read forms of git tag, worktree, and branch; three tasks in two layers

On wip/kanban/stories/S-0623-flai-guard-refuses-the-orchestrator-git-tag-list-and-git-worktree-list-which.md.

## Entries

### 2026-11-09T21:56:17Z planner-S-0623
**Proposed solution (the goal asks for one).** `flai/internal/guard/guard.go` judges git by subcommand name against `gitReads`, which has no `tag`, `worktree`, or `branch`. The orchestrator's rule also lists `tag` under `gitPublishes`. The guard will judge these three by form, as `flagReads` does for flai's commands:

- `git tag` that lists passes.
- `git worktree list` passes.
- `git branch` that lists passes.
- Every writing form keeps its refusal. `git tag v1.2.3` keeps the orchestrator's `publishesThrough` reason.

**Tasks and layers.**

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2766 The guard passes git tag, worktree, and branch in the forms that only read | `flai/internal/guard/guard.go`, `guard_test.go` | none |
| 2 | T-2767 Document the git forms flai guard passes as reads | `flai/cmd/guard.go`, `docs/users/flai-reference.md`, `design/system/flai-cli.md` | T-2766 |
| 2 | T-2768 Close I-0289 | the issue file, `design/issues/summary.md` | T-2766 |

T-2767 and T-2768 share no path, so they run together.

**Figures.**

- Forecast: 1h30m, raised from flai's 3m, which counted no touches.
- Cost of delay: 75 USD a week, from the issue's 30m `time_lost_per_cycle`. Unchanged.

**Assumptions.**

1. `git branch` is included even though the instances show only `tag` and `worktree`. The analysis names it alongside them, and it has the same split between reading and writing forms. Say if it should be left out.
2. The rule is shared, so sub-agents, the planner, the analyzer, and the curator gain these reads too, not just the orchestrator.
3. No ADR. This widens a list of reads within the existing rule, as S-0612 did for `HostAPIReads`.
4. The redirection refusals (`2>&1`) in the same analysis belong to I-0230 and stay out of scope.
5. The orchestrator asked for this run because the operator wrote on a thread called "asked", but S-0623 had no thread and no earlier plan. I planned it from scratch.

### 2026-11-09T21:56:50Z orchestrator
Approved. The tasks, touches, and figures fit S-0623.

- The plan answers the goal's request for a solution. T-2766 changes the guard, T-2767 documents it, and T-2768 closes I-0289. Together they cover criteria 1 and 2.
- The touches match the tasks' paths. No folder touch.
- 1h30m and 75 USD a week fit a guard change with table tests and a regenerated reference.
- Assumption 1 stands: I-0289's evidence names `branch` beside `tag` and `worktree` ("gitReads has no tag, worktree, or branch").
- Assumptions 2 to 4 stand as written.

I finalize S-0623 now. `flai promote --drafts` lists it as complete.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-09T21:57:17Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: tasks, touches, forecast, and value fit S-0623.
