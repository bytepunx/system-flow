---
id: T-1771
type: task
nature: feature
title: flai check reports formulas.file on each formula that could not run
status: backlog
parent: S-0421
owner: alex
created: 2026-10-16T05:47:23Z
updated: 2026-10-16T05:47:23Z
transitions: []
stream: S-0421
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/formulas_test.go]
after: [T-1769]
---
# T-1771 flai check reports formulas.file on each formula that could not run

## Work

Call `formulas.Validate` from `flai check` and report each problem as `formulas.file`, on its field's line. Add a test per problem.

Waits for the loader. Runs beside the command task, since they share no path.

## Done when

`go test ./internal/check/...` passes.

## Notes
