---
id: T-3611
type: task
nature: improvement
title: Run the analyzer and write its report
status: in-progress
parent: S-0850
owner: alex
created: 2026-11-23T06:00:22Z
updated: 2026-11-27T22:34:34Z
transitions:
  - to: ready
    at: 2026-11-27T22:34:34Z
    by: agent-S-0850
  - to: in-progress
    at: 2026-11-27T22:34:34Z
    by: agent-S-0850
stream: S-0850
tags: []
---
# T-3611 Run the analyzer and write its report

## Work

Run the role `analyzer` with this prompt:

> Analyze this project with the focus "": bottlenecks, intent, or risk, and all three when the focus is empty. Write one report under design/analysis, and file or bump an issue for each actionable finding, linked from the report.

## Done when

The role's run ends with its work done.

## Notes
