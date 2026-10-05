---
id: S-0224
type: story
nature: feature
title: The analyzer files its actionable findings as issues with their impact, and the issue step turns them into draft stories
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:17Z
updated: 2026-10-05T05:47:33Z
transitions: []
tags: [flai]
topics: [template]
touches: [flai/internal/issues, flai/internal/harness, ".claude/agents/analyzer.md", design/system/continuous-improvement.md, flai/internal/mcpserver, template/root/.claude/agents/analyzer.md, design/system/strategic-agents.md, docs/users/flai.md, flai/cmd/issue.go, flai/cmd/issue_test.go, flai/internal/guard, design/system/flai-cli.md, docs/users/flai-reference.md, template/CHANGELOG.md]
after: [S-0223, S-0203]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 76.01
  by: planner-S-0224
  at: 2026-10-05T05:44:32Z
forecast:
  duration: 40m
  delivery: 2026-10-05T19:55:00Z
  basis: "Its own forecast of 40m; 14th in the pull order with an in-progress limit of 3, behind S-0276, S-0217, S-0218, S-0219, S-0220, S-0221, S-0222, S-0226, S-0212, S-0213, S-0214, S-0215, S-0216 and S-0223."
  by: flai
  at: 2026-10-05T05:47:33Z
---
# S-0224 The analyzer files its actionable findings as issues with their impact, and the issue step turns them into draft stories

## Goal

The analyzer's report is for reading; its actionable findings become issues so that the continuous-improvement flow (S-0198, the draft story from an issue) turns them into stories with cost of delay already in place, and the analyzer never authors stories itself.

## Acceptance criteria
- [ ] For each actionable finding the analyzer files an issue (`flai issue new`, class from the finding: `defect`, `efficiency`, `impression` for risks without a measured instance) with an `## Impact` section giving time lost per cycle, or revenue or penalty, and the evidence, deduplicated against open issues by linking the report to an existing issue and bumping it instead
- [ ] The story the issue step creates from such an issue carries the impact as cost of delay inputs and links the report
- [ ] The report links each issue it filed; the issue's Remediation section links the report
- [ ] Tests cover a new issue, a bumped duplicate, and the inputs carried over

## Tasks
- T-0918 flai issue new and bump write an issue's impact and link the analyzer's report
- T-0922 flai issue story links the analyzer's report from the draft story it makes and carries the impact over
- T-0925 MCP issue_new and issue_bump file and bump issues with their impact and report
- T-0930 The analyzer's prompt and definition file each actionable finding as an issue and its report links each issue, and flai guard lets it
- T-0937 Design and user guide say how the analyzer files findings as issues and how the issue step carries their impact and report

## Notes

### Planning

Planned by planner-S-0224 on 2026-10-05.

Touches:

- Declared, kept: `flai/internal/issues`, `flai/internal/harness`, `.claude/agents/analyzer.md`, `design/system/continuous-improvement.md`, `flai/internal/mcpserver`, `template/root/.claude/agents/analyzer.md`, `design/system/strategic-agents.md`, `docs/users/flai.md`.
- Layout: `flai/cmd/issue.go` and `flai/cmd/issue_test.go`, where `flai issue new`, `bump`, and `story` live; today `new` and `bump` take no impact or report.
- Design: `flai/internal/guard`, since `design/system/strategic-agents.md` gives each strategic role guard rules of its own and the analyzer's must pass its issue calls and refuse story creation.
- Co-change: `design/system/flai-cli.md` (57% of the commits that changed the declared touches) and `docs/users/flai-reference.md` (22%), which document the issue commands and MCP tools; `template/CHANGELOG.md` (7%), for the template's analyzer definition.
- Left out: `docs/operators`, `design/system/flaiover-dashboard.md`, and `docs/users/flaiover.md` co-change often but nothing here changes the dashboard or settings.

Figures:

- Forecast 40m, as `flai forecast` gives it once the touches above are in (size 18: 4 criteria, 14 touches, 132 s a unit), up from the 15m it gave on the 8 declared touches. Kept rather than adjusted: S-0203, which built the issue step on the same code in four tasks, took 25m, and this story adds a fifth task for the MCP tools and the guard. The delivery, 2026-10-05T20:01Z, is flai's, behind S-0223 and 14 others in the pull order.
- Cost of delay 76.01 USD a week, from `flai cod`: the story has no inputs of its own, so it takes its share of E-0016's 1500 USD a week, 1h15m of 24h40m over 17 open stories. It replaces the 60.98 planner-E-0016 set on 2026-10-04, which came from an earlier share.

Assumptions:

- S-0223, which this story waits for, makes the analyzer's prompt, its `.claude/agents/analyzer.md` definition, its `analyze` guard rules, and the report under `design/analysis/`; this story extends them.
- S-0203 already reads `## Impact` into cost of delay inputs; what is missing is a way for the analyzer to write it, the report link, and the deduplication.
- Deduplication is the analyzer's judgement from `flai issue list --json`, backed by flai bumping an open issue of the same title when `--report` is given, as close-out findings are recorded once.
