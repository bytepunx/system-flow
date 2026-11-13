---
id: S-0738
type: story
nature: feature
title: The rules Repo.Move enforces are named rules from a vocabulary flai implements, listed on the baseline's transitions, each refusing or warning
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T15:58:55Z
updated: 2026-11-13T07:23:59Z
transitions: []
tags: [cli]
topics: [cli]
touches: [flai/internal/workflow/rules.go, flai/internal/workflow/rules_test.go, flai/internal/workflow/workflow.go, flai/internal/workflow/baseline/board.md, flai/internal/workflow/baseline_test.go, flai/internal/workitem/rules.go, flai/internal/workitem/rules_test.go, flai/internal/workitem/workitem_test.go, flai/internal/workitem/worktree_test.go, design/system/workflow.md]
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
      seconds: 12
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 15
          cache_read: 108315
          cache_write: 18020
          cost: 0.0361
draft: true
forecast:
  duration: 29m
  delivery: 2026-11-13T15:10:00Z
  basis: "Its own forecast of 29m; 12th in the pull order with an in-progress limit of 3, behind S-0659, S-0241, S-0660, S-0661, S-0662, S-0663, S-0730, S-0733, S-0734, S-0735, S-0736 and S-0737."
  by: flai
  at: 2026-11-13T07:23:59Z
---
# S-0738 The rules Repo.Move enforces are named rules from a vocabulary flai implements, listed on the baseline's transitions, each refusing or warning

## Goal

The checks `Repo.Move` makes in code today become named rules that the workflow's transitions list. flai implements each rule, and the baseline names them on the transitions where they apply today, so a move refuses and warns exactly as it does now, and a workflow of the operator's own puts a rule on any transition it names.

## Acceptance criteria

- [ ] `flai/internal/workflow` names the rule vocabulary the ADR fixes: acceptance criteria present, not a draft, at least one task, no open question in the narrative, nothing uncommitted in the worktree, every criterion ticked, every child closed, and a reason, each refusing or warning as the ADR says
- [ ] The baseline's transitions list those rules where `Repo.Move` applies them today: to `ready`, to `review`, to `done`, to `cancelled`, and `review` back to `in-progress`
- [ ] `Repo.Move` runs the rules the transition lists, in their order, and each refusal and warning keeps today's text word for word
- [ ] A workflow whose transition names a rule the vocabulary lacks is refused when it is parsed, naming the rule and the transition
- [ ] Tests cover each rule refusing and, where it warns, warning, the uncommitted-worktree rule and the every-child-closed rule included, which have no `workitem` test today
- [ ] `design/system/workflow.md` § Transitions and who makes them lists the baseline's rules by name

## Tasks

- T-3189 The workflow package names the rule vocabulary, each rule refusing or warning, and refuses a transition naming a rule it lacks
- T-3197 The baseline workflow lists today's rules on the transitions where Repo.Move applies them
- T-3199 workflow.md lists the baseline's rules by name on each transition
- T-3202 Repo.Move runs the rules its transition lists, in their order, with today's refusals and warnings word for word, and each rule is tested

## Notes

- Today's rules are in `flai/internal/workitem/rules.go`, `Repo.Move`, lines 58 to 118: the reason at 58, criteria and draft at 72 to 80, the task and the open question at 82 to 91, the worktree at 95 to 102, the unticked warning at 105, and the children and criteria at done at 110 to 118.

### Planning

- Touches, all files:
  - `flai/internal/workflow/rules.go`, `rules_test.go`, `workflow.go`, `baseline/board.md`, `baseline_test.go`: layout. The vocabulary and the baseline's rules live where S-0737 puts the parser and the baseline.
  - `flai/internal/workitem/rules.go`, `rules_test.go`, `workitem_test.go`, `worktree_test.go`: layout. `Repo.Move` and the tests of its rules; `TestMoveRules` is in `workitem_test.go`, the worktree rule's helpers in `worktree.go`.
  - `design/system/workflow.md`: design (criterion 6).
  - Left out from `flai touches suggest`: `work-management.md` here and in the template (17% to 21%). Its definitions of ready and done keep their words, since the rules do not change.
- `after`: S-0737. This story adds rules to the parser and baseline S-0737 writes, and `Repo.Move` reads them through the workflow `Repo` carries from S-0737 on.
- Tasks: T-3189 (the vocabulary); then T-3197 (the baseline lists the rules) and T-3199 (`workflow.md`) together; then T-3202 (`Repo.Move` runs them). T-3199 shares no path with T-3197.
- Forecast 29m, as `flai forecast` gives it for size 16. The rules keep their code and their words; only how they are chosen changes.
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
