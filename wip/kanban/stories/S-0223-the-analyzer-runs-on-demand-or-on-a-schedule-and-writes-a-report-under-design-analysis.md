---
id: S-0223
type: story
nature: feature
title: The analyzer runs on demand or on a schedule and writes a report under design/analysis
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:17Z
updated: 2026-10-05T05:47:33Z
transitions: []
tags: [flai, dashboard]
touches: [flai/internal/serve, flai/internal/harness, flai/internal/hostapi, flai/cmd, flai/internal/mcpserver, ".claude/agents", template, design/analysis, flai/internal/guard, flai/internal/manifest, flai/internal/check, CLAUDE.md, design/system/repository-layout.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, flaiover/src/routes/docs, flai/internal/analysis, ".claude/settings.json", design/README.md, docs/users/conventions.md, docs/users/index.md, design/adrs, design/system/project-manifest.md, design/system/flaiover-dashboard.md, docs/users/flai-reference.md, docs/users/flaiover.md, docs/operators/settings.md]
after: [S-0206, S-0207, S-0211]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 121.95
  by: planner-E-0016
  at: 2026-10-04T04:48:21Z
forecast:
  duration: 2h30m
  delivery: 2026-10-05T18:54:00Z
  basis: "Its own forecast of 2h30m; 13th in the pull order with an in-progress limit of 3, behind S-0276, S-0217, S-0218, S-0219, S-0220, S-0221, S-0222, S-0226, S-0212, S-0213, S-0214, S-0215 and S-0216."
  by: flai
  at: 2026-10-05T05:47:33Z
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

## Notes
