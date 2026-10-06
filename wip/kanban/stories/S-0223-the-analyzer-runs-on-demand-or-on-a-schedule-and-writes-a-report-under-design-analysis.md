---
id: S-0223
type: story
nature: feature
title: The analyzer runs on demand or on a schedule and writes a report under design/analysis
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:17Z
updated: 2026-10-06T11:14:34Z
transitions:
  - to: ready
    at: 2026-10-05T06:13:35Z
    by: alex
tags: [flai, dashboard]
topics: [analysis]
touches: [flai/internal/serve, flai/internal/harness, flai/internal/hostapi, flai/cmd, flai/internal/mcpserver, ".claude/agents", template, design/analysis, flai/internal/guard, flai/internal/manifest, flai/internal/check, CLAUDE.md, design/system/repository-layout.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, flaiover/src/routes/docs, flai/internal/analysis, ".claude/settings.json", design/README.md, docs/users/conventions.md, docs/users/index.md, design/adrs, design/system/project-manifest.md, design/system/flaiover-dashboard.md, docs/users/flai-reference.md, docs/users/flaiover.md, docs/operators/settings.md]
after: [S-0206, S-0207, S-0211]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 138.36
  by: planner-S-0223
  at: 2026-10-05T05:49:29Z
forecast:
  duration: 2h
  delivery: 2026-10-06T15:32:00Z
  basis: "Its own forecast of 2h; 5th in the pull order with an in-progress limit of 3, behind S-0221, S-0222, S-0224 and S-0226."
  by: flai
  at: 2026-10-06T11:14:34Z
---
# S-0223 The analyzer runs on demand or on a schedule and writes a report under design/analysis

## Goal

The analyzer reads the metrics, the design, the code, and the issues, and writes what it finds: bottlenecks (from cumulative flow, time in state, waiting, holds), gaps between `design/system` and the code, and technical and security risks. It writes a report; it does not author stories.

## Acceptance criteria
- [ ] An `analyze` host action gates it; `flai analyze [--focus bottlenecks|intent|risk]`, hostapi `analyze.run`, the MCP tool, and a manifest schedule (`analysis.schedule`) start a run; `flai serve` records it as other runs
- [ ] The report is `design/analysis/<date>-<focus>.md` with front matter (`title`, `updated`, `status`, `focus`, the window), sections per finding with evidence (metric figures, file paths, design sections quoted), severity, and estimated impact (time lost per cycle, or revenue or penalty when the design states them); `design/analysis/README.md` indexes them and the layout tables list the folder
- [ ] Its prompt primes with `--role analyze` and uses `flai stats --json`, `doc_search`, and the explorer; it edits nothing but its report, and `flai guard` enforces that
- [ ] Its log entry names the report and the run's cost; the dashboard's documents page shows the reports
- [ ] `design/system/strategic-agents.md` and the user guide describe it; tests cover a run, the guard, and the schedule

## Tasks
- T-0949 system-flow.yaml takes analysis.agent and analysis.schedule, and flai check reports a bad one
- T-0957 The analyzer's prompt, its claude-code definition, and settings that run flai guard on its file edits
- T-0960 design/analysis holds the analyzer's reports and their index, flai check validates them, and the layout tables list the folder
- T-0968 flai serve starts the analyzer behind the analyze host action, records its run, and logs the report it wrote with the run's cost
- T-0971 flai analyze [--focus] and the MCP tool analyze start an analyzer run through flai serve
- T-0972 flai serve starts the analyzer when analysis.schedule comes round, while the analyze host action is on
- T-0973 flai guard holds an analyzer session to its report: file edits only under design/analysis, flai reads, threads, and activity_log
- T-0974 The dashboard's documents page lists and shows the analyzer's reports under design/analysis
- T-0975 An ADR, strategic-agents.md, flai-cli.md, the manifest design, the template changelog, and the user and operator guides describe the analyzer

## Notes

### Planning

Planned by planner-S-0223 on 2026-10-05.

Touches:

- Declared, kept: `flai/internal/serve`, `flai/internal/harness`, `flai/internal/hostapi`, `flai/cmd`, `flai/internal/mcpserver`, `.claude/agents`, `template`, `design/analysis`, `flai/internal/guard`, `flai/internal/manifest`, `flai/internal/check`, `CLAUDE.md`, `design/system/repository-layout.md`, `design/system/strategic-agents.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, `flaiover/src/routes/docs`.
- Layout: `flai/internal/analysis`, a new package beside `flai/internal/experiment`, which is how `design/experiments` was added (S-0194); `design/README.md`, `docs/users/conventions.md`, and `docs/users/index.md`, the other layout tables that list `design/experiments/`.
- Design (ADR-0082 § 7, the planner's precedent): `.claude/settings.json`, whose `Edit|Write|NotebookEdit` hook runs the guard only for `FLAI_ROLE=plan` and must cover `analyze` too.
- Co-change (seeded by the declared touches, 491 of 866 commits): `design/system/flaiover-dashboard.md` (16%), `docs/users/flai-reference.md` (16%), `docs/users/flaiover.md` (12%), `design/adrs` (its README 12%), `docs/operators/settings.md` (8%), and `design/system/project-manifest.md` (4%), which the new `analysis` block, the host action, and the documents page reach.
- Left out: the co-change list's issues, conventions, workflow, and work-hierarchy documents, and `flaiover/src/lib/server/agent.ts`; the criteria do not reach them. The Analyzer page and its settings are S-0228's and S-0229's.

Forecast: 2h, against flai's 1h13m (132 s per unit over 14 large feature stories on this model, times size 33). Adjusted up because the story joins what S-0208, the planner's host action and run (65m of agent time), and S-0211, the planner's schedule (48m), each did, and adds a report folder with its own check, as S-0194 did for experiments, and a guard role; S-0218's planner set 1h40m for a story of the same shape without the schedule or the folder. Delivery 18:10Z: it starts about 15:07Z, 13th in the pull order with an in-progress limit of 3, and takes 2h times the cycle factor 1.52.

Cost of delay: 138.36 USD a week, flai's figure, kept. The story has no inputs of its own, so it takes its share of E-0016's 1500 USD a week, 2h of 21h41m forecast over the epic's 17 open stories without inputs. It replaces planner-E-0016's 121.95, which was worked out over more open stories and a 2h30m duration.

Plan: four layers. T-0949 and T-0957 wait for nothing; T-0960 waits for T-0949 and T-0968 for both; T-0971 and T-0972 wait for T-0968, and T-0973 and T-0974 for T-0960; T-0975 waits for the seven code tasks. The plan's thread on this story gives the assumptions.
