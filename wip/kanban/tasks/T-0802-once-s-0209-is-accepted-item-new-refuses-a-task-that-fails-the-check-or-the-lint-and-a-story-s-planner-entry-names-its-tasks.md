---
id: T-0802
type: task
nature: feature
title: Once S-0209 is accepted, item_new refuses a task that fails the check or the lint, and a story's planner entry names its tasks
status: backlog
parent: S-0255
owner: alex
created: 2026-10-04T04:05:35Z
updated: 2026-10-04T04:05:35Z
transitions: []
stream: S-0255
tags: []
touches: [flai/internal/mcpserver/items_write_test.go, flai/internal/serve/plan_test.go]
after: [T-0799]
---

# T-0802 Once S-0209 is accepted, item_new refuses a task that fails the check or the lint, and a story's planner entry names its tasks

## Work

S-0209's T-0792 makes the MCP tool `item_new` check an item with a body through `itemnew.Create`, and its T-0793 makes a planner run's activity entry name the planned item and the items under it created or changed during the run. Both are what S-0255's criteria 3 and 4 need for a story's tasks (TH-0099). Once S-0209 is accepted and this branch synced, add task cases to their tests: `item_new` type task with a body that brings a lint finding is refused and leaves nothing; a planner run on a story names the tasks created and edited during it, and not one left alone. If either does not cover tasks, make it. Waits for T-0799, whose `itemnew` tests these extend to the MCP path, and for S-0209's acceptance.

## Done when

- the task cases pass and fail without the behaviour they pin
- `go test ./internal/mcpserver/ ./internal/serve/` passes

## Notes
