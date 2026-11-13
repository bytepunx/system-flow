---
id: T-3224
type: task
nature: feature
title: item_move refuses a transition its caller may not make, naming it and who may, and the MCP server's other tools read a state's kind
status: backlog
parent: S-0746
owner: alex
created: 2026-11-12T16:07:42Z
updated: 2026-11-12T16:10:58Z
transitions: []
stream: S-0746
tags: [cli]
touches: [flai/internal/mcpserver/work.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/installed.go]
---
# T-3224 item_move refuses a transition its caller may not make, naming it and who may, and the MCP server's other tools read a state's kind

## Work

Replace the references to a state's name in these files with reads of the workflow: `item_move` refuses a transition its caller may not make, by the workflow package's mover check that S-0737 adds, naming the transition and who may make it. It keeps refusing a move to an accepted state, and keeps refusing the move to the wait state while the worktree has uncommitted changes, as today. `permission_prompt`, `folder.go`, and `installed.go` tell an open story from a closed one by kind. Test `item_move` on a workflow with a state added between `in-progress` and `review`.

## Done when

These files name no state, and `work_test.go` and `server_test.go` pass on both workflows.

## Notes
