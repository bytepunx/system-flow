---
id: T-0890
type: task
nature: feature
title: "flai promote --drafts lists each draft story with what it lacks to be finalized: criteria, touches, forecast, and value"
status: backlog
parent: S-0219
owner: alex
created: 2026-10-05T04:46:38Z
updated: 2026-10-05T04:46:38Z
transitions: []
stream: S-0219
tags: [flai]
touches: [flai/internal/workitem/finalizable.go, flai/internal/workitem/finalizable_test.go, flai/cmd/promote.go, flai/cmd/promote_drafts_test.go]
---
# T-0890 flai promote --drafts lists each draft story with what it lacks to be finalized: criteria, touches, forecast, and value

## Work

The second criterion, with `finalize_drafts`, finalizes a draft whose criteria, touches, forecast, and value are complete and consistent. Completeness is arithmetic and belongs in flai; consistency, whether the criteria, touches, and forecast describe the same work, stays the orchestrator's judgement.

- Add `flai/internal/workitem/finalizable.go`: a draft story is complete when it has every section of the template, a goal, at least one acceptance criterion as a checkbox, at least one touch, a forecast `duration` and `delivery`, a cost of delay `value`, and an open parent epic or none. It returns each thing it lacks. A story that is not a draft is not listed.
- Add `--drafts [--json]` to S-0217's `flai promote` (`flai/cmd/promote.go`): it lists each draft in the backlog, complete or with what it lacks, in the project's policy order, and writes nothing. The orchestrator finalizes a complete one it judges consistent and opens a thread on an incomplete one quoting what this lists.

This task waits for no other task of this story; it needs S-0217's `flai promote`, which S-0219 waits for through S-0218. It runs with the other first-layer tasks, whose paths it does not share.

## Done when

- fixture tests pin a complete draft, and a draft lacking each of: a criterion, touches, a forecast, and a value
- `flai promote --drafts` changes no file, and `--json` gives each draft's ID, whether it is complete, and what it lacks
- `go test ./internal/workitem/ ./cmd/` passes

## Notes
