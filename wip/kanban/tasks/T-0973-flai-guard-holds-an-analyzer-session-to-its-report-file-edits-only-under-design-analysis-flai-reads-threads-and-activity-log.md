---
id: T-0973
type: task
nature: feature
title: "flai guard holds an analyzer session to its report: file edits only under design/analysis, flai reads, threads, and activity_log"
status: in-progress
parent: S-0223
owner: alex
created: 2026-10-05T05:47:47Z
updated: 2026-10-06T20:28:08Z
transitions:
  - to: ready
    at: 2026-10-06T20:28:07Z
    by: agent-S-0223
  - to: in-progress
    at: 2026-10-06T20:28:08Z
    by: agent-S-0223
stream: S-0223
tags: [flai]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go]
after: [T-0960]
---
# T-0973 flai guard holds an analyzer session to its report: file edits only under design/analysis, flai reads, threads, and activity_log

## Work

Add the role `analyze` to `flai guard` (`flai/internal/guard/guard.go`, beside `RolePlan`, `MCPPlans`, `cliPlans`, and `planning`), read from `FLAI_ROLE` as `plan` is (`flai/cmd/guard.go`):

- Edit, Write, and NotebookEdit only on a file under the manifest's design folder's `analysis/` (`analysis.Dir()`), refused elsewhere with the folder it may write.
- The MCP reads, `activity_log`, `inbox`, `thread_open`, `thread_reply`, and `wait_for_events`; every item write (`item_new`, `item_edit`, `item_move`) refused, saying the analyzer authors no stories.
- The flai CLI reads, `flai stats` among them, and no flai write in a shell command.

The explorer it hands search to keeps the sub-agent rules it has. Issues are S-0224's: this task allows none.

It waits for T-0960, whose `analysis.Dir()` it reads.

## Done when

- a test lets an analyzer session write a report under `design/analysis/` and refuses one outside it, naming the folder
- a test refuses `item_new`, `item_edit`, and `item_move` to it, and allows `activity_log`, `thread_open`, and `flai stats --json`
- the planner's rules are unchanged, in their tests
- `go test ./internal/guard/ ./cmd/` passes

## Notes
