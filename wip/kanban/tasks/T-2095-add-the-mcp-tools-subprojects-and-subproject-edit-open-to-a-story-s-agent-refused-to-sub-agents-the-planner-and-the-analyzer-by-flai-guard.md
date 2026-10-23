---
id: T-2095
type: task
nature: improvement
title: Add the MCP tools subprojects and subproject_edit, open to a story's agent, refused to sub-agents, the planner, and the analyzer by flai guard
status: backlog
parent: S-0491
owner: alex
created: 2026-10-23T20:39:51Z
updated: 2026-10-23T20:39:51Z
transitions: []
stream: S-0491
tags: [cli]
touches: [flai/internal/mcpserver/folder.go, flai/internal/mcpserver/subproject.go, flai/internal/mcpserver/subproject_test.go, flai/internal/mcpserver/server_test.go, flai/internal/guard/guard.go, flai/internal/guard/subproject_test.go, docs/users/flai.md]
after: [T-2093]
---
# T-2095 Add the MCP tools subprojects and subproject_edit, open to a story's agent, refused to sub-agents, the planner, and the analyzer by flai guard

## Work

Add two MCP tools in `flai/internal/mcpserver/subproject.go`, registered from `folder.go` as `addSharedTools` is:

- `subprojects`, a read open to every agent: the `projects[]` entries as `flai subproject list --json` prints them, with `project` only to choose the project when the server serves more than one.
- `subproject_edit`, with `op` (`add`, `set`, `remove`), `name`, and the optional `path`, `kind`, `tags`, and `new_name`, as `settings.subproject` takes them. It runs `flai subproject <op> --json` from T-2093 in the project's folder, so that the folder is created and the manifest written the one way; it commits nothing itself, as `shared_paths_edit` does not: a story's agent commits the change with `flai task done` on its branch, and the orchestrator's change on the main checkout is committed as the task that gives it the permission says. Unlike `shared_paths_edit`, it is not refused to a served session: the sub-projects are the project's to shape, not the operator's list of shared paths.

In `flai/internal/guard/guard.go`: the story's agent may call `subproject_edit` and run `flai subproject add`, `set`, and `remove` with no rule; read-only sub-agents are refused both, as they are `flai components config`, with `subprojects` and `flai subproject list` added to the reads (`MCPReads`, `cliReads`); the planner and the analyzer are refused `subproject_edit` and the writing commands by their existing default. The orchestrator stays refused by its default in this task; the next one gives it a permission. Tests in `flai/internal/guard/subproject_test.go` show each verdict, and `flai/internal/mcpserver/subproject_test.go` follows `shared_test.go` with a fixture, `f.call`, and `StartedByServe`.

Add the two names to the sorted list in `server_test.go` and two rows to the tools table in `docs/users/flai.md` § Serving agents over MCP, beside `shared_paths` and `shared_paths_edit`.

It waits for T-2093, whose command it runs, and runs beside the host API task, which touches other files.

## Done when

- `subprojects` lists the manifest's sub-projects, and `subproject_edit` with `op: add` writes the entry, creates the folder, and commits nothing, shown by `subproject_test.go`.
- `flai guard` lets a story's agent call `subproject_edit` and run `flai subproject add`, refuses an explorer or verifier both, lets every agent call `subprojects` and run `flai subproject list`, and refuses the planner, the analyzer, and the orchestrator `subproject_edit`, each shown by a test.
- `server_test.go`'s tool list and the table in `docs/users/flai.md` name both tools, and `go test ./internal/mcpserver/... ./internal/guard/...` passes.

## Notes

Precedents: `flai/internal/mcpserver/shared.go` and `shared_test.go`; `guard.go` § `subAgent`, which refuses `flai components config`, and § `MCPReads`.
