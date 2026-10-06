---
id: T-1035
type: task
nature: remediation
title: flai guard refuses a sub-agent's write under .claude/ at once while auto-approve is off, so its layer never waits on a permission thread
status: backlog
parent: S-0299
owner: alex
created: 2026-10-06T21:08:01Z
updated: 2026-10-06T21:08:01Z
transitions: []
stream: S-0299
tags: [flai, guard]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, flai/internal/hostapi/writes.go, docs/users/flai-reference.md]
---
# T-1035 flai guard refuses a sub-agent's write under .claude/ at once while auto-approve is off, so its layer never waits on a permission thread

## Work

I-0093's cause: a task sub-agent of S-0223 called `Write` and `Edit` on files under `template/root/.claude/`. Claude Code passed each call to flai's `permission_prompt`, which opened a thread and held the call until the owner answered. The owner was away, so the whole layer waited. The guard already sees every `Edit`, `Write`, and `NotebookEdit` (the `.claude/settings.json` matcher `Edit|Write|NotebookEdit`) and knows a sub-agent's call by its `agent_id`. That makes it the place to refuse the call before it can be held.

- In `Guard.Decide` (`flai/internal/guard/guard.go`), when a sub-agent makes the call (`AgentID` set) and it is an `Edit`, `Write`, `MultiEdit`, or `NotebookEdit` of a file with a `.claude` folder in its path, refuse it unless the project's `auto-approve` host action is on. The refusal names the sub-agent and the path. It says that a file in a `.claude/` folder is written only with the operator's approval on a thread, which would hold this call and the layer with it until they answer. It tells the sub-agent to put the file's whole new content in its final message for the story's agent to write. Leave the story's agent's own calls, and every write outside a `.claude/` folder, to `permission_prompt` and the session's permissions as now.
- Give `Guard` what it needs to know whether auto-approve is on. Set it in `flai/cmd/guard.go` from the host configuration, the same way `mcpAutoApprove` in `flai/cmd/mcp.go` reads it (`hostapi.ActionAutoApprove`). When the configuration cannot be read, count auto-approve as off, as `permission_prompt` does.
- Name the rule in `flai guard`'s help and regenerate `docs/users/flai-reference.md` from it. Say in `ActionAutoApprove`'s description in `flai/internal/hostapi/writes.go` that turning it on also lets a story's sub-agents make these writes.
- This task waits for no other: it shares no path with the start-prompt task and runs beside it.

## Done when

- A test in `flai/internal/guard/guard_test.go` reproduces I-0093. It refuses a sub-agent's `Write` of `<worktree>/template/root/.claude/agents/analyzer.md` and its `Edit` of `<worktree>/.claude/settings.json` while auto-approve is off, naming the path and the final-message route. It allows both with auto-approve on. It allows the story's agent's own write and a sub-agent's write outside a `.claude/` folder.
- A test in `flai/cmd/guard_test.go` feeds `flai guard` such a hook input from a sub-agent and gets exit 2 with the refusal; with auto-approve enabled for the project it gets exit 0.
- `flai guard --help` and `docs/users/flai-reference.md` name the rule, and the reference check passes.
- `go test ./internal/guard/... ./cmd/... ./internal/hostapi/...` passes in `flai/`.

## Notes

Drafted by the planner for S-0299. The project's `.claude/settings.json` needs no change: its matcher already sends these calls to the guard.
