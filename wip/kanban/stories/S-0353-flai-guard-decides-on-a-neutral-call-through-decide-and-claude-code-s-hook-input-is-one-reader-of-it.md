---
id: S-0353
type: story
nature: improvement
title: flai guard decides on a neutral call through Decide, and Claude Code's hook input is one reader of it
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:46:31Z
updated: 2026-10-08T08:52:06Z
transitions: []
tags: [cli]
topics: [agents]
touches: [flai/internal/guard/call.go, flai/internal/guard/call_test.go, flai/internal/guard/claudecode.go, flai/internal/guard/claudecode_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/shared.go, flai/internal/guard/shared_test.go, flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, design/system/flai-cli.md]
after: [S-0351]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
finalized:
  by: alex
  at: 2026-10-08T08:52:06Z
---
# S-0353 flai guard decides on a neutral call through Decide, and Claude Code's hook input is one reader of it

## Goal

`flai guard` reads Claude Code's hook JSON and keys every rule on Claude Code's tool names, `mcp__flai__` prefixes, and `agent_id` ([ADR-0060](../../../design/adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md)). [ADR-0130](../../../design/adrs/0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md) splits it: one decision, `Decide(Call) Verdict`, on a neutral call that says who calls (the session's agent, a sub-agent, a strategic role), what (a shell line, a file write, a flai tool), and with which arguments; and one input reader per harness that maps its hook input onto that call. This story builds the decision and Claude Code's reader, with every rule and every refusal text as they are.

## Acceptance criteria

- [ ] `guard.Call` carries the session, whether a sub-agent made the call, the role from `FLAI_ROLE`, the story from `FLAI_STORY`, the kind (`Shell`, `FileWrite`, `FlaiTool`, `Other`), the shell line, the path, the flai tool's name without any harness prefix, and its arguments; `guard.Verdict` is allow, deny with a reason, or ask with a reason.
- [ ] `guard.Decide(Call, running)` holds every rule `guard.go` applies today, for a story's agent, its sub-agents, the planner, the orchestrator, and the analyzer, and no rule reads a Claude Code field or tool name.
- [ ] A Claude Code reader maps the hook's `hook_event_name`, `session_id`, `tool_name`, `tool_input`, `agent_id`, and `agent_type` onto a `Call`, including `SubagentStart` and `SubagentStop` onto the running sub-agent record, and maps a verdict back onto exit 2 and the reason on standard error; input it cannot read still passes, as today.
- [ ] `flai guard` takes `--harness <name>`, default `claude-code`, and refuses a harness with no reader, naming the ones it has; every existing test in `flai/internal/guard` and `flai/cmd` passes with its expectations unchanged.
- [ ] `design/system/flai-cli.md` describes `Call`, `Decide`, and the reader per harness.

## Tasks

Drafted by the planner; see the children.
- T-1391 guard.Call, guard.Verdict, and guard.Decide hold every rule of flai guard with no Claude Code field in them
- T-1392 A Claude Code reader maps the hook's input onto a Call and the verdict onto its exit, and flai guard takes --harness
- T-1393 flai-cli.md describes the guard's Call, Decide, and one reader per harness

## Notes

- No rule returns `ask` yet. Today Claude Code itself sends a protected write to `permission_prompt`; the next story makes that the guard's hold-and-ask and fills in the protected list.
- No change to `.claude/settings.json` or the template: the hook command stays `flai guard`, and `--harness` defaults to `claude-code`.
