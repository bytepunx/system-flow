---
id: T-1110
type: task
nature: improvement
title: permission_prompt handles an Edit, Write, MultiEdit, or NotebookEdit of every path Claude Code protects inside an in-progress story's worktree, and still refuses .git
status: in-progress
parent: S-0286
owner: alex
created: 2026-10-06T22:53:35Z
updated: 2026-10-07T01:04:24Z
transitions:
  - to: ready
    at: 2026-10-07T01:04:24Z
    by: agent-S-0286
  - to: in-progress
    at: 2026-10-07T01:04:24Z
    by: agent-S-0286
stream: S-0286
tags: [flai]
touches: [flai/internal/protected/protected.go, flai/internal/protected/protected_test.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go]
after: [T-1106]
---
# T-1110 permission_prompt handles an Edit, Write, MultiEdit, or NotebookEdit of every path Claude Code protects inside an in-progress story's worktree, and still refuses .git

## Work

Put the list of paths Claude Code protects, as the ADR names them, in one place that both `permission_prompt` and the acceptance preview read. The planned place is a new leaf package, `flai/internal/protected`. It has a function that says whether a path, relative to a worktree or a repository root, is protected, and whether it is `.git`. `mcpserver` does not import `preview` today, so neither package should hold the list.

In `flai/internal/mcpserver/permission.go`, replace the `.claude` folder test in the route with the shared one:

- An Edit, Write, MultiEdit, or NotebookEdit of any protected path inside an in-progress story's worktree is handled by the same rules as before: auto-approve, or a thread to the story's owner.
- `.git`, or anything under it, is still refused at once.
- A path in the main checkout, or outside the worktree, is still refused at once.

Update the tool's description and `permissionScope` to say what it now covers. These texts reach Claude Code and the agent.

Waits for T-1106, because the ADR decides the list and what is left out.

## Done when

- `flai/internal/protected` holds the list and its tests: each protected path is matched, and `.git`, its contents, and ordinary paths are not.
- `permission_test.go` covers a write to `.mcp.json` and to another protected path outside `.claude/`, each allowed under auto-approve and asked on a thread without it. It also covers refusals of a path under `.git`, of a protected path in the main checkout, and of a path outside the worktree.
- `scripts/flai-test.sh` passes.

## Notes

Drafted by planner-S-0286. The package name and its place are the plan's proposal; the story's agent may place the list elsewhere if it keeps one list for both readers, and updates the touches with `flai touches`.
