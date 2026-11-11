---
id: E-0034
type: epic
nature: feature
title: An operator defines the board's states, transitions, rules, and timers in a workflow file that flai interprets and enforces, starting from a baseline that is today's board
status: backlog
owner: alex
created: 2026-11-10T07:55:17Z
updated: 2026-11-10T07:55:17Z
transitions: []
tags: [cli, dashboard, template]
topics: [cli, template]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 20
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 7
          cache_read: 227338
          cache_write: 1251
          cost: 0.0646
---
# E-0034 An operator defines the board's states, transitions, rules, and timers in a workflow file that flai interprets and enforces, starting from a baseline that is today's board

## Outcome

An operator changes the board by editing a workflow file, not flai's code. The file, under `design/workflows/`, names the states, the transitions between them, who may make each, the rules each must pass, and the timers measured across them. flai interprets it and enforces it everywhere a move is made: `flai move`, `flai accept`, MCP's `item_move`, the dashboard, `flai serve`, and `flai check`. A project with no workflow file runs the baseline flai carries, as a project with no role file does ([ADR-0207](../../../design/adrs/0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md)), and the baseline is today's board exactly: the same six states, the same transitions, the same rules, the same figures.

What is fixed in code today, and becomes data:

| Today | Where | In the workflow file |
|-------|-------|----------------------|
| Six states, a closed list | `flai/internal/workitem/item.go` | The states, each with a name, a kind, the item types it takes, whether it has a WIP limit, and whether it has an order |
| The transition map | `flai/internal/workitem/rules.go` | The transitions, each with its from and to, who may make it, and whether it needs a reason |
| The rules in `Repo.Move` | `flai/internal/workitem/rules.go` | The rules each transition names, from a vocabulary flai implements, each refusing or warning |
| What a state means to flai | About 320 references in 20 packages | The state's kind and the transition's effect, which the code reads in place of a state's name |
| Lead, cycle, queue, and review time | `flai/internal/metrics`, `design/system/metrics.md` | The timers, each from one point in the workflow to another, or the sum of the time in the states it names |
| The columns, three lists | `flaiover/src` | Nothing: the dashboard reads the columns and the timers from flai |

A state's kind is what lets flai and the charts work on a workflow they have never seen: a queue, active work, a wait on someone else, accepted, or cancelled. A transition's effect names what flai does with it besides the move: the pull that starts a story's agent, the acceptance that merges and archives, a send back. A state's name is then only a name.

The timers are what makes time tracking flexible. Every state change is already recorded in an item's `transitions` ([ADR-0004](../../../design/adrs/0004-workflow-states-and-transitions.md)), so a timer is a reading of that history, and adding or changing one needs no new data and applies to past items too. The baseline's timers are the metrics as they stand. An operator adds a timer of their own, such as the time from the first move to review to the acceptance, and `flai stats` and the dashboard's charts report it beside the others.

Done when an operator adds a state between `in-progress` and `review` with a rule and a timer of its own, in the workflow file alone, and the board, the moves, the check, the agents' pull, the acceptance, and the charts all follow it with no change to flai or flaiover.

## Stories

## Notes

- The states and the state machine are ADR-0003's and ADR-0004's closed lists. The first story writes the ADR that supersedes that part of them and fixes the file's format, where it lives under `design/workflows/`, and the vocabulary of kinds, effects, and rules. No other story starts before it is accepted.
- The metrics definitions are the contract between `flai stats` and the dashboard and change only with an ADR. The timers' story writes it, and keeps every figure the baseline reports today byte for byte, tested against the existing fixtures.
- The rules are a vocabulary flai implements, as a role file's conditions are (`flai/internal/roles/conditions.go`), not an expression language. The baseline's are today's: acceptance criteria present, not a draft, at least one task, no open question in the narrative, nothing uncommitted in the worktree, every criterion ticked, every child closed. A rule the vocabulary lacks is a story that adds it.
- A role file's `when: to-review` and `when: to-backlog` become one condition that takes any state the workflow names.
- `flai check` validates the file and the items against it: a state no transition reaches, a state with no way to an accepted or cancelled one, a workflow with no pull or no acceptance, a timer that names a state the file lacks, and an item in a state the file lacks, with what to do about each.
- Changing a workflow with items on the board needs a rule of its own: renaming or removing a state says where its items go and how its history reads in the timers. The ADR decides it.
- The file is code the operator runs, as every file under `design/workflows/` is: it changes through a story the operator accepts, and no role may edit it.
- A proposed order, each a story or a few, for the planner to draft:
  1. The ADR and the design: the file's format, the kinds, the effects, the rule vocabulary, the timers, and migration.
  2. flai reads the baseline workflow it carries, and the states and the transition map come from it. Nothing else changes.
  3. The rules in `Repo.Move` become named rules the baseline's transitions list.
  4. Each package that names a state reads its kind or a transition's effect instead: `serve`, `check`, `cmd`, `preview`, `messages`, `mergequeue`, `planning`, `mcpserver`, `hostapi`, and the rest, one story for each group that can land alone.
  5. The timers: `flai stats` computes the baseline's figures from the file, then an operator's own.
  6. The dashboard reads the columns, their kinds, and the timers from flai, for the board and for every chart.
  7. A project's own workflow file: `flai check` validates it, `flai upgrade` keeps it, the template ships the baseline, and the operators' and users' guides say how to change it.
  8. The epic's closing story makes the change in Outcome's last paragraph in a project made from the template, end to end.
- Out of scope: more than one workflow in a project, a workflow per epic, and a workflow editor in the dashboard. Epics and tasks keep following the one workflow, with the item types each state and transition takes saying where they differ, as a task's move from `in-progress` to `done` does today.
