---
id: S-0533
type: story
nature: improvement
title: The orchestration page shows its settings panel first, above the action's state, the run, the last decisions, and the activity log, so the operator reaches the permissions and the policy without scrolling past them
status: done
parent: E-0030
owner: alex
created: 2026-10-28T21:01:22Z
updated: 2026-10-28T21:33:32Z
transitions:
  - to: ready
    at: 2026-10-28T21:16:58Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-28T21:17:05Z
    by: agent-S-0533
  - to: review
    at: 2026-10-28T21:27:35Z
    by: system-flow
  - to: done
    at: 2026-10-28T21:33:32Z
    by: orchestrator
tags: [flaiover, orchestrator, settings]
topics: [dashboard]
touches: [flaiover/src/routes/workflow/orchestrator, flaiover/src/lib/components, design/system/flaiover-dashboard.md, docs, flaiover/src/lib/components/StrategicAgentPanel.svelte, flaiover/src/lib/components/OrchestratorPanel.svelte, flaiover/src/lib/components/PlannerPanel.svelte, flaiover/src/lib/components/AnalyzerPanel.svelte, flaiover/src/lib/components/OrchestratorPanel.svelte.test.ts, flaiover/src/routes/workflow/orchestrator/+page.svelte, flaiover/src/routes/workflow/orchestrator/orchestrator.svelte.test.ts, flaiover/src/routes/workflow/planner/+page.svelte, flaiover/src/routes/workflow/planner/planner.svelte.test.ts, flaiover/src/routes/workflow/analyzer/+page.svelte, flaiover/src/routes/workflow/analyzer/analyzer.svelte.test.ts, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 642
  priced_by: harness
  turns:
    - day: 2026-10-28
      ceremony: 2
      hand_edits: 1
      work: 45
  models:
    - model: claude-opus-5-5
      input: 102
      output: 24697
      cache_read: 8634288
      cache_write: 229712
      cost: 4.0589
  strategic:
    - kind: orchestrator
      seconds: 1102
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 179
          output: 3166
          cache_read: 21252101
          cache_write: 47925
          cost: 4.6872
        - model: claude-sonnet-5-5
          input: 8
          output: 59
          cache_read: 80282
          cache_write: 42924
          cost: 0.085
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-10-28T21:13:23Z
  value: 12.5
  by: planner-S-0533
  at: 2026-10-28T21:15:39Z
forecast:
  duration: 30m
  delivery: 2026-10-29T11:52:00Z
  basis: "flai forecast: median 89 s per unit of size over 100 done improvement stories on claude-opus-5-5, times size 20 (4 criteria, 16 touches); 38th in the pull order with an in-progress limit of 5."
  by: planner-S-0533
  at: 2026-10-28T21:15:57Z
---
# S-0533 The orchestration page shows its settings panel first, above the action's state, the run, the last decisions, and the activity log, so the operator reaches the permissions and the policy without scrolling past them

## Goal

On `/workflow/orchestrator` the settings panel, `StrategicSettings.svelte` told `orchestration` with the permissions, `orchestration.policy`, and the `orchestration.release` keys (S-0229), sits below the `orchestrate` action's state, the run under way with its stream, the last decisions, and the activity log, so the operator scrolls past everything that changes to reach what they came to change. The operator asked on 2026-10-28 for the settings to come first. This story moves the settings panel to the top of the page, under the heading and beside the action's state and its Stop and Start, with the run, the decisions, and the activity log below it, and does the same on the planner's and the analyzer's pages where they share the panel (S-0259), so the three agents' pages read alike.

## Acceptance criteria
- [x] On `/workflow/orchestrator` the settings panel is the first section under the heading and the action's state; the run, the last decisions, and the activity log follow, and the page's test asserts the order
- [x] The planner's and the analyzer's pages put their settings first the same way, with their tests asserting it
- [x] The notice shown when `settings.get` has no `host.strategic` stays where the panel is
- [x] `flaiover-dashboard.md`'s rows for the three pages and the operator docs describe the order

## Tasks
- T-2278 StrategicAgentPanel renders a settings snippet between the host action and Running now, and the three agent panels pass it through
- T-2279 The Orchestrator page passes its settings to the panel, so they follow the host action's state and come before the run, the decisions, and the activity log
- T-2280 The Planner and Analyzer pages pass their settings to the panel, so they follow the host action's state as on the Orchestrator page
- T-2281 flaiover-dashboard.md and the flaiover guide say the three agent pages show their settings first, under the host action

## Notes

Asked by the operator in conversation on 2026-10-28 while reviewing the orchestrator's token use; not an ADR, a page layout.

### Planning

Touches, by source:

| Touch | Source | Why |
|---|---|---|
| `flaiover/src/routes/workflow/orchestrator` | declared | Folder kept as declared. T-2279's files narrow it in the claim (ADR-0096); no file beyond them is expected there |
| `flaiover/src/lib/components` | declared | Folder kept as declared. T-2278's files narrow it in the claim |
| `docs` | declared | Folder kept as declared. T-2281 names `docs/users/flaiover.md`, the one guide that describes the pages' order; `docs/operators/settings.md` names the pages but not their order |
| `design/system/flaiover-dashboard.md` | declared | The § Strategic agents' pages table (T-2281) |
| `flaiover/src/lib/components/StrategicAgentPanel.svelte` | layout | The shared panel renders the host action, then Running now; the settings snippet goes between them (T-2278) |
| `flaiover/src/lib/components/OrchestratorPanel.svelte`, `PlannerPanel.svelte`, `AnalyzerPanel.svelte` | layout | Each wraps the shared panel and passes the new snippet through (T-2278) |
| `flaiover/src/lib/components/OrchestratorPanel.svelte.test.ts` | layout | The panel-level order test (T-2278) |
| `flaiover/src/routes/workflow/orchestrator/+page.svelte`, `planner/+page.svelte`, `analyzer/+page.svelte`, and their `*.svelte.test.ts` | goal, criteria, layout | Each page renders `StrategicSettings` in a divider at its foot today; criteria 1 to 3 move it and test the order (T-2279, T-2280) |
| `docs/users/flaiover.md` | criteria | Its Planner, Orchestrator, and Analyzer sections say the settings come "below that" (T-2281) |

`flai touches suggest` listed co-change paths across the declared folders, led by `design/system/flai-cli.md`, `design/system/workflow.md`, and the ADR index. A page layout changes none of them, so none was added. A component examples file under `design/components/flaiover/` may be added by T-2278 if `flai check` asks; task done widens the touches then.

Overlap: S-0420 (ready, held on S-0419) replaces these three pages with one role page and touches the same panel, pages, and tests. Neither needs the other on main, so neither names the other in `after`. Whichever starts first holds the other on overlap; the pull order decides.

Forecast: 30m, delivery 2026-10-29T11:52:00Z, as `flai forecast` gives it (size 20: 4 criteria, 16 touches). It stands. The three declared folders double-count files already named, which overstates the size by 3. The four tasks, each closed with its own test run, and three page tests to reorder make up for that.

Cost of delay: 12.5 USD a week, as `flai cod` works it out from `time_lost_per_cycle: 5m` (set by the orchestrator on TH-0708, as recommended) at 150 USD an hour over a 168h cycle. It stands: a page layout with no revenue or penalty.

### Accepted by the orchestrator

- Verified: 2c8baf8f159e82d5874ccd33aff81b3b936e2a69
- Batch: 20261028-212920
- Enqueued verified: fd9ba3c72f9eefd355f8d8b6b8149cd44a2df286
- At: 2026-10-28T21:33:32Z

Verdict: meets all criteria
- 1: flaiover/src/lib/components/StrategicAgentPanel.svelte, flaiover/src/lib/components/OrchestratorPanel.svelte, flaiover/src/lib/components/OrchestratorPanel.svelte.test.ts, flaiover/src/routes/workflow/orchestrator/+page.svelte, flaiover/src/routes/workflow/orchestrator/orchestrator.svelte.test.ts
- 2: flaiover/src/lib/components/PlannerPanel.svelte, flaiover/src/lib/components/AnalyzerPanel.svelte, flaiover/src/routes/workflow/planner/+page.svelte, flaiover/src/routes/workflow/planner/planner.svelte.test.ts, flaiover/src/routes/workflow/analyzer/+page.svelte, flaiover/src/routes/workflow/analyzer/analyzer.svelte.test.ts
- 3: flaiover/src/routes/workflow/orchestrator/+page.svelte, flaiover/src/routes/workflow/planner/+page.svelte, flaiover/src/routes/workflow/analyzer/+page.svelte, flaiover/src/routes/workflow/orchestrator/orchestrator.svelte.test.ts, flaiover/src/routes/workflow/planner/planner.svelte.test.ts, flaiover/src/routes/workflow/analyzer/analyzer.svelte.test.ts
- 4: design/system/flaiover-dashboard.md, docs/users/flaiover.md
