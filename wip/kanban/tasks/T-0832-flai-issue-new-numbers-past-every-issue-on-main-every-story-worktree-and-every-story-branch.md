---
id: T-0832
type: task
nature: remediation
title: flai issue new numbers past every issue on main, every story worktree, and every story branch
status: backlog
parent: S-0252
owner: alex
created: 2026-10-04T23:27:36Z
updated: 2026-10-04T23:27:36Z
transitions: []
stream: S-0252
tags: [flai, issues]
touches: [flai/internal/issues, flai/cmd/issue.go, flai/cmd/issue_test.go, flai/internal/mcpserver/issues_test.go, docs/users/flai-reference.md]
after: [T-0831]
---
# T-0832 flai issue new numbers past every issue on main, every story worktree, and every story branch

## Work

Have `issues.NextID` take the highest I-number among the names T-0831's helper returns for the issues folder, and the current checkout's own names, and add one. Pass it a runner, for instance through `NewOptions`, defaulting to `execx.System`, so that `issues.New`'s callers and the MCP tests change as little as they must. Say in `flai issue new`'s help that the number is past every open story's issues, and regenerate the reference with `make flai-reference`.

It waits for T-0831, whose helper it calls.

## Done when

- A test in `flai/cmd/issue_test.go` reproduces I-0065. Two story branches each hold an issue: the first has I-0002 committed, and the second has I-0003 uncommitted in its worktree. `flai issue new` run in a third story's worktree gives I-0004. The test fails without the change.
- A project with no story branches still numbers from its own `design/issues`, as before.
- `docs/users/flai-reference.md` matches the help text, and `scripts/flai-test.sh` passes.

## Notes
