---
id: TH-0614
title: "Plan for S-0432: six tasks in four layers, touches widened to check.go and serve/curate.go, forecast 1h10m, cost of delay 95.02 USD a week"
anchor:
  path: wip/kanban/stories/S-0432-the-curator-runs-after-each-acceptance-and-on-a-schedule-marks-an-entry-the-design-now-records-absorbed-and-compacts-an-entry-no-one-has-cited-within-the-window.md
  item: S-0432
status: resolved
participants: [planner-S-0432, orchestrator]
created: 2026-10-23T17:11:22Z
updated: 2026-10-23T17:11:31Z
---

# TH-0614 Plan for S-0432: six tasks in four layers, touches widened to check.go and serve/curate.go, forecast 1h10m, cost of delay 95.02 USD a week

On wip/kanban/stories/S-0432-the-curator-runs-after-each-acceptance-and-on-a-schedule-marks-an-entry-the-design-now-records-absorbed-and-compacts-an-entry-no-one-has-cited-within-the-window.md.

## Entries

### 2026-10-23T17:11:22Z planner-S-0432
The story's six tasks stand as drafted; I revisited each against the outcome and widened two. No task to split, merge, or drop.

## Tasks, order, and layers

1. **Layer 1:** T-1831, the ADR (triggers, citation, absorption, compaction, nothing deleted). Touches `design/adrs` and its index.
2. **Layer 2, together:** T-1832 (`digest.schedule`, `digest.compact_after`, their `manifest.digest` check, the template manifest) and T-1833 (`flai digest absorb`, `cited`, `compact` in `flai/internal/digest/decay.go` and `flai/cmd/digest.go`). Both wait for T-1831; no path in common.
3. **Layer 3, together:** T-1834 (the curator scheduler in `flai/internal/serve`, after T-1832 for the schedule setting) and T-1835 (the prompt, the guard, and `curator.md` in both copies, after T-1833 for the commands it names). No path in common.
4. **Layer 4:** T-1836, the design and the user and operator docs, after T-1834 and T-1835.

## What changed

- T-1832 gains `flai/internal/check/check.go` and `flai/internal/check/digest_test.go`: `manifest.analysis` and `manifest.planning` are reported from `check.go`, which loops the manifest's `Errors()`, so `manifest.digest` is reported there too. The drafting planner left `check.go` out on the belief that the manifest package reports it; the code says otherwise.
- T-1834 gains `flai/internal/serve/curate.go` and its test (S-0431's files): the start it makes takes the trigger `accepted S-nnnn` or `schedule <spec>` and journals it, as `serve/analyze.go` does with `AgentRun.Trigger`.
- Story touches: the same four, now 32 files and one folder (`design/adrs`, shared, kept because the ADR's name is made by `flai adr new`).
- Forecast: 1h10m, from flai's 1h1m raised for the four touches at flai's rate; delivery 2026-10-24T11:11Z.
- Cost of delay value: 95.02 USD a week, the 1h10m share of E-0025's 600 USD a week over its eight open stories.

## Assumptions

- Acceptance detection reuses `moveTriggers` in `flai/internal/serve/replan.go` (a story moved to `done`), and coalescing follows its `enqueue`; the schedule reuses the `cadence` type there, as `analysis_schedule.go` does.
- `template/root/system-flow.yaml.tmpl` documents neither `planning` nor `analysis` today; T-1832 adds a commented `digest` block there because criterion 6 asks for it. If you would rather keep strategic settings out of the template manifest, as they are now, say so and I drop that touch and the words from criterion 6.
- `flai/internal/check/digest_test.go` is S-0428's file for the entry check; S-0428 lands first (this story is after S-0431), so T-1832 adds its cases to it rather than making a second digest test file.
- S-0430 delivers that an absorbed entry never loads in a pack and a compacted one loads as its title; criterion 3 relies on it.
- The write to `.claude/agents/curator.md` goes through `permission_prompt`; the story's agent makes it itself and waits for your answer.

Nothing here needs an answer to go on: the epic's cost of delay input was set after TH-0476's confirmation.

### 2026-10-23T17:11:28Z orchestrator
Approved: the plan fits S-0432.

- **Tasks and layers.** Six tasks in four layers. The pairs that share a layer touch disjoint paths: T-1832 with T-1833, and T-1834 with T-1835.
- **Touches.** 32 files. The four you added are grounded in the code: `manifest.digest` is reported from `check.go`, and the curator's start lives in `serve/curate.go`. The one folder, `design/adrs`, is kept for the ADR's file name.
- **Forecast.** 1h10m, flai's 1h1m raised at flai's rate for the four touches.
- **Cost of delay.** 95.02 USD a week, from `flai cod` on E-0025's input.

**The template manifest.** Keep the commented `digest` block. Criterion 6 asks for it, so dropping it would change the story's words, which is the operator's call.

**One note for review.** S-0432 writes `.claude/agents/curator.md`, so its acceptance will be the operator's (ADR-0106).

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-23T17:11:31Z orchestrator
Resolved: Plan approved under plan_backlog_stories.
