---
id: T-0811
type: task
nature: feature
title: flai promote --candidates lists the backlog stories that could go to ready, by the project's policy, and why each other one cannot
status: backlog
parent: S-0217
owner: alex
created: 2026-10-04T19:09:05Z
updated: 2026-10-04T19:09:17Z
transitions: []
stream: S-0217
tags: [flai]
touches: [flai/internal/workitem/promote.go, flai/internal/workitem/promote_test.go, flai/internal/workitem/hold.go, flai/cmd/promote.go, flai/cmd/promote_test.go, flai/cmd/root.go]
after: [T-0809, T-0810]
---
# T-0811 flai promote --candidates lists the backlog stories that could go to ready, by the project's policy, and why each other one cannot

## Work

Add `flai/internal/workitem/promote.go`. It checks each backlog story, and a story is a candidate when it:

- is not a draft
- meets the definition of ready that `flai move` enforces: a goal, criteria as checkboxes, and an open parent epic (reuse that check, do not copy it)
- would not be held if it were in ready: no overlap, no unfinished `after`, and not without touches. Use `Holds.Of` (`hold.go`, line 134), changing it only if it assumes ready.
- has a forecast `duration` and a cost of delay `value`

Each other story is listed with every reason it is not a candidate.

Order the candidates by `orchestration.policy` with T-0810's function. Add the command `flai promote --candidates [--limit N] [--json]` in `flai/cmd/promote.go`, registered in `root.go`. It lists the candidates in order with their figures, then the others with their reasons, and writes nothing.

This task waits for T-0809, for the policy setting, and for T-0810, for the ordering function. It runs with T-0812, whose paths it does not share.

## Done when

- fixture tests pin a candidate, and a story refused for each reason: a draft, not ready, held, no forecast, and no value
- the candidates follow the manifest's policy, and `--limit` caps them
- `flai promote --candidates` changes no file
- `go test ./internal/workitem/ ./cmd/` passes

## Notes
