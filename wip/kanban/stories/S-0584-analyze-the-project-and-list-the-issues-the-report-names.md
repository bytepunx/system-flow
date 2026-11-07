---
id: S-0584
type: story
nature: improvement
title: Analyze the project and list the issues the report names
status: ready
owner: alex
created: 2026-11-07T08:03:23Z
updated: 2026-11-07T08:03:40Z
transitions:
  - to: ready
    at: 2026-11-07T08:03:40Z
    by: system-flow
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
formula: analysis-cadence
---
# S-0584 Analyze the project and list the issues the report names

## Goal

Run the analyzer on its focus, or on all three, and list the open issues its report names, for the operator to choose stories from at acceptance.

## Acceptance criteria

- [ ] Run the analyzer and write its report
- [ ] List the open issues the analyzer's report names

## Tasks
- T-2545 Run the analyzer and write its report
- T-2547 List the open issues the analyzer's report names

## Notes

Made by `flai formula run analysis-cadence` from `design/workflows/formulas/analysis-cadence.md`, with the inputs `focus=`.
