---
id: T-2275
type: task
nature: improvement
title: The wait_for_events and inbox tool descriptions name the agent kind and its fields, and a test per to reads an agent event back through each tool
status: backlog
parent: S-0532
owner: alex
created: 2026-10-28T21:09:18Z
updated: 2026-10-28T21:09:18Z
transitions: []
stream: S-0532
tags: [flai, mcp, events]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go]
after: [T-2271]
---
# T-2275 The wait_for_events and inbox tool descriptions name the agent kind and its fields, and a test per to reads an agent event back through each tool

## Work

Name the event in the MCP descriptions:

- `wait_for_events` in `server.go` and its folder server twin in `folder.go`: an event of kind `agent` for each start, end, failure, or stop of a role's run or a story's agent, with the item or role, `to` (`started`, `ended`, `failed`, `stopped`), and `why`, reported once like the others under the 50-event cap;
- `inboxDescription` in `folder.go`, for `changes`, the same;
- the server's instructions in `server.go`, one sentence saying what an `agent` event is for: seeing an agent fail or restart without reading the serve journal.

Then read the events back through the tools: in `server_test.go`, append an event of each `to` to the log and read it through `wait_for_events`, once held and once already behind the cursor, and through `inbox`, each reported once; and one test in `folder_test.go` that an event in one project of a folder is reported with that project.

This task waits for T-2271 alone, and runs beside the serve tasks: its tests write the log directly.

## Done when

- [ ] Both tools' descriptions, in the single-project and the folder server, name the kind, `to`'s four values, and the fields
- [ ] A test per `to` reads an `agent` event back through `wait_for_events` and through `inbox`, once each
- [ ] The folder server reports an `agent` event with its project
- [ ] `flai test` passes on the changed paths

## Notes

Planned by planner-S-0532.
