---
id: T-1392
type: task
nature: improvement
title: A Claude Code reader maps the hook's input onto a Call and the verdict onto its exit, and flai guard takes --harness
status: backlog
parent: S-0353
owner: alex
created: 2026-10-08T08:47:03Z
updated: 2026-10-08T08:54:18Z
transitions: []
stream: S-0353
tags: [cli]
touches: [flai/internal/guard/claudecode.go, flai/internal/guard/claudecode_test.go, flai/internal/guard/guard.go, flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/operators/settings.md, docs/users/flai-reference.md]
after: [T-1391]
---
# T-1392 A Claude Code reader maps the hook's input onto a Call and the verdict onto its exit, and flai guard takes --harness

## Work

- `claudecode.go`: a `Reader` for Claude Code's hook JSON. It maps `tool_name` (`Bash`, `Edit`, `MultiEdit`, `Write`, `NotebookEdit`, `mcp__flai__*`) onto `Kind` and the bare tool, `tool_input`'s `command`, `file_path`, `notebook_path`, and flai tool fields onto the call, and `agent_id` onto `Subagent`. `SubagentStart` and `SubagentStop` become record updates. A verdict becomes exit 0, or exit 2 with the reason on standard error.
- `subagents.go` records running sub-agents from the reader's start and stop, not from Claude Code's event names.
- `guard.go`'s inline mapping from T-1391 is replaced by the reader. Input the reader cannot parse passes, as today.
- `flai guard --harness <name>`, default `claude-code`; an unknown name exits 0 with a warning on standard error naming the readers, since the guard fails open (ADR-0060), and `flai guard --help` says so.

It waits for T-1391, whose `Call` and `Decide` it feeds, and changes `guard.go` after it.

## Done when

- `claudecode_test.go` covers each tool kind, a sub-agent's call, a start and a stop, and unreadable input; `guard_test.go` in `flai/cmd` covers an unknown `--harness`.
- `flai test flai/internal/guard/ flai/cmd/` passes with no expectation changed.

## Notes

Layer 2 of S-0353.
