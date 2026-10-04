---
id: S-0224
type: story
nature: feature
title: The analyzer files its actionable findings as issues with their impact, and the issue step turns them into draft stories
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:17Z
updated: 2026-10-04T04:48:22Z
transitions: []
tags: [flai]
touches: [flai/internal/issues, flai/internal/harness, ".claude/agents/analyzer.md", design/system/continuous-improvement.md, flai/internal/mcpserver, template/root/.claude/agents/analyzer.md, design/system/strategic-agents.md, docs/users/flai.md]
after: [S-0223, S-0203]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 60.98
  by: planner-E-0016
  at: 2026-10-04T04:48:22Z
forecast:
  duration: 1h15m
  delivery: 2026-10-05T09:30:00Z
  basis: "Builds on S-0203's issue-to-draft path (1555 s) with an Impact section, deduplication and prompt rules; serial on flai/internal/harness after S-0222."
  by: planner-E-0016
  at: 2026-10-04T04:44:53Z
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

## Notes
