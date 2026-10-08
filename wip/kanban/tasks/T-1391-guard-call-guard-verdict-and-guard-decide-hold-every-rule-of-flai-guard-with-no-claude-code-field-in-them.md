---
id: T-1391
type: task
nature: improvement
title: guard.Call, guard.Verdict, and guard.Decide hold every rule of flai guard with no Claude Code field in them
status: backlog
parent: S-0353
owner: alex
created: 2026-10-08T08:46:53Z
updated: 2026-10-08T08:46:53Z
transitions: []
stream: S-0353
tags: [cli]
touches: [flai/internal/guard/call.go, flai/internal/guard/call_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/shared.go, flai/internal/guard/shared_test.go]
---
# T-1391 guard.Call, guard.Verdict, and guard.Decide hold every rule of flai guard with no Claude Code field in them

## Work

- `call.go`: `Call` (`Session`, `Subagent`, `Role`, `Story`, `Kind` of `Shell`, `FileWrite`, `FlaiTool`, `Other`, `Command`, `Path`, `Tool`, `Args`), `Verdict` (`Allow`, `Deny`, `Ask`, with `Reason`), and `Decide(c Call, running Running) Verdict`.
- Move each rule out of `guard.go` into `Decide` or the helpers it calls, keyed on `Kind` and the bare tool name: the sub-agent's flai reads and refused writes, the `git` and `flai` command splitting in `shared.go`, the `.claude/` write while auto-approve is off (ADR-0102), `wait_for_events` while a sub-agent runs (ADR-0092), and each strategic role's permissions. Refusal texts stay word for word.
- `guard.go` keeps today's entry point by building a `Call` from the hook input inline, so every test in `guard_test.go`, `orchestrate_test.go`, and `thread_test.go` passes unchanged; T-1392's reader replaces that inline mapping.
- `call_test.go` drives `Decide` directly with neutral calls for one rule of each caller.

First layer.

## Done when

- `grep` finds no `tool_name`, `agent_id`, `mcp__`, `Bash`, `Edit`, or `Write` literal in `call.go`.
- `flai test flai/internal/guard/` passes with no expectation changed.

## Notes

Layer 1 of S-0353.
