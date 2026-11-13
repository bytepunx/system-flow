---
id: S-0736
type: story
nature: feature
title: An ADR fixes the workflow file's format and place under design/workflows, its state kinds, transition effects, rules, and timers, and how a workflow changes with items on the board
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T15:58:31Z
updated: 2026-11-13T13:22:24Z
transitions: []
tags: [cli]
topics: [cli, template, conventions]
touches: [design/adrs/README.md, design/adrs/0003-work-item-hierarchy.md, design/adrs/0004-workflow-states-and-transitions.md, design/system/workflow.md, design/system/work-hierarchy.md, design/workflows/README.md]
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
          input: 4
          output: 16
          cache_read: 108315
          cache_write: 18020
          cost: 0.0361
draft: true
forecast:
  duration: 45m
  delivery: 2026-11-13T17:01:00Z
  basis: "Its own forecast of 45m; 9th in the pull order with an in-progress limit of 3, behind S-0661, S-0745, S-0760, S-0759, S-0751, S-0757, S-0663, S-0662, S-0241 and S-0730."
  by: flai
  at: 2026-11-13T13:22:24Z
---
# S-0736 An ADR fixes the workflow file's format and place under design/workflows, its state kinds, transition effects, rules, and timers, and how a workflow changes with items on the board

## Goal

Before any code reads a workflow file, one accepted ADR settles what the file is, so that every later story of E-0034 builds to one format. It supersedes the closed list of states and the state machine in ADR-0003 and ADR-0004, and the design says how the board is defined.

## Acceptance criteria

- [ ] An accepted ADR fixes the workflow file's path under `design/workflows/` and its format: each field with its type, meaning, and allowed values, as ADR-0207 §2 does for a role file
- [ ] It names the state kinds (a queue, active work, a wait on someone else, accepted, cancelled), what flai does with each kind, and the kind of each of today's six states, with the item types, WIP limit, and order each state takes
- [ ] It names the transition effects (the pull that starts a story's agent, the acceptance that merges and archives, a send back, and any other today's code acts on) and who may make a transition, and maps each of today's transitions to them, the moves back one column of ADR-0055 and a task's move from `in-progress` to `done` included
- [ ] It names the rule vocabulary flai implements, each rule refusing or warning, and lists the rules each baseline transition names: acceptance criteria present, not a draft, at least one task, no open question in the narrative, nothing uncommitted in the worktree, every criterion ticked, every child closed, and a reason where one is needed
- [ ] It defines a timer, from one point in the workflow to another or as the sum of the time in the states it names, and gives lead, cycle, queue, and review time as the baseline's timers, leaving the change to `metrics.md` to the timers' story and its own ADR
- [ ] It says how a workflow changes with items on the board: renaming or removing a state names where its items go and how their history reads in the timers
- [ ] It says that a project with no workflow file runs the baseline flai carries, which is today's board exactly, and that the file changes only through a story the operator accepts, outside every role's `files`
- [ ] ADR-0003 and ADR-0004 carry the new ADR in `superseded_by`, or it refines them, as the ADR decides
- [ ] `design/system/workflow.md`, `design/system/work-hierarchy.md` § States, and `design/workflows/README.md` describe the workflow file and link the ADR

## Tasks

- T-3187 Write the ADR, proposed, that fixes the workflow file's format, kinds, effects, rules, timers, and migration
- T-3193 Settle the ADR's open choices with the operator, accept it, and mark ADR-0003 and ADR-0004 superseded or refined by it
- T-3200 workflow.md, work-hierarchy.md § States, and the workflows README describe the workflow file and link the accepted ADR

## Notes

- The epic's Notes list what the ADR decides. The rules are a vocabulary flai implements, as a role file's conditions are (`flai/internal/roles/conditions.go`), not an expression language.
- No other story of E-0034 starts before this one is accepted.

### Planning

- Touches, all files:
  - `design/adrs/README.md`: co-change. `flai adr new` and `flai adr accept` write the index again.
  - `design/adrs/0003-work-item-hierarchy.md`, `design/adrs/0004-workflow-states-and-transitions.md`: design (criterion 8). They change only by `superseded_by`.
  - `design/system/workflow.md`, `design/system/work-hierarchy.md`, `design/workflows/README.md`: design (criterion 9).
  - No folder touch on `design/adrs`. The ADR's file takes the number `flai adr new` gives it when the story runs.
  - Left out from `flai touches suggest`: `flai-cli.md`, the user and operator guides, `flaiover-dashboard.md`, and the template. They describe what is built, and the stories that build it change them.
- No `after`: the story reads no other story's output. Every other E-0034 story builds on this ADR, most through S-0737.
- Tasks: T-3187 (the ADR, proposed), then T-3193 (the operator's answers, the acceptance, ADR-0003 and ADR-0004), then T-3200 (the design). One task per layer: each reads the one before.
- Forecast 45m, up from the 28m `flai forecast` gives for size 15. The ADR settles as much as ADR-0207 did (32 KB, eight decisions), and T-3193 waits on the operator in the middle. The wait is wall time, not agent time, but the ADR's rewrite after the answers is agent time.
- Cost of delay: not set yet. E-0034 has no inputs and neither does this story, so `flai cod` refuses it. I asked on TH-1097 and recommended `time_lost_per_cycle: 2h` on the epic. Once the operator answers, the value is this story's share by forecast.
