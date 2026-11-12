---
id: S-0738
type: story
nature: feature
title: The rules Repo.Move enforces are named rules from a vocabulary flai implements, listed on the baseline's transitions, each refusing or warning
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T15:58:55Z
updated: 2026-11-12T15:58:55Z
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
draft: true
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

## Notes

- Today's rules are in `flai/internal/workitem/rules.go`, `Repo.Move`, lines 58 to 118: the reason at 58, criteria and draft at 72 to 80, the task and the open question at 82 to 91, the worktree at 95 to 102, the unticked warning at 105, and the children and criteria at done at 110 to 118.
