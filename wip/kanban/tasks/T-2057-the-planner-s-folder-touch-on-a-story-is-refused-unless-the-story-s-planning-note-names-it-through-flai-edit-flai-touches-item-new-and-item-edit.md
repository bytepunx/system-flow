---
id: T-2057
type: task
nature: improvement
title: The planner's folder touch on a story is refused unless the story's Planning note names it, through flai edit, flai touches, item_new, and item_edit
status: backlog
parent: S-0468
owner: alex
created: 2026-10-22T21:05:36Z
updated: 2026-10-22T21:05:36Z
transitions: []
stream: S-0468
tags: [cli]
touches: [flai/internal/itemedit/itemedit.go, flai/internal/itemedit/planning_test.go, flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go, flai/cmd/touches.go, flai/cmd/touches_test.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go]
after: [T-2056]
---
# T-2057 The planner's folder touch on a story is refused unless the story's Planning note names it, through flai edit, flai touches, item_new, and item_edit

## Work

Waits for the task rule, whose `workitem.FolderTouch` this rule reads, and shares `itemedit.go` and `cmd/touches.go` with it. Builds the ADR's point 3.

- Who the planner is: `FLAI_ROLE=plan` in the environment (`harness.roleEnv` sets it for the planner's session, `cmd/edit.go` and `mcpserver/server.go` read it already; the stdio MCP server `.mcp.json` starts inherits it), else an author whose name starts with `planner-` (ADR-0143 holds a served agent to its own name). Put the test of it in one place the three paths call.
- The note: the story's Notes section, under a `### Planning` heading, names the folder's path, in backticks or not. Read the body as it will be after the edit, so that one `item_edit` carrying both `body` and `touches` passes, and say in the refusal that the note goes under `### Planning` in the Notes.
- `flai/internal/itemedit/itemedit.go` `Apply`, for a story, by the planner: refuse a folder touch (`FolderTouch` true) that the Planning note does not name, writing nothing, and name the folder, its file count, and the note that would admit it. `flai edit --touches` and MCP `item_edit` come through here; `items_write.go` changes only if the refusal needs to be given as data, as other refusals are.
- `flai/internal/itemnew/itemnew.go` `Create`, for a story the planner drafts with touches (an epic's planner does this through `item_new`): the same rule against the body given.
- `flai/cmd/touches.go`, for a story, by the planner: the same rule before the save.
- Not refused: the operator's and a story's agent's story touches, whoever wrote the body; a file; an absent path; a folder inside `claims.shared`, such as `design/adrs`.

Tests: `itemedit/planning_test.go` (the planner refused without the note, accepted with it, accepted when the body in the same call adds it; another author accepted), `itemnew_test.go` (a planner's draft story), `cmd/touches_test.go` under `FLAI_ROLE=plan`, `items_write_test.go` for `item_new` and `item_edit` by an agent named `planner-S-0001`.

## Done when

- The planner's story folder touch is refused, writing nothing, by `flai edit --touches`, `flai touches`, `item_new`, and `item_edit` unless the story's `### Planning` names the folder, and accepted when it does.
- Another author's story folder touch, and a planner's file, absent path, or shared folder, are accepted as before.
- The tests above pass, `go vet`, `gofmt`, and `golangci-lint` are clean, and `flai test` on the changed paths passes.

## Notes

Layer 3, beside the check rules task; no path in common with it.
