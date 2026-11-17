---
id: TH-1282
title: "Plan for S-0833: the publish commit waits out a held index lock, three tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0833-flai-release-pending-fails-when-another-process-holds-the-main-checkout-s-git.md
  item: S-0833
status: resolved
participants: [planner-S-0833, orchestrator]
created: 2026-11-17T08:30:38Z
updated: 2026-11-17T08:33:41Z
---

# TH-1282 Plan for S-0833: the publish commit waits out a held index lock, three tasks in two layers

On wip/kanban/stories/S-0833-flai-release-pending-fails-when-another-process-holds-the-main-checkout-s-git.md.

## Entries

### 2026-11-17T08:30:38Z planner-S-0833
Plan: route the publish commit of `flai release --pending` through `storygit.RunPastIndexLock`, the helper that fixed the same race in `flai accept` (I-0100, S-0307). Today `computeApplyAndTagPending` in `flai/cmd/release.go` runs `git add -A` and `git commit` straight through the runner.

Tasks:

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3535 publish commit waits out a held lock, with a test that holds it | `flai/cmd/release.go`, `flai/cmd/release_lock_test.go` | none |
| 2 | T-3538 `docs/users/flai.md` says so | `docs/users/flai.md` | T-3535 |
| 2 | T-3539 close I-0359 | I-0359's file, `design/issues/summary.md` | T-3535 |

Figures:

- Forecast 18m, as `flai forecast` gives it.
- Cost of delay 40 USD a week, as `flai cod` works it from your 16m input.

The reasons are in the story's Notes under Planning.

Assumptions:

- A wait of about nine seconds is enough. Both instances saw the lock gone within seconds.
- `git status`, `git tag`, and `git push` take no index lock, so the story leaves them alone.
- The story has no `after`. S-0811 (ready) also changes the publish path, and a trial merge judges that overlap.

Out of scope. I propose these and have drafted neither:

1. The publish runs `git add -A`, so it commits every uncommitted file in the main checkout. That is how I-0359's second instance left `wip/` and `design/digest/` staged. My recommendation is to leave this as it is in S-0833: the publish may mean to commit what main holds, as acceptance does under ADR-0248. If you want the publish to stage only its own files, as `flai accept` does, that is a separate story.
2. About a dozen other commands commit in the main checkout without the helper. Among them are `flai/cmd/formula.go`, `flai/cmd/subproject.go`, `flai/internal/itemedit/itemedit.go`, `flai/internal/itemnew/itemnew.go`, `flai/internal/serve/replan.go`, and `flai/internal/adr/adr.go`. No issue records a failure in any of them yet, so I recommend a story only once one does.

### 2026-11-17T08:32:57Z orchestrator
Approved. The plan fits S-0833:

- T-3535 routes the publish's `git add -A` and `git commit` through `storygit.RunPastIndexLock`, with a test that holds the lock, which meets criterion 1. T-3538 documents it and T-3539 closes I-0359, which meets criterion 2. Both wait for T-3535 and share no path, so they run together.
- Every touch is a file one task names. There is no folder touch and no `after`. A trial merge judges the overlap with S-0811.
- The 18m forecast comes from `flai forecast` and the 40 USD a week value from `flai cod`.

On the two out-of-scope proposals, I take your recommendation. Neither becomes part of S-0833, and I draft no story for either:

1. The publish's `git add -A` stays as it is in S-0833. A story to stage only the publish's own files is for the operator to ask for.
2. The other commands that commit in the main checkout get a story once an issue records a failure in one of them.

I am finalizing S-0833 under `finalize_drafts`.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-17T08:33:41Z orchestrator
Resolved: Plan approved under plan_backlog_stories; the two out-of-scope proposals stay out of S-0833, as the planner recommended.
