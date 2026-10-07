---
id: T-1083
type: task
nature: improvement
title: MCP issue_new and issue_bump take commit, and issue_close and adr_new are new tools
status: in-progress
parent: S-0275
owner: alex
created: 2026-10-06T22:52:46Z
updated: 2026-10-07T08:21:02Z
transitions:
  - to: ready
    at: 2026-10-07T08:21:02Z
    by: agent-S-0275
  - to: in-progress
    at: 2026-10-07T08:21:02Z
    by: agent-S-0275
stream: S-0275
tags: [mcp]
touches: [flai/internal/mcpserver/issues.go, flai/internal/mcpserver/issues_test.go, flai/internal/mcpserver/adr.go, flai/internal/mcpserver/adr_test.go, flai/internal/mcpserver/folder.go]
after: [T-1072, T-1074]
---
# T-1083 MCP issue_new and issue_bump take commit, and issue_close and adr_new are new tools

## Work

Criterion 2, over MCP. It waits for T-1072, the commit-and-widen helper, and for T-1074, the issue paths. It shares no path with T-1077, so the two run together: the MCP tools call the packages, not the CLI.

- Give `issue_new` and `issue_bump` a `commit` input that does what `--commit` does: write in the story's worktree (`issues.RepoFor`), commit the files written on the story branch with the story's prefix, and widen the story's touches. Their output gains the commit and the touches added.
- Add `issue_close` in `issues.go`: an issue ID, a reason, and `commit`.
- Add `adr_new` in the new `adr.go`: the decision sentence, status, supersedes or refines, a body, and `commit`.
- Register both in `folder.go` beside the issue tools. Each tool's description says what it writes and what `commit` does.
- Keep `adr_new` from writing the main checkout's ADR folder while a story is named: it writes in that story's worktree, as the issue tools do.

## Done when

- Tests in `issues_test.go` and `adr_test.go` call each tool with and without `commit` against a temporary story worktree. They assert the files written, the commit on the story branch, and the widened touches.
- The tool list test in `folder_test.go` or `server_test.go`, whichever lists the tools, names the two new tools. Widen this task's touches to that file if it changes.
- `scripts/flai-test.sh` passes.

## Notes

Written by the planner. `flai guard` gives sub-agents `issue list` only, and the analyzer `issue_new` and `issue_bump`. Leave its lists as they are unless a test shows a story's agent is refused the new tools.
