---
id: TH-1100
title: "Plan for S-0726: flai stream sync continues a rebase whose stop the agent resolved"
anchor:
  path: wip/kanban/stories/S-0726-a-git-rebase-continue-run-by-hand-after-a-sync-s-stop-meets-later-stops-on.md
  item: S-0726
status: resolved
participants: [planner-S-0726, orchestrator]
created: 2026-11-12T15:57:17Z
updated: 2026-11-12T16:02:20Z
---

# TH-1100 Plan for S-0726: flai stream sync continues a rebase whose stop the agent resolved

On wip/kanban/stories/S-0726-a-git-rebase-continue-run-by-hand-after-a-sync-s-stop-meets-later-stops-on.md.

## Entries

### 2026-11-12T15:57:17Z planner-S-0726
S-0726 is planned: 29 file touches, a 2h forecast, a cost of delay value of 12.50 USD a week, and five tasks in four layers. The story names no `after`.

Proposed fix, from I-0351's one instance (S-0689): when `flai stream sync` is run again on a rebase whose stop the agent resolved and staged, it continues that rebase itself through `ContinueOverResolved`. Later stops on issue files and `design/issues/summary.md` are then settled as in a sync that never stopped. It finishes with the steps after a clean rebase. A stop that is not resolved is still refused. Every "run git rebase --continue" advice becomes "git add it, then run flai stream sync again".

| Layer | Task | Waits for | What it does |
|-------|------|-----------|--------------|
| 1 | T-3154 | none | ADR and `design/system/flai-cli.md`, `workflow.md` |
| 2 | T-3155 | T-3154 | `storygit/sync.go` continue, with a real-git test of I-0351 |
| 2 | T-3157 | T-3154 | user docs and the git and work-management conventions, with their template copies and `template/CHANGELOG.md` |
| 3 | T-3156 | T-3155 | command help, stop report, every other advice site, close-out scripts, regenerated reference, command test |
| 4 | T-3158 | T-3156, T-3157 | `flai issue close I-0351` |

T-3155 and T-3157 share no path, so they run together.

Assumptions:

- Plain `flai stream sync` continues, rather than a new `--continue` flag. The advice already ends "then run flai stream sync again", so agents need learn no new command. The story's agent may choose the flag in T-3154's ADR instead.
- The ADR number 0291 in the touches is predicted. `flai task done` narrows it to the file written.
- The undone check (ADR-0155) can recover the branch's old base and head from the rebase's state. T-3154 decides how.
- S-0709 also changes `flai stream sync`. This story reads nothing S-0709 writes, so the two are not ordered, and a trial merge judges them.
- Forecast raised from flai's 54m to 2h. The size counts many one-line advice changes, but not the depth of the rebase loop.

I would split, merge, or drop nothing. Nothing here needs an answer from the operator to go on.

### 2026-11-12T16:02:07Z orchestrator
**Approved.** The plan fits S-0726.

- The five tasks cover both criteria. T-3155 reproduces I-0351 with a real-git test (criterion 1), and T-3158 closes I-0351 (criterion 2).
- The layers hold. T-3155 and T-3157 share no path. Every task touch is a file, and no folder touch is kept.
- The 2h forecast and the 12.50 USD a week value stand, for the reasons in `### Planning`.
- The ADR decides between plain `flai stream sync` and a `--continue` flag. That choice is left to the story's agent in T-3154, as you propose.

S-0726 is finalized and in ready already.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-12T16:02:20Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: tasks, touches, and figures fit S-0726.
