---
id: T-1730
type: task
nature: feature
title: The planner's and the analyzer's guard rules are read from planner.md and analyzer.md through the role path, with their guard tests unchanged
status: backlog
parent: S-0416
owner: alex
created: 2026-10-16T05:44:09Z
updated: 2026-10-16T05:44:09Z
transitions: []
stream: S-0416
tags: [cli, template]
touches: [flai/internal/guard/guard.go, design/roles/planner.md, design/roles/analyzer.md, template/root/design/roles/planner.md, template/root/design/roles/analyzer.md]
after: [T-1729]
---
# T-1730 The planner's and the analyzer's guard rules are read from planner.md and analyzer.md through the role path, with their guard tests unchanged

## Work

Route `plan` and `analyze` through the role path from T-1729. Remove from `guard.go` the fixed lists the role files now state: `cliPlans`, `MCPAnalyzes`, `cliAnalyzes`, `cliAuthors`, and the planner's and analyzer's tool lists. Keep the named conditions in Go: draft-only `item_new`, moves to backlog, and `inReports`. Correct the two role files, here and in the template, wherever the tests show they differ from today's rules. The refusal texts stay word for word.

Waits for T-1729. T-1731 waits for this task, since both change `guard.go`.

## Done when

`guard_test.go`, `call_test.go`, `writes_test.go`, `authors_test.go`, and `flai/cmd/guard_test.go` pass unchanged.

## Notes
