---
id: S-0748
type: story
nature: feature
title: A role file's to-review and to-backlog conditions become one condition that names any state the workflow has, and the baseline role files use it
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:00:44Z
updated: 2026-11-13T22:47:30Z
transitions: []
tags: [cli, template]
topics: [cli, template, orchestration, conventions]
touches: [flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/roles/roles.go, flai/internal/roles/start.go, flai/internal/check/check_roles.go, flai/internal/check/roles_test.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/replan.go, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/planner.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/planner.md, design/workflows/roles/orchestrator.md, design/workflows/roles/planner.md, template/CHANGELOG.md, template/template.yaml, design/system/strategic-agents.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, docs/operators/settings.md, docs/operators/index.md, docs/users/flai.md]
after: [S-0737]
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
      seconds: 11
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 13
          cache_read: 108313
          cache_write: 18018
          cost: 0.0361
draft: true
forecast:
  duration: 53m
  delivery: 2026-11-14T03:17:00Z
  basis: "Its own forecast of 53m; 13th in the pull order with an in-progress limit of 3, behind S-0783, S-0241, S-0736, S-0737, S-0738, S-0739, S-0741, S-0742, S-0743, S-0744, S-0746 and S-0747."
  by: flai
  at: 2026-11-13T22:47:30Z
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

- T-3231 The roles package has one condition naming the state a move goes to, and treats to-review and to-backlog as the ADR decides
- T-3238 flai guard and flai serve's triggers read the new move condition, and a role starts on a move into a state of the workflow's own
- T-3239 The baseline orchestrator and planner role files use the new move condition in the template, in flai's embedded copies, and in this project
- T-3240 strategic-agents.md, the roles READMEs, and the operators' and users' guides list the move condition in place of to-review and to-backlog

## Notes

- Today `to-review` and `to-backlog` are the `matches` map in `conditions.go` at lines 208 and 209, and the condition table at 251 to 273. `baseline_test.go` holds the embedded copies to the template's.
- The role files are code the operator runs, changed here through a story the operator accepts, as ADR-0207 §7 asks.

### Planning

- Touches, all files:
  - Layout: `roles/conditions.go` and the files that load and start roles; `check/check_roles.go` for the old names' finding; `guard/roles.go` and `serve`'s `orchestrate_reasons.go` and `replan.go`, the three readers of the conditions.
  - Design: the baseline role files in the template, their embedded copies, and this project's own, with their READMEs, `strategic-agents.md`, and `template/CHANGELOG.md` (criteria 2 and 5).
  - Co-change: `template/template.yaml` (18% in `flai touches suggest`), bumped with any template change.
  - Search: `docs/operators/settings.md`, `docs/operators/index.md`, and `docs/users/flai.md` name `to-review` or `to-backlog` today.
  - `orchestrate_reasons.go` and `replan.go` are also S-0742's; the trial merge judges that.
- `after`: S-0737. The condition refuses a state the workflow lacks, read from the workflow S-0737 gives. The condition's name is S-0736's ADR's, which S-0737 already waits for.
- Tasks: T-3231 (the condition) first; then T-3238 (the guard and serve), T-3239 (the role files), and T-3240 (the documents) together, each reading T-3231's condition and sharing no path with another.
- Forecast 53m, `flai forecast`'s rate for size 29, up from its 46m for the first draft's 25 once the template version and the guides were added.
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
