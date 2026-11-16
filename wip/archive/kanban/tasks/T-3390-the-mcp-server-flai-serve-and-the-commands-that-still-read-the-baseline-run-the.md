---
id: T-3390
type: task
nature: feature
title: The MCP server, flai serve, and the commands that still read the baseline run the project's workflow
status: done
parent: S-0755
owner: alex
created: 2026-11-15T21:53:39Z
updated: 2026-11-15T22:58:55Z
transitions:
  - to: ready
    at: 2026-11-15T22:33:57Z
    by: agent-S-0755
  - to: in-progress
    at: 2026-11-15T22:33:57Z
    by: agent-S-0755
  - to: done
    at: 2026-11-15T22:58:55Z
    by: agent-S-0755
stream: S-0755
tags: []
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/workflow_test.go, flai/internal/serve/flow.go, flai/internal/serve/flow_test.go, flai/internal/serve/kluster.go, flai/internal/serve/replan.go, flai/internal/serve/role_events.go, flai/internal/inbox/inbox.go, flai/internal/planning/cod.go, flai/internal/itemedit/forecasts.go, flai/internal/preview/accept.go, flai/cmd/task_done.go, flai/cmd/accept_overlap.go, flai/cmd/accept.go, flai/cmd/accept_overlap_test.go, flai/cmd/cod.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/cmd/project.go, flai/internal/itemedit/itemedit.go, flai/internal/mcpserver/folder.go, flai/internal/planning/forecast.go]
after: [T-3253, T-3259, T-3260]
usage:
  source: log
  seconds: 1498
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 187
      output: 63828
      cache_read: 13162515
      cache_write: 386209
      cost: 6.1827
---
# T-3390 The MCP server, flai serve, and the commands that still read the baseline run the project's workflow

## Work

Waits for T-3253, whose load and `Repo.Flow()` it reads, and for T-3259 and T-3260, which run their tests in packages this task's edits make `cmd` import half done. ADR-0307 §2 and §3.

- The MCP server reads the file again at each tool call when it changed, and refuses the call with the finding while it is refused.
- `flai serve`: `openRepo` in `serve/flow.go` sets `repo.Workflow` only when the options give one, so the loaded file stands; a refused open is logged once per distinct finding.
- Each live path that still reads `workflow.Baseline()` or a package `baselineFlow` where a repository is at hand reads its `Flow()`: `inbox.Describe`, `serve/kluster.go`'s direct read, `planning.Replay` from `serve/replan.go`, `metrics/costorder.go`'s puller, `planning/cod.go`'s forecaster, `itemedit/forecasts.go`'s `editable`, `preview.ArchivedWith`, `cmd/task_done.go`'s `printSyncChecks`, `cmd/accept_overlap.go`'s `Covering`, and the guard's and the roles' callers that leave their workflow nil.

## Done when

A test shows the MCP server running a project's file edited while it runs, and refusing a call on a broken one; a test shows `openRepo` keeping the loaded file; no non-test caller of a baseline wrapper is left where a repository is at hand; `flai test` passes on the changed packages.

## Notes
