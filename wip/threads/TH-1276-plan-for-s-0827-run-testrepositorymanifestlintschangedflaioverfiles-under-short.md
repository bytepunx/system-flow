---
id: TH-1276
title: "Plan for S-0827: run TestRepositoryManifestLintsChangedFlaioverFiles under -short, two tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0827-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md
  item: S-0827
status: resolved
participants: [planner-S-0827, orchestrator]
created: 2026-11-17T08:28:54Z
updated: 2026-11-17T08:29:57Z
---

# TH-1276 Plan for S-0827: run TestRepositoryManifestLintsChangedFlaioverFiles under -short, two tasks in two layers

On wip/kanban/stories/S-0827-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md.

## Entries

### 2026-11-17T08:28:54Z planner-S-0827
S-0827 is planned with two tasks in two layers.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3523 Run TestRepositoryManifestLintsChangedFlaioverFiles under -short | `flai/internal/verify/select_test.go` | none |
| 2 | T-3524 Close I-0327 with what fixed it | the I-0327 issue file, `design/issues/summary.md` | T-3523 |

T-3524 waits for T-3523 because its close reason names the fix.

The fix: drop the `testing.Short()` gate from the test and keep its "monorepo not present" skip. The test reads only `system-flow.yaml` and the file tree through `os.DirFS`, so it needs no git.

Figures:

- Forecast 30m, up from flai's 4m, because flai's figure leaves out the close-out. S-0693 made the same change in 21m with a 6m47s close-out.
- Cost of delay value 22.50 USD a week, as `flai cod` gives it.

Assumptions:

- This redoes S-0693's fix. S-0693 was cancelled (TH-1066) because its branch undid other work on main, not because the fix was wrong. Its gate is still on main.
- The test's `want` may be out of date with today's manifest. If so, T-3523 sets `want` to what `system-flow.yaml` selects now, and records that in its notes.
- The `-short` gate in `paths_test.go`'s `realRepo` runs real git, so it stays.

I would split, merge, or drop nothing.

### 2026-11-17T08:29:41Z orchestrator
Approved. The plan fits S-0827:

- T-3523 meets criterion 1, and T-3524 meets criterion 2. T-3524 waits for T-3523 because its close reason names the fix.
- Each touch is a file that one of the two tasks names. There is no folder touch and no `after`.
- The 30m forecast rests on S-0828's twin, S-0693, which took 21m. The 22.50 USD a week value comes from `flai cod`.

I am finalizing S-0827 under `finalize_drafts`.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-17T08:29:57Z orchestrator
Resolved: Plan approved under plan_backlog_stories. Correction to my reply: S-0693 is S-0827's own earlier attempt at this fix, not S-0828's.
