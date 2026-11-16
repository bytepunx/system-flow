---
id: S-0803
type: story
nature: improvement
title: Analyze the project and list the issues the report names
status: ready
owner: alex
created: 2026-11-16T06:00:53Z
updated: 2026-11-16T06:00:54Z
transitions:
  - to: ready
    at: 2026-11-16T06:00:54Z
    by: flai-serve
tags: []
touches: [wip/kanban/stories/S-0803-analyze-the-project-and-list-the-issues-the-report-names.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
formula: analysis-cadence
---
# S-0803 Analyze the project and list the issues the report names

## Goal

Run the analyzer on its focus, or on all three, and list the open issues its report names, for the operator to choose stories from at acceptance.

## Acceptance criteria

- [ ] Run the analyzer and write its report
- [ ] List the open issues the analyzer's report names

## Tasks
- T-3426 Run the analyzer and write its report
- T-3427 List the open issues the analyzer's report names

## Notes

Made by `flai formula run analysis-cadence` from `design/workflows/formulas/analysis-cadence.md`, with the inputs `focus=`.
