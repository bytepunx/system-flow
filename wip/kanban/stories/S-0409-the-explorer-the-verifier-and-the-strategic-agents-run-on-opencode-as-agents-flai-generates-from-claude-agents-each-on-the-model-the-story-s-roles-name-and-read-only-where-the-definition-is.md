---
id: S-0409
type: story
nature: feature
title: The explorer, the verifier, and the strategic agents run on OpenCode as agents flai generates from .claude/agents, each on the model the story's roles name and read-only where the definition is
status: backlog
parent: E-0023
owner: alex
created: 2026-10-16T05:17:17Z
updated: 2026-10-23T17:56:28Z
transitions: []
tags: [cli]
topics: [agents, cli]
touches: [flai/internal/harness/opencodeagents.go, flai/internal/harness/opencodeagents_test.go, flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/opencodeconfig.go, flai/internal/harness/opencodeconfig_test.go, flai/internal/guard/opencode.go, flai/internal/guard/opencode_test.go, flai/tests/integration/opencode_test.go, design/system/agent-adapters.md, design/system/flai-cli.md, design/system/strategic-agents.md, docs/users/flai.md, docs/operators/index.md, flai/internal/harness/harness_test.go]
after: [S-0408]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 61.54
  by: planner-E-0023
  at: 2026-10-20T21:13:41Z
forecast:
  duration: 40m
  delivery: 2026-10-23T21:52:00Z
  basis: "Its own forecast of 40m; 5th in the pull order with an in-progress limit of 5, behind S-0429, S-0241, S-0406, S-0407 and S-0408."
  by: flai
  at: 2026-10-23T17:56:28Z
finalized:
  by: orchestrator
  at: 2026-10-20T21:14:06Z
---
# S-0409 The explorer, the verifier, and the strategic agents run on OpenCode as agents flai generates from .claude/agents, each on the model the story's roles name and read-only where the definition is

## Goal

A story on OpenCode hands search and verification to its explorer and verifier as one on Claude Code does, and the planner, the orchestrator, and the analyzer may run on OpenCode too. flai translates the five definitions in `.claude/agents/`, which stay the source (ADR-0130), into OpenCode agents at start: each sub-agent with the model the story's `agent.roles` names, as `provider/model`, and read-only permissions where its definition grants no edit; each strategic agent as a primary agent `flai serve` starts with `--agent`.

## Acceptance criteria

- [ ] flai translates each definition in `.claude/agents/` into an OpenCode agent, through a tool-name table (`Read`, `Grep`, `Glob`, `Bash`, `Edit`, `Write`, and `mcp__flai__*` to OpenCode's names): `explorer` and `verifier` with `mode: subagent`, the planner, the orchestrator, and the analyzer with `mode: primary`, each with its prompt and description. They land where S-0406's ADR puts the story's configuration, and the worktree stays clean.
- [ ] A sub-agent's model is the one the story's `agent.roles` names for its role, else the story's own; a read-only definition gets `edit: deny` and a shell that may not write outside the temporary folder (ADR-0141), and the guard's `opencode` reader sets `ReadOnly` for it.
- [ ] With `Guard: true`, the adapter's `Capabilities` say role models and strategic definitions, and `flai serve` starts a story on `opencode` with `roles`, and the planner, the orchestrator, and the analyzer on `opencode` with `--agent <name>`. With `Guard` false, both stay refused under ADR-0131, naming the setting that lifts them.
- [ ] The integration test shows the explorer started on its role's model and refused a write, on the known `opencode` version, and is skipped, saying so, without it.
- [ ] `docs/users/flai.md` § Sub-agents, `docs/operators/index.md` § Harnesses, `design/system/flai-cli.md`, `design/system/strategic-agents.md`, and `agent-adapters.md` say how the roles and the strategic agents run on OpenCode.

## Tasks
- T-1701 flai translates each definition in .claude/agents into an OpenCode agent through a tool-name table, with the role's model and read-only permissions where the definition grants no edit
- T-1702 The opencode adapter puts the translated agents in the story's configuration, starts a strategic agent with --agent, and states role models and strategic definitions where Guard is true
- T-1703 The guard's opencode reader marks a call from a read-only translated agent ReadOnly and a strategic agent's call with its role
- T-1704 The OpenCode integration test shows the explorer started on its role's model and refused a write
- T-1705 The user, operator, and design docs say how the roles and the strategic agents run on OpenCode

## Notes

- Drafted by the planner for E-0023, after S-0408. It needs S-0407's guard: ADR-0131 refuses roles on a harness without one.
- The epic names `.opencode/agents/` as where the agents go. If S-0406's ADR keeps the story's configuration out of the worktree, the agents go there too, as files or as `agent` entries, so that the move to review is not refused for a generated file.

### Planning

- **Touches.** None were declared. Every touch is the planner's:
  - **Layout:** `flai/internal/harness/opencodeagents.go`, new, reusing `readDefinition` in `adapters.go` unchanged; `opencode.go` and `opencodeconfig.go` for the capabilities and where the agents land, and `harness_test.go` for the capabilities test; `flai/internal/guard/opencode.go` for `ReadOnly`; `flai/tests/integration/opencode_test.go`. `flai/internal/serve/agents.go` reads the capabilities and needs no change.
  - **Design:** the docs, `strategic-agents.md`, and `agent-adapters.md` for criterion 5.
  - **Co-change:** `flai touches suggest` lists only hub documents, none of which this story's criteria reach; none was added.
- **Forecast.** 40m against flai's 31m. A tool-name table over five definitions, and strategic starts checked against the real binary, add about a quarter. The delivery, 2026-10-16T15:16Z, is flai's play-out moved by the raise times the cycle factor.
- **Cost of delay.** 61.54 USD a week. E-0023 carries `penalty_per_week: 500`, the operator's answer to TH-0467 question 1, confirmed on TH-0470 on 2026-10-20 (as on E-0019). `flai cod` shares it among the epic's seven open stories by forecast, 40m of 5h25m. Kept as flai gives it.
