---
id: TH-0143
title: "S-0224 plan: five tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0224-the-analyzer-files-its-actionable-findings-as-issues-with-their-impact-and-the-issue-step-turns-them-into-draft-stories.md
  item: S-0224
status: open
participants: [planner-S-0224]
created: 2026-10-05T05:45:56Z
updated: 2026-10-05T05:45:56Z
---

# TH-0143 S-0224 plan: five tasks in four layers

On wip/kanban/stories/S-0224-the-analyzer-files-its-actionable-findings-as-issues-with-their-impact-and-the-issue-step-turns-them-into-draft-stories.md.

## Entries

### 2026-10-05T05:45:56Z planner-S-0224
Plan for S-0224 (planner-S-0224). Touches, forecast, and cost of delay are written, with their reasons under `### Planning` in the story's Notes.

Tasks, layer by layer:

- Layer 1: T-0918, `flai issue new` and `bump` take the impact (`--revenue-per-week`, `--penalty-per-week`, `--time-lost-per-cycle`, `--evidence`, written as `## Impact`) and `--report` (an instance line and a `## Remediation` link); with `--report`, an open issue of the same title is bumped, not duplicated.
- Layer 2, together: T-0922, `flai issue story` links the report from the draft story and the impact carries over (after T-0918); T-0925, MCP `issue_new` and `issue_bump` (after T-0918). They share no path.
- Layer 3: T-0930, the analyzer's prompt, its `.claude/agents/analyzer.md` and template copy, and the `analyze` guard rules let it file and bump issues and refuse story creation; the report links each issue (after T-0918 and T-0925).
- Layer 4: T-0937, `continuous-improvement.md`, `strategic-agents.md`, `flai-cli.md`, and the user guides (after T-0922 and T-0930).

Figures: forecast 40m, as flai gives it (S-0203 took 25m on the same code in four tasks); cost of delay 76.01 USD a week, the story's share of E-0016's 1500, replacing planner-E-0016's 60.98.

Assumptions:

- S-0223 makes the analyzer's prompt, definition, `analyze` guard rules, and its report under `design/analysis/`, and this story extends them. If S-0223 places the prompt outside `flai/internal/harness`, T-0930's touches should follow it.
- Writing `.claude/agents/analyzer.md` in this repository needs you to paste it, since agents' writes there are refused.
- Deduplication is the analyzer's judgement against `flai issue list --json`, with flai bumping an exact-title match when `--report` is given.
- I added the MCP tools (T-0925) so that an analyzer without a shell can file issues too. If the CLI alone is enough, I would drop T-0925 and `flai/internal/mcpserver` from the touches. My recommendation is to keep it.

I added topic `template` to the story for T-0930. Nothing here needs an answer before the story is finalized.
