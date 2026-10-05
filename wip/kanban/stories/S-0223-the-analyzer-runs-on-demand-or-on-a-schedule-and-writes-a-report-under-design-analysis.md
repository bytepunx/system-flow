---
id: S-0223
type: story
nature: feature
title: The analyzer runs on demand or on a schedule and writes a report under design/analysis
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:17Z
updated: 2026-10-05T01:13:13Z
transitions: []
tags: [flai, dashboard]
touches: [flai/internal/serve, flai/internal/harness, flai/internal/hostapi, flai/cmd, flai/internal/mcpserver, ".claude/agents", template, design/analysis, flai/internal/guard, flai/internal/manifest, flai/internal/check, CLAUDE.md, design/system/repository-layout.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, flaiover/src/routes/docs]
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
  delivery: 2026-10-05T15:45:00Z
  basis: "Its own forecast of 2h30m; 20th in the pull order with an in-progress limit of 3, behind S-0266, S-0268, S-0253, S-0257, S-0244, S-0262, S-0258, S-0260, S-0212, S-0213, S-0214, S-0215, S-0216, S-0217, S-0218, S-0219, S-0220, S-0221 and S-0222."
  by: flai
  at: 2026-10-05T01:13:13Z
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

## Notes
