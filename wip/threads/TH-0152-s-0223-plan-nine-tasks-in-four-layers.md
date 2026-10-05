---
id: TH-0152
title: "S-0223 plan: nine tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0223-the-analyzer-runs-on-demand-or-on-a-schedule-and-writes-a-report-under-design-analysis.md
  item: S-0223
status: resolved
participants: [planner-S-0223, alex]
created: 2026-10-05T05:50:04Z
updated: 2026-10-05T05:51:49Z
---

# TH-0152 S-0223 plan: nine tasks in four layers

On wip/kanban/stories/S-0223-the-analyzer-runs-on-demand-or-on-a-schedule-and-writes-a-report-under-design-analysis.md.

## Entries

### 2026-10-05T05:50:04Z planner-S-0223
Plan for S-0223: nine tasks in four layers. The touches, the forecast, and the cost of delay are under ### Planning in the story's Notes.

Layers:

- Layer 1: T-0949 (`analysis.agent` and `analysis.schedule` in the manifest, `manifest.analysis` in flai check) and T-0957 (the analyzer's prompt, `.claude/agents/analyzer.md`, and the settings hook for `FLAI_ROLE=analyze`). They wait for nothing and share no path.
- Layer 2: T-0960 (the `design/analysis` folder, its README index, the report check, and the layout tables) waits for T-0949, because both change `check.go`. T-0968 (the `analyze` host action, `analyze.run`, the serve run, and the activity entry naming the report) waits for T-0949 and T-0957.
- Layer 3: T-0971 (`flai analyze --focus` and the MCP tool) and T-0972 (the schedule) wait for T-0968. T-0973 (the guard role `analyze`) and T-0974 (the documents page) wait for T-0960.
- Layer 4: T-0975 (the ADR, strategic-agents.md, flai-cli.md, project-manifest.md, the user and operator guides, the template changelog) waits for the seven code tasks.

Assumptions:

- A run with no `--focus` covers all three focuses and is named `<date>-all.md`. The window is front matter `from` and `to`.
- `flai check` validates each report's front matter, and reports a report that `design/analysis/README.md` does not list. This follows the `design/experiments` precedent.
- The guard lets the analyzer write only under `design/analysis/`, with no item writes and no issues. Filing issues is S-0224's work, and that story widens the guard.
- The activity entry names the newest file under `design/analysis/` changed during the run.
- The documents page already lists every folder under design. T-0974 is mainly tests that prove it, plus the dashboard docs.
- The dashboard's Analyze button, the Analyzer page, and the settings for `analysis.*` are S-0228's and S-0229's, not this story's.
- A story's agent cannot write `.claude/` (I-0069). T-0957 sends the operator the whole of `.claude/settings.json` and `.claude/agents/analyzer.md` to paste.
- S-0223 and S-0218 both touch the guard, the harness, serve, and `.claude/settings.json`. The pull hold will run them one after the other, so I gave S-0223 no `after` on S-0218.

Nothing to split, merge, or drop.

### 2026-10-05T05:51:49Z alex
Resolved.
