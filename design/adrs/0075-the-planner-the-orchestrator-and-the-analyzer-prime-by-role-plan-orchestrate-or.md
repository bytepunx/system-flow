---
id: ADR-0075
title: "The planner, the orchestrator, and the analyzer prime by role, plan, orchestrate, or analyze, with a full pack of their own and no story"
status: accepted
date: 2026-10-03
supersedes: []
superseded_by: []
refines: [ADR-0059, ADR-0068]
topics: [conventions, cli, planning, orchestration, analysis]
---

# ADR-0075 The planner, the orchestrator, and the analyzer prime by role, plan, orchestrate, or analyze, with a full pack of their own and no story

## Context

E-0016 adds three agents that work above a story. The planner drafts an epic's stories and enriches a story. The orchestrator keeps work moving within the permissions the operator sets. The analyzer reads the metrics, the design, and the code, and files issues. `flai prime` primed only an agent with a story: `--story` for the story's agent, and `--role explore|verify` with `--story` for the sub-agents it hands work to (ADR-0059, ADR-0068). None of the three has a story of its own. The orchestrator and the analyzer work on the whole project, and the planner on one epic or one story. S-0196 had added `orchestrator`, `planner`, and `analyzer` to the roles `flai check` accepts in a convention's front matter, but no convention listed them and `flai prime` took none, so nothing gave these agents the rules they work by.

## Decision

The strategic agents prime by role: `plan` for the planner, `orchestrate` for the orchestrator, and `analyze` for the analyzer. The roles are verbs, like `explore` and `verify`, and replace the nouns S-0196 had added. `flai check` accepts `story`, `explore`, `verify`, `plan`, `orchestrate`, and `analyze` in a convention's `roles`, and warns (`conventions.roles`) about any other value.

`flai prime --role plan --epic E-nnnn` or `--story S-nnnn`, `flai prime --role orchestrate`, `flai prime --role analyze`, and the MCP tool `prime` with `role` and, for the planner, `epic` or `story`, print a strategic agent's pack:

- Its topics are the role's, `planning`, `orchestration`, or `analysis`, from the source kind `role`, and for the planner the item's as well.
- The conventions whose `roles` are empty or list the role, with the sections its topics leave out taken out. No README.
- The open issues.
- For the planner, what the item names, whole, as a story's pack loads it: a story's with its epic and tasks, an epic's alone. Then the sections ranked highest against the item, while the budget has room.
- Briefs of the design, tech, and ADRs its topics select, and of the ADRs one link step reaches.
- A catalog of the rest.

The budget is the story's agent's, the project's `prime.budget` or 80 KB, not half of it; `--budget` sets it. The pack's JSON names the planner's epic or story as `item` and has no `story`. The orchestrator and the analyzer refuse an item. The planner refuses no item, a task, both `--epic` and `--story`, and an ID of the other type than its flag.

The baseline gains `strategic-agents.md`, order 140, read by `plan`, `orchestrate`, and `analyze`, and the baseline's `roles` name the strategic agents on the files they read, as the designer decided on TH-0081 (2026-10-03). This refines ADR-0068's list of roles and ADR-0059's priming by role; the rest of both stands.

## Consequences

- Besides `strategic-agents.md` and the conventions without roles, the planner reads `work-management.md`; the orchestrator `work-management.md` and `decisions.md`; the analyzer `continuous-improvement.md`, `code-quality.md`, and `telemetry.md`.
- `telemetry.md` gains roles and so leaves the planner's and the orchestrator's packs; the story's agent and the sub-agents read it as before. The story's agent does not read `strategic-agents.md`.
- A project's own convention that lists `orchestrator`, `planner`, or `analyzer` draws a `conventions.roles` warning until it uses the verb. The baseline never listed them.
- `design/system/metrics.md` carries the topic `analysis`, so the analyzer's pack briefs it. `work-hierarchy.md` and `workflow.md` carry `all`, which every strategic pack selects already. The design E-0016's later stories write, such as `strategic-agents.md`, carries the role topics.
- `design/system/conventions.md`, `agent-context.md`, `flai-cli.md`, and the user guide describe the roles and the packs. Template 1.0.38 carries the convention and the roles.

## Alternatives considered

- Keep the nouns S-0196 had added: the other roles say what the agent does, and `--role` reads as an action.
- A pack shaped like a sub-agent's, briefs only in half the budget: these agents are top level, with no story's agent holding the detail for them, and the planner must read what its item names to plan it.
- An epic's pack that loads what its stories name as well: on E-0016 it was 131 KB, past what a tool result carries. The planner reads the epic's stories as items instead.
