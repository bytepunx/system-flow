---
id: ADR-0245
title: "A formula's role step that names the built-in analyzer starts it through the analyzer's own start, with the instance's focus and the instance as the run's item, so its end is the agent event on the instance"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0239]
topics: [orchestration, cli]
---

# ADR-0245 A formula's role step that names the built-in analyzer starts it through the analyzer's own start, with the instance's focus and the instance as the run's item, so its end is the agent event on the instance

## Context

[ADR-0239](0239-a-formula-s-instance-is-worked-by-its-story-s-agent-step-by-step-a-role-step.md) §2 lets the agent of a formula's instance start, with `role_run` or `flai role run`, the role a role step of its formula names, on the instance, and wait for the run's end as an `agent` event on the instance. The run goes through `serve.RunRole`, a defined role's start.

The `analysis-cadence` formula (S-0424) has a role step that names the built-in `analyzer`. `flai check` accepted it, since the analyzer has a role file, but the step could not be worked:

- `serve.RunRole` refuses every built-in role, naming the command that starts it (`definedRole`).
- `flai analyze` and the MCP tool `analyze` refuse an agent `flai serve` started, and a story's agent is one.
- The analyzer's run is recorded as the project's analyzer, on no item, and its `agent` events are on the analyzer's role (`id` `analyzer`, `type` `role`). The instance's agent waits for an event on its instance, so it would never see the end.

The planner, the orchestrator, and the curator have the same refusal, and a formula has no reason to start them: each keeps a start of its own with its own gates.

## Decision

**A role step that names the built-in analyzer starts it through the analyzer's own start, with the instance's focus and the instance as the run's, so its end is the `agent` event on the instance. No other built-in role is started for a role step, and `flai check` reports a role step that names one.**

1. `runTrigger` in `flai/cmd/role.go` gives the agent of a formula's instance the trigger `formula <name> <S-nnnn>` for the analyzer, as for any role a role step of its formula names, on the instance alone.
2. `serve.RunRole`, given that trigger and the name `analyzer`, does not start a defined role. It starts the analyzer through the analyzer's own start (`analyzeStep`, then `Analyze`'s path): `analyzeCheck`'s checks, the `analyze` host action, one analyzer per project, the budget state, and the analyzer's agent and prompt.
3. The focus is the run's input `focus` when the agent gives one. Otherwise it is the `focus` input that `flai formula run` recorded in the instance's Notes (``with the inputs `focus=<value>` ``). Empty means all three. Other inputs, the step's prompt among them, are not given to the analyzer: its own prompt and focus are what it runs on.
4. An item other than the instance is refused, naming the instance.
5. The run is the project's analyzer run, recorded under `analyzer` in `serve/agents.json` with a new field, `instance`, the instance's ID. `item` stays empty, so the run is still told apart as the analyzer's. Its `agent` events, `started` and then `ended` or `failed`, are on the instance (`type` `story`), by `analyzer`, with `fields.role` `analyze`. The analyzer's other runs keep their events on the analyzer's role.
6. The analyzer's activity entry, its report, and the issues it charges are as for any analyzer run. The run charges nothing to the instance until S-0425.
7. The MCP tool `role_run` journals such a start under the `analyze` action, method `mcp.role_run`, and answers with the role `analyzer`, the instance as item, and the action `analyze`.
8. `serve.RunRole` refuses the planner, the orchestrator, and the curator for a role step: "<role> is a built-in role a role step cannot start: only the analyzer is started for a role step, through flai analyze's start". It still refuses every built-in role to the operator's own `flai role run`, naming its command (`flai analyze` for the analyzer).
9. `flai check` (`formulas.file`) reports a role step that names the planner, the orchestrator, or the curator, even where the project has the role's file.

## Consequences

- `analysis-cadence` can be worked step by step: its agent starts the analyzer, holds `wait_for_events`, and takes the analyzer's end on its instance as the step's exit, `failed` when the run wrote no report.
- The analyzer still runs once per project at a time. A role step asked while another analyzer runs is refused, and the step fails as ADR-0239 §3 says.
- The `analyze` host action gates the analyzer from a formula as it does from the operator and the schedule. The exception to the refusal of an agent `flai serve` started stays narrow: one agent, its instance, and the analyzer only where its formula names it.
- `serve/agents.json` gains an optional `instance` field on the analyzer's run. Readers that do not know it, the dashboard among them, pass over it.
- A formula that names the planner, the orchestrator, or the curator in a role step no longer passes `flai check`.

## Alternatives considered

- **Setting the analyzer run's `item` to the instance.** `item` is how `serve` tells a planner run from the analyzer's (`analyzes`, `put`, `of`, `agentAbout`). A run with an item would be filed as a planner run of the instance. A field of its own keeps the record's shape.
- **A defined role file for the analyzer's step.** A project role that runs an analysis would duplicate the analyzer's prompt, checks, report detection, and issue charging, and would bypass `analyze`'s single-run lock.
- **Lifting `flai analyze`'s refusal for an instance's agent.** The step would call another command than `role_run`, against ADR-0239 §2, and the run would still end on the analyzer's role, not the instance.
- **Taking the focus only from the run's inputs.** The step's task gives only the filled prompt, so the agent would have to parse the focus from the prompt or the Notes itself. flai reads the Notes it wrote, and an input the agent gives still wins.
- **Starting every built-in role for a role step.** The planner, the orchestrator, and the curator each have gates and triggers of their own, and no formula needs them. Refusing them, and reporting them in `flai check`, keeps the exception to the analyzer.
