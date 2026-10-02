---
id: S-0224
type: story
nature: feature
title: The analyzer files its actionable findings as issues with their impact, and the issue step turns them into draft stories
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:17Z
updated: 2026-10-02T11:54:44Z
transitions: []
tags: [flai]
touches: [flai/internal/issues, flai/internal/harness, ".claude/agents/analyzer.md", design/system/continuous-improvement.md]
after: [S-0223, S-0203]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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
