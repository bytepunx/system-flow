---
id: TH-1239
title: "Plan for S-0798: the close-out syncs when its first flai verify stops at sync, in the scripts and a container story's close-out, in two layers"
anchor:
  path: wip/kanban/stories/S-0798-a-close-out-asked-for-a-story-whose-branch-is-behind-main-stops-at-sync-without.md
  item: S-0798
status: resolved
participants: [planner-S-0798, orchestrator]
created: 2026-11-16T01:45:08Z
updated: 2026-11-16T01:46:42Z
---

# TH-1239 Plan for S-0798: the close-out syncs when its first flai verify stops at sync, in the scripts and a container story's close-out, in two layers

On wip/kanban/stories/S-0798-a-close-out-asked-for-a-story-whose-branch-is-behind-main-stops-at-sync-without.md.

## Entries

### 2026-11-16T01:45:08Z planner-S-0798
The plan has four tasks in two layers. They fix the cause: a close-out asked for a branch behind main now syncs before it verifies, instead of leaving a failed report.

## Tasks and layers

| Layer | Task | What it changes | after |
|-------|------|-----------------|-------|
| 1 | T-3399 | `scripts/close-out.sh` and `template/root/scripts/close-out.sh` sync and verify again, within the three rounds, when the round's first `flai verify` stops at sync. A new test, `flai/internal/closeout/script_test.go`, reproduces I-0383 with a stub `flai.sh` | none |
| 1 | T-3400 | A container story's close-out, `flai/internal/serve/verifyjob_container.go` (ADR-0274), does the same on the host. Its test is in `verifyjob_test.go` | none |
| 1 | T-3401 | A new ADR (0309 if free) refining ADR-0156 and ADR-0291 | none |
| 2 | T-3402 | The close-out rule in `work-management.md` and its template copy, `template/CHANGELOG.md`, `docs/users/flai.md`, and both `scripts/README.md`. Then it closes I-0383 | T-3399, T-3400, T-3401 |

Layer 1's three tasks share no file and can run together. T-3399 and T-3400 print the same round line: `close-out: the branch is behind main; flai stream sync <story> (round <n> of <rounds>)`. T-3402 waits for all three because it describes their behaviour and cites the ADR.

## Figures

- Forecast 1h, delivery 2026-11-16T04:56Z. I raised it from flai's 33m. The fix covers two scripts and a Go loop, the script needs its first test harness, and a `flai/` plus `template/` story's close-out runs integration and smoke.
- Cost of delay 300 USD a week, from flai's input of 2h a cycle. It stands as flai gives it.

## Assumptions

- The issue's interim alternatives are not planned. Those are keeping a passing report when a close-out stops at sync before any tier ran, or having the queue's `verify_failed` tell that case apart. Once the close-out syncs first, a failed report at sync is left only when the sync conflicts, and the story then needs its agent. Say if you want either one as well; it would be a fifth task in `flai/internal/verify` or `flai/internal/mergequeue/eligible.go`.
- The `verify_behind_main` blocker's words (`flai/internal/preview/accept.go`, line 671) and `flai accept`'s help already say the close-out syncs the branch onto main. They become true with this change, so the plan leaves them as they are. ADR-0291 is accepted, so the new ADR refines it rather than editing it.
- `flai stream sync` refuses uncommitted changes, so the close-out commits what is outstanding before the new sync, as its commit step does.
- The ADR's number depends on S-0796's ADR-0308 landing first. `flai task done` widens the touches to whatever file is committed.

No question is open. I'd split, merge, or drop nothing.

### 2026-11-16T01:46:24Z orchestrator
Approved: the plan fits S-0798.

- **Tasks.** T-3399 and T-3400 fix the cause I-0383 records in both close-outs, the script and a container story's, each with a test that reproduces it first. That meets criterion 1. T-3402 closes I-0383, which is criterion 2.
- **Layers.** Layer 1's three tasks share no path. T-3402 waits for them because it describes their behaviour and cites the ADR.
- **Touches.** Every task touch is in the story's touches. No folder touch, no `after`.
- **Figures.** Forecast 1h, raised from flai's 33m with its reason stated. Value 300 USD a week, from the input flai carried from I-0383. No input is asked for.
- **Assumptions.** Leaving out the issue's interim alternatives is right: once the close-out syncs first, a stop at sync remains only on a conflict. S-0799 takes up the stale-report wording on its own.

One note for the story's agent: S-0797's plan also names ADR-0309. `flai adr new` takes the next free number, so whichever lands second gets 0310, and `flai task done` widens the touches.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-16T01:46:42Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, layers, touches, and figures fit S-0798.
