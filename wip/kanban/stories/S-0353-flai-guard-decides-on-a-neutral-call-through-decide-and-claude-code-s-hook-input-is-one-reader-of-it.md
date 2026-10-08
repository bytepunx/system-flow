---
id: S-0353
type: story
nature: improvement
title: flai guard decides on a neutral call through Decide, and Claude Code's hook input is one reader of it
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:46:31Z
updated: 2026-10-08T10:29:12Z
transitions: []
tags: [cli]
topics: [agents]
touches: [flai/internal/guard/call.go, flai/internal/guard/call_test.go, flai/internal/guard/claudecode.go, flai/internal/guard/claudecode_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/shared.go, flai/internal/guard/shared_test.go, flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/operators/settings.md, docs/users/flai-reference.md, design/system/flai-cli.md]
after: [S-0351]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 55.31
  by: planner-E-0019
  at: 2026-10-08T09:00:21Z
forecast:
  duration: 50m
  delivery: 2026-10-08T21:35:00Z
  basis: "Its own forecast of 50m; 22nd in the pull order with an in-progress limit of 5, behind S-0232, S-0338, S-0342, S-0337, S-0347, S-0343, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0289, S-0304, S-0305, S-0306, S-0313, S-0349, S-0350, S-0351 and S-0352."
  by: flai
  at: 2026-10-08T10:29:12Z
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
- [ ] `flai guard` takes `--harness <name>`, default `claude-code`; a harness with no reader passes the call with a warning on standard error naming the readers there are, as the guard fails open on input it cannot read (ADR-0060); every existing test in `flai/internal/guard` and `flai/cmd` passes with its expectations unchanged.
- [ ] `design/system/flai-cli.md` describes `Call`, `Decide`, and the reader per harness.

## Tasks

- T-1391 guard.Call, guard.Verdict, and guard.Decide hold every rule of flai guard with no Claude Code field in them
- T-1392 A Claude Code reader maps the hook's input onto a Call and the verdict onto its exit, and flai guard takes --harness
- T-1393 flai-cli.md describes the guard's Call, Decide, and one reader per harness

## Notes

- No rule returns `ask` yet. Today Claude Code itself sends a protected write to `permission_prompt`; the next story makes that the guard's hold-and-ask and fills in the protected list.
- No change to `.claude/settings.json` or the template: the hook command stays `flai guard`, and `--harness` defaults to `claude-code`.

### Planning

Planned by planner-E-0019 on 2026-10-08. Every touch is a file; no folder touch is kept; `call.go`, `claudecode.go`, and their tests are new. Three layers, one task each: `Decide` (T-1391), the reader (T-1392), which replaces T-1391's inline mapping in `guard.go`, and the document (T-1393).

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/guard/call.go`, `call_test.go` | design | ADR-0130's `Call`, `Verdict`, and `Decide` (T-1391) |
| `flai/internal/guard/guard.go`, `guard_test.go`, `shared.go`, `shared_test.go` | layout | The rules and the shell-line splitting live here (T-1391) |
| `flai/internal/guard/claudecode.go`, `claudecode_test.go`, `subagents.go`, `subagents_test.go` | design, layout | Claude Code's reader, and the running sub-agent record its start and stop feed (T-1392) |
| `flai/cmd/guard.go`, `guard_test.go` | layout | `--harness` (T-1392) |
| `docs/operators/settings.md`, `docs/users/flai-reference.md` | layout | `--harness` is a new flag: `make flai-reference` regenerates both (T-1392) |
| `design/system/flai-cli.md` | design | ADR-0130's consequences (T-1393) |

`touches suggest` listed the user and operator guides and the dashboard's design. None is taken beyond the generated two: nothing a user does changes.

Forecast: 50m. `flai forecast` gave 26m, 78 s per unit over 51 improvement stories, size 20. It is raised because `guard.go` is 1,550 lines with a 1,289-line test, and every rule moves with its refusal text unchanged.

Cost of delay: 55.31 USD a week, as `flai cod` works it out: 50m of the 7h32m forecast over E-0019's 12 stories, of the operator's 500 USD a week penalty. It stands.
