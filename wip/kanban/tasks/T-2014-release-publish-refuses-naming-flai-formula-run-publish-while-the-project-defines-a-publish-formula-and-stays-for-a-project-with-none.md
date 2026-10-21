---
id: T-2014
type: task
nature: improvement
title: release_publish refuses, naming flai formula run publish, while the project defines a publish formula, and stays for a project with none
status: backlog
parent: S-0423
owner: alex
created: 2026-10-21T13:10:20Z
updated: 2026-10-21T13:27:19Z
transitions: []
stream: S-0423
tags: [cli]
touches: [flai/internal/mcpserver/orchestrate.go, flai/internal/mcpserver/orchestrate_test.go]
after: [T-2015]
---
# T-2014 release_publish refuses, naming flai formula run publish, while the project defines a publish formula, and stays for a project with none

## Work

In `flai/internal/mcpserver/orchestrate.go`, `releasePublish` reads the project's formulas with `flai/internal/formulas` before `mayPublish`. When one named `publish` exists, it refuses, changing nothing, saying that publishing is the `publish` formula and naming `flai formula run publish --ready`. With no such formula it publishes as today, so a project that has not taken the baseline formula keeps the tool.

Waits for the ADR (T-2015), which records the decision it enforces. Runs beside T-1780, since they share no path.

## Done when

`go test ./internal/mcpserver/...` passes, with a test for each case. The three existing `release_publish` tests are unchanged.

## Notes

Criterion 3: the orchestrator reaches the remote through the formula and no other way, and `release_publish` stays only where there is no `publish` formula.
