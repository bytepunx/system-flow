---
id: T-1070
type: task
nature: improvement
title: The MCP tool test runs the tiers for paths in a story's worktree and answers pass or the first findings
status: done
parent: S-0273
owner: alex
created: 2026-10-06T22:52:05Z
updated: 2026-10-07T00:20:50Z
transitions:
  - to: ready
    at: 2026-10-07T00:11:06Z
    by: agent-S-0273
  - to: in-progress
    at: 2026-10-07T00:11:06Z
    by: agent-S-0273
  - to: done
    at: 2026-10-07T00:20:50Z
    by: agent-S-0273
stream: S-0273
tags: [mcp, testing]
touches: [flai/internal/mcpserver/folder.go, flai/internal/mcpserver/test.go, flai/internal/mcpserver/test_test.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder_test.go]
after: [T-1067, T-1068]
usage:
  source: log
  seconds: 584
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 50
      output: 21936
      cache_read: 2941926
      cache_write: 70807
      cost: 1.4364
---
# T-1070 The MCP tool test runs the tiers for paths in a story's worktree and answers pass or the first findings

## Work

Criterion 2, the MCP half. This task waits for T-1067 and T-1068: it reads the manifest's tiers and runs them through the verify package. It shares no path with the command task, so the two run together.

- Add the tool `test` in a new `flai/internal/mcpserver/test.go`, the way `criteria.go` adds `criteria_tick`, and register it in `folder.go`'s `addProjectTools`.
- Arguments:
  - `story`: run in that story's worktree; with none, run in the project's main checkout
  - `paths`: paths or packages
  - `all`
  - `max`
- With no paths, it takes what the story's branch changed against the main branch.
- Answer the verify package's result, the same as `flai test --json`, so the two never disagree.
- Write its description as the other tools' are: what it runs, what it answers, and that it replaces running `go test`, vitest, golangci-lint, or gofmt by hand between tasks.
- Respect the call's context, so a cancelled call stops the tier's processes.

## Done when

- [ ] The tool is listed. Called on a fixture project with a failing test, it answers that test's findings; with passing paths, it answers pass (tests in `test_test.go`).
- [ ] Called with `story`, it runs in that story's worktree.
- [ ] `go test -race ./internal/mcpserver/...` and golangci-lint pass.

## Notes

Drafted by the planner for S-0273.
