---
id: S-0553
type: story
nature: improvement
title: flai generates the planner's, the orchestrator's, the analyzer's, and the curator's harness definitions from their role files alone, flai check holds a .claude/agents copy to them, and a session started by hand reads the same definition
status: backlog
parent: E-0024
owner: alex
created: 2026-10-30T20:45:17Z
updated: 2026-10-31T07:23:38Z
transitions: []
tags: [cli, template]
topics: [cli, orchestration, planning, template, conventions]
touches: [design/adrs, design/adrs/README.md, design/system/agent-adapters.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md, flai/internal/roles/definition.go, flai/internal/roles/definition_test.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go, flai/internal/guard/writes.go, flai/internal/guard/writes_test.go, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/cmd/role.go, flai/cmd/role_test.go, ".claude/agents/planner.md", ".claude/agents/orchestrator.md", ".claude/agents/analyzer.md", ".claude/agents/curator.md", template/root/.claude/agents/planner.md, template/root/.claude/agents/orchestrator.md, template/root/.claude/agents/analyzer.md, template/root/.claude/agents/curator.md]
after: [S-0419]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 63
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 23
          output: 134
          cache_read: 863799
          cache_write: 99239
          cost: 0.2168
cost_of_delay:
  value: 58.47
  by: planner-E-0024
  at: 2026-10-30T20:46:31Z
forecast:
  duration: 1h5m
  delivery: 2026-10-31T20:50:00Z
  basis: "Its own forecast of 1h5m; 33rd in the pull order with an in-progress limit of 5, behind S-0424, S-0425, S-0548, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0506, S-0508, S-0509, S-0510, S-0516, S-0517, S-0518, S-0520, S-0420, S-0421, S-0422, S-0423, S-0540, S-0543, S-0545, S-0546, S-0547, S-0549, S-0550, S-0551 and S-0552."
  by: flai
  at: 2026-10-31T07:23:38Z
finalized:
  by: orchestrator
  at: 2026-10-30T20:48:03Z
---
# S-0553 flai generates the planner's, the orchestrator's, the analyzer's, and the curator's harness definitions from their role files alone, flai check holds a .claude/agents copy to them, and a session started by hand reads the same definition

## Goal

S-0419 made each built-in strategic role's file body the prompt of the definition flai passes in `--agents`. The definition's description, tools, and model still come from the hand-written `.claude/agents/<name>.md`, and `subAgents` refuses a start, naming `flai upgrade`, when that file is missing. So a change to a strategic agent's prompt is still made twice, in the role file and in `.claude/agents` (D-0411). ADR-0207 § 4 has flai generate the whole definition from the role file. This story does that for the planner, the orchestrator, the analyzer, and the curator. It settles, by a live trial, how a session the operator starts by hand reads the same definition. After it, a role file is the one place a strategic agent's definition is written. This is the definition-generation story that D-0351 and TH-0760 name.

## Acceptance criteria

- [ ] A session `flai serve` starts on `claude-code` for the planner, the orchestrator, the analyzer, or the curator gets its whole definition from its role file, the project's or the embedded baseline: `name`, `definition.description`, `definition.tools`, `definition.model` with the role's agent's model over it, and the body with its project additions as the prompt. It starts with no `.claude/agents/<name>.md` present, and the refusal naming `flai upgrade` for a missing strategic definition is gone.
- [ ] `TestClaudeCodeStartsThePlanner`, `TestClaudeCodeStartsTheOrchestrator`, `TestClaudeCodeStartsTheAnalyzer`, and `TestClaudeCodeStartsTheCurator` in `flai/internal/harness/harness_test.go` pass with their `.claude/agents` fixtures removed and the same start prompt text.
- [ ] `flai guard` reads a strategic role's tools from its role file, not from `.claude/agents/<name>.md`, and its guard tests pass unchanged.
- [ ] A live trial with Claude Code settles whether a session started by hand follows a symbolic link under `.claude/agents`. A new ADR refining ADR-0207 § 4 records the result, where flai writes the generated definition, and what the operator does once. `flai role definition <name>` prints the generated definition, and writes it there.
- [ ] `flai check` reports `roles.definition` when a built-in role's `.claude/agents/<name>.md` exists and is neither the link the ADR names nor the definition generated from its role file.
- [ ] D-0411's double edit ends. The strategic-agents convention, here and in the template, says a strategic agent's prompt is changed in its role file alone. `design/system/strategic-agents.md`, `design/system/agent-adapters.md`, `design/workflows/roles/README.md`, `docs/operators/index.md`, and `template/CHANGELOG.md` say how a definition is generated and read.

## Tasks
- T-2347 flai generates a built-in role's whole harness definition from its role file, and subAgents passes it with no .claude/agents read
- T-2348 A live Claude Code trial settles whether a session started by hand follows a symbolic link under .claude/agents, and an ADR refining ADR-0207 § 4 records where the generated definition lives
- T-2349 flai guard reads a strategic role's tools from its role file rather than from .claude/agents
- T-2350 flai check reports roles.definition for a built-in role's .claude/agents copy that is neither the ADR's link nor the definition its role file generates
- T-2351 flai role definition prints a role's generated definition and writes it where the ADR puts it
- T-2352 The planner's, orchestrator's, analyzer's, and curator's .claude/agents files, here and in the template, become what the ADR says, written by the operator
- T-2353 The strategic-agents convention and design, the roles README, the operator docs, and the changelog say a strategic agent's definition is generated from its role file and changed there alone

## Notes

- Drafted by the planner of E-0024 on 2026-10-30, on the operator's question on TH-0760. D-0351 put the generation in a story of its own after S-0419, and no story held it until now.
- Out of scope: the explorer's and the verifier's definitions, which are sub-agents' and not roles (ADR-0207 § 2), stay hand-written in `.claude/agents`. The start prompts and the built-in roles' triggers are the next story's.
- The `.claude/agents` pair is protected (ADR-0106). Its write is a task of its own, so the rest can close while the write waits on the operator, as S-0536 and S-0537 did (D-0411).

### Planning

- **After.** S-0419, whose `startOf` and role-file body in `subAgents` this story builds on, and which is in review. Nothing else is needed on main.
- **Touches.** 31, all files but one.
  - **From layout** (S-0419's worktree, read 2026-10-30): `flai/internal/harness/adapters.go`, where `subAgents` and `readDefinition` read `.claude/agents` and refuse naming `flai upgrade`; `harness_test.go`, whose start tests write the fixtures; `flai/internal/guard/writes.go` and its test, where `ReadOnlyAgents` parses a definition's `tools:` line; `flai/internal/check/check.go` and `roles_test.go` for `roles.definition`; `flai/cmd/role.go` and its test for `flai role definition`. `flai/internal/roles/definition.go` and its test are new, for the generation both the harness and the command use.
  - **From design:** ADR-0207 § 4, which names the generation and leaves the hand-started session open; `agent-adapters.md` and `strategic-agents.md`, which describe the definitions; the convention pair, the roles README pair, and `docs/operators/index.md`, which name `.claude/agents`; the four strategic `.claude/agents` files in both copies.
  - **Co-change, added:** `design/system/flai-cli.md` (38% of the seed commits), which describes `flai role`, and `docs/users/flai.md` and `flai-reference.md` for the new subcommand.
  - **Co-change, not added:** `docs/operators/settings.md`, `flaiover` docs, and `workflow.md`: no setting, page, or workflow changes. `flai/internal/harness/harness.go`: the start prompts are S-0554's; `flai task done` widens the touches if the request's fields move.
  - **Folder touch kept:** `design/adrs`, since the new ADR's file name is not known until it is written.
- **Forecast.** flai's 54m raised to 1h5m for the live trial, a short Claude Code session the touches do not count, as S-0423 and S-0424 were raised for their live runs. The delivery is flai's play-out, 33rd in the pull order, moved by the 11 minutes times the cycle factor 1.15.
- **Cost of delay.** 58.47 USD a week, worked out by hand as the planners of S-0414 to S-0425 did, because `flai cod` refuses while E-0024 has no inputs (TH-0567, question 3, confirmed; TH-0587 asks the operator to write 4h on the epic). 4h lost per cycle at `planning.hour_rate` 150 over the default 168h cycle gives E-0024 600 USD a week. S-0553's share is its 1h5m of the 11h7m forecast for E-0024's nine open stories, S-0419 to S-0425, S-0553, and S-0554: 65 / 667 × 600.
