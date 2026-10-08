---
id: S-0354
type: story
nature: improvement
title: Asking the owner before a protected write is one hold-and-ask that the guard's ask verdict and permission_prompt share, and the protected list is flai's, kept per harness
status: backlog
parent: E-0019
owner: alex
created: 2026-10-08T08:47:53Z
updated: 2026-10-08T09:08:56Z
transitions: []
tags: [cli]
topics: [agents]
touches: [flai/internal/ask/ask.go, flai/internal/ask/ask_test.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/protected/protected.go, flai/internal/protected/protected_test.go, flai/internal/guard/call.go, flai/internal/guard/call_test.go, flai/internal/guard/claudecode.go, flai/internal/guard/claudecode_test.go, docs/users/flai.md, docs/operators/settings.md, design/system/flai-cli.md]
after: [S-0353]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 38.72
  by: planner-E-0019
  at: 2026-10-08T09:00:22Z
forecast:
  duration: 35m
  delivery: 2026-10-08T22:09:00Z
  basis: "Its own forecast of 35m; 28th in the pull order with an in-progress limit of 5, behind S-0232, S-0322, S-0341, S-0344, S-0345, S-0338, S-0346, S-0342, S-0337, S-0334, S-0343, S-0297, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0289, S-0304, S-0305, S-0306, S-0313, S-0347, S-0348, S-0349, S-0350, S-0351, S-0352 and S-0353."
  by: flai
  at: 2026-10-08T09:08:56Z
finalized:
  by: alex
  at: 2026-10-08T08:52:09Z
---
# S-0354 Asking the owner before a protected write is one hold-and-ask that the guard's ask verdict and permission_prompt share, and the protected list is flai's, kept per harness

## Goal

`permission_prompt` holds a protected write, asks the story's owner or the project's owner on a thread, and answers within four minutes ([ADR-0086](../../../design/adrs/0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md), [ADR-0097](../../../design/adrs/0097-permission-prompt-takes-an-answer-from-the-story-s-owner-or-the-project-s-owner.md), [ADR-0124](../../../design/adrs/0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md)). All of it lives in Claude Code's MCP tool, and the protected paths are Claude Code's list. [ADR-0130](../../../design/adrs/0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md) makes asking the guard's hold-and-ask, with `permission_prompt` as Claude Code's form of it, and the protected list flai's policy with each harness's own files in it. This story moves the asking into one package that both call, gives `Decide` its `ask` verdict, and keys the protected list by harness, with what a story's agent sees unchanged.

## Acceptance criteria

- [ ] A package `flai/internal/ask` holds the hold-and-ask: the question on a thread on the story, the answerers (the story's owner and the project's owner), the bound passed by the caller, the refusal that leaves the thread open, and the answer taken when the same request comes again; `permission_prompt` maps Claude Code's input and output onto it with the four-minute bound, and every test in `permission_test.go` passes unchanged in its expectations.
- [ ] `guard.Decide` returns `ask` for a write by a story's own agent to a protected path inside its worktree while auto-approve is off. Claude Code's reader passes it on to Claude Code, which calls `permission_prompt` as today. `guard.Hold`, for a reader that holds the call itself, asks through `flai/internal/ask` with the harness's bound, and a test drives it with a fake reader.
- [ ] `flai/internal/protected` keeps its paths per harness, Claude Code's as today, with `Path` and `Changed` answering for the union, so that the acceptance preview of ADR-0106 is unchanged; a test pins Claude Code's list.
- [ ] `docs/users/flai.md` renames "Writes to paths Claude Code protects" to say flai's protected paths and describes the hold-and-ask; `design/system/flai-cli.md` describes the package, the verdict, and the list per harness.

## Tasks

- T-1394 The hold-and-ask moves out of permission_prompt into flai/internal/ask, with the bound passed by the caller
- T-1395 The protected list keeps its paths per harness, Claude Code's as today, and Path and Changed answer for the union
- T-1396 guard.Decide returns ask for a story agent's protected write, Claude Code's reader passes it on, and guard.Hold asks for a reader that holds
- T-1397 The user and design documents describe flai's protected paths and the hold-and-ask

## Notes

- No other harness is adopted in E-0019, so the list holds Claude Code's paths only. `.codex/`, `.opencode/`, `.agents/plugins/`, and `AGENTS.md` come with their harness's adapter, as ADR-0130 has it, rather than now: adding them here would leave a story changing `AGENTS.md` to the operator's acceptance alone.
- The `claude --version` check in `claudecheck.go` stays Claude Code's.

### Planning

Planned by planner-E-0019 on 2026-10-08. Every touch is a file; no folder touch is kept; `flai/internal/ask` is a new package of two files. Three layers: T-1394 and T-1395 together, in different packages; then T-1396, which calls both; then T-1397.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/ask/ask.go`, `ask_test.go` | design | ADR-0130's hold-and-ask, out of `permission.go` (T-1394) |
| `flai/internal/mcpserver/permission.go`, `permission_test.go` | layout | Today's `askOperator`, `awaitAnswer`, and the rest (T-1394) |
| `flai/internal/protected/protected.go`, `protected_test.go` | design | The list per harness (T-1395) |
| `flai/internal/guard/call.go`, `call_test.go`, `claudecode.go`, `claudecode_test.go` | design | S-0353's files: the `ask` verdict and `Hold` (T-1396) |
| `docs/users/flai.md`, `docs/operators/settings.md` | layout | The renamed section, and the one link to its anchor outside it (T-1397) |
| `design/system/flai-cli.md` | design | ADR-0130's consequences (T-1397) |

`flai/internal/preview/accept.go` reads `protected.Changed` and is left out: T-1395 keeps its signature and meaning. `touches suggest` listed generic documents; none else is taken.

Forecast: 35m. `flai forecast` gave 23m, 78 s per unit over 51 improvement stories, size 17. It is raised for moving `permission.go`'s 465 lines of thread handling into a new package with its tests unchanged.

Cost of delay: 38.72 USD a week, as `flai cod` works it out: 35m of the 7h32m forecast over E-0019's 12 stories, of the operator's 500 USD a week penalty. It stands.
