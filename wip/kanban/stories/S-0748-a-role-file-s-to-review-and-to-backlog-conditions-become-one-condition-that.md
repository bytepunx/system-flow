---
id: S-0748
type: story
nature: feature
title: A role file's to-review and to-backlog conditions become one condition that names any state the workflow has, and the baseline role files use it
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:00:44Z
updated: 2026-11-12T16:00:44Z
transitions: []
tags: [cli, template]
topics: [cli, template, orchestration, conventions]
touches: [flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/roles/roles.go, flai/internal/roles/start.go, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/planner.md, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/replan.go, flai/internal/check/check_roles.go, flai/internal/check/roles_test.go, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/planner.md, template/root/design/workflows/roles/README.md, design/workflows/roles/orchestrator.md, design/workflows/roles/planner.md, design/workflows/roles/README.md, design/system/strategic-agents.md, template/CHANGELOG.md]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0748 A role file's to-review and to-backlog conditions become one condition that names any state the workflow has, and the baseline role files use it

## Goal

A role file names the state a move goes to with one condition that takes any state the workflow has, in place of the two fixed `to-review` and `to-backlog` conditions, so a project's role reacts to a move into a state of its own workflow. The baseline role files say the same as today with the new condition, and a role file that still names the old ones keeps working or is told what to write.

## Acceptance criteria

- [ ] `flai/internal/roles/conditions.go` has one condition, as the ADR names it, that qualifies a `moved` event and the `item_move` and `flai move` calls with the state the move goes to, and refuses at load a state the workflow lacks
- [ ] The baseline role files in `template/` and in `flai/internal/roles/baseline/`, byte for byte, and this project's own under `design/workflows/roles/`, use the new condition, and the orchestrator's and the planner's triggers and grants are what they are today
- [ ] A role file naming `to-review` or `to-backlog` either works as before or gets a `roles.file` finding naming the condition to write instead, as the ADR decides
- [ ] `flai guard` and flai serve's triggers read the new condition, and a test shows a role started on a move into a state added between `in-progress` and `review`
- [ ] `design/system/strategic-agents.md` § Role files and `design/workflows/roles/README.md`, here and in the template, list the condition, and `template/CHANGELOG.md` records the change

## Tasks

## Notes

- Today `to-review` and `to-backlog` are the `matches` map in `conditions.go` at lines 208 and 209, and the condition table at 251 to 273. `baseline_test.go` holds the embedded copies to the template's.
- The role files are code the operator runs, changed here through a story the operator accepts, as ADR-0207 §7 asks.
