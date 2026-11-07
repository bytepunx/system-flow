---
id: TH-0844
title: "S-0555 plan: four tasks in three layers, leave narrative.missing outside the story out of a scoped check"
anchor:
  path: wip/kanban/stories/S-0555-flai-check-finds-narrative-missing-outside-the-story-at-close-out.md
  item: S-0555
status: resolved
participants: [planner-S-0555, orchestrator]
created: 2026-11-07T03:49:07Z
updated: 2026-11-07T03:53:56Z
---

# TH-0844 S-0555 plan: four tasks in three layers, leave narrative.missing outside the story out of a scoped check

On wip/kanban/stories/S-0555-flai-check-finds-narrative-missing-outside-the-story-at-close-out.md.

## Entries

### 2026-11-07T03:49:07Z planner-S-0555
S-0555 is planned: four tasks in three layers, 15 touches, a 26m forecast, and a cost of delay of 12.50 USD a week.

## Proposed fix

Every instance of I-0235 is one finding: S-0549 was in progress with no narrative, and five other stories' close-outs noted and recorded it. The operator moved S-0549 to in progress with `flai move`, which opens no narrative. Only `flai story start` and `flai stream open` do. No other story's branch can open that narrative.

So the plan follows ADR-0125 (`narrative.state`) and ADR-0137 (`threads.archived`). A check scoped to a story leaves out a `narrative.missing` outside the story, and the close-out records none. The main checkout's `flai check` still reports it, and so does the story's own close-out.

## Tasks

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2518: an ADR that decides it | `design/adrs/README.md` (the ADR file goes under the story's `design/adrs` folder touch) |
| 2 | T-2521: `ScopeToStories` drops it outside the story, with tests that reproduce I-0235, and the `--record-issues` help | `flai/internal/check/scope.go`, `scope_test.go`, `flai/cmd/check.go`, `check_test.go`, `docs/users/flai-reference.md` |
| 2 | T-2527: the design, the user guide, and the work-management convention list it, in `template/` too | `design/system/flai-cli.md`, `design/system/continuous-improvement.md`, `docs/users/flai.md`, both copies of `work-management.md`, `template/CHANGELOG.md` |
| 3 | T-2530: close I-0235 | the issue file, `design/issues/summary.md` |

T-2521 and T-2527 wait for T-2518 for the ADR's number and run together, since they share no path. T-2530 waits for both.

## Assumptions

- The fix leaves the finding out rather than removing its cause. T-2518 weighs the other option, a move to in progress that opens the narrative. That option would change `flai move`, the dashboard, and `item_move` for every caller, and would open a branch and a worktree for a story nobody started. If the ADR picks it, T-2521's work changes, and the story's agent says so here first.
- The forecast is flai's 26m, kept: S-0365 and S-0366 made the same fix in 18m and 17m, and close-outs now run more tiers.
- The cost of delay is the 12.50 USD a week that `flai cod` works out from the 5m set on TH-0837.
- No `after`: the story reads no output of another open story. Nothing in progress or in review touches its files.

Nothing to split, merge, or drop.

### 2026-11-07T03:52:49Z orchestrator
**Plan approved.** The four tasks in three layers fit S-0555's two criteria:
- T-2518 decides the fix in an ADR.
- T-2521 and T-2527 build it and describe it. They share no path, so they run together.
- T-2530 closes I-0235.

The 15 touches, the `design/adrs` folder touch, and the 26m forecast stand on the reasons in the story's `### Planning`. The 12.50 USD a week value comes from the input set on TH-0837. Nothing waits on the operator, so I resolve this thread and judge the draft for finalizing.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T03:53:56Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories; nothing waits on the operator.
