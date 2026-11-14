---
id: S-0736
type: story
nature: feature
title: An ADR fixes the workflow file's format and place under design/workflows, its state kinds, transition effects, rules, and timers, and how a workflow changes with items on the board
status: review
parent: E-0034
owner: alex
created: 2026-11-12T15:58:31Z
updated: 2026-11-14T18:03:13Z
transitions:
  - to: ready
    at: 2026-11-14T17:38:24Z
    by: alex
  - to: in-progress
    at: 2026-11-14T17:38:32Z
    by: agent-S-0736
  - to: review
    at: 2026-11-14T17:56:47Z
    by: agent-S-0736
  - to: in-progress
    at: 2026-11-14T17:58:41Z
    by: alex
  - to: review
    at: 2026-11-14T18:03:13Z
    by: alex
blocked:
  - from: 2026-11-14T17:48:02Z
    until: 2026-11-14T17:54:50Z
    reason: awaiting the operator on TH-1198
tags: [cli]
topics: [cli, template, conventions]
touches: [design/adrs/README.md, design/system/workflow.md, design/system/work-hierarchy.md, design/workflows/README.md, design/adrs/0301-the-board-is-a-workflow-file-design-workflows-board-md-whose-front-matter-names.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 692
  priced_by: harness
  turns:
    - day: 2026-11-14
      ceremony: 3
      work: 36
  models:
    - model: claude-opus-5-5
      input: 138
      output: 54553
      cache_read: 14668132
      cache_write: 375812
      cost: 6.8349
  strategic:
    - kind: planner
      seconds: 31
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 577
          cache_read: 190240
          cache_write: 4297
          cost: 0.084
    - kind: orchestrator
      seconds: 386
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 93
          output: 605
          cache_read: 4039786
          cache_write: 513027
          cost: 1.3482
        - model: claude-sonnet-5-5
          input: 7
          output: 54
          cache_read: 46400
          cache_write: 45995
          cost: 0.0721
cost_of_delay:
  value: 80.17
  by: system-flow
  at: 2026-11-14T17:44:01Z
forecast:
  duration: 45m
  delivery: 2026-11-14T18:27:00Z
  basis: "Its own forecast of 45m; 2nd in the pull order with an in-progress limit of 3, behind S-0781 and S-0241."
  by: flai
  at: 2026-11-14T17:36:12Z
finalized:
  by: alex
  at: 2026-11-14T17:38:22Z
---
# S-0736 An ADR fixes the workflow file's format and place under design/workflows, its state kinds, transition effects, rules, and timers, and how a workflow changes with items on the board

## Goal

Before any code reads a workflow file, one accepted ADR settles what the file is, so that every later story of E-0034 builds to one format. It supersedes the closed list of states and the state machine in ADR-0003 and ADR-0004, and the design says how the board is defined.

## Acceptance criteria

- [x] An accepted ADR fixes the workflow file's path under `design/workflows/` and its format: each field with its type, meaning, and allowed values, as ADR-0207 §2 does for a role file
- [x] It names the state kinds (a queue, active work, a wait on someone else, accepted, cancelled), what flai does with each kind, and the kind of each of today's six states, with the item types, WIP limit, and order each state takes
- [x] It names the transition effects (the pull that starts a story's agent, the acceptance that merges and archives, a send back, and any other today's code acts on) and who may make a transition, and maps each of today's transitions to them, the moves back one column of ADR-0055 and a task's move from `in-progress` to `done` included
- [x] It names the rule vocabulary flai implements, each rule refusing or warning, and lists the rules each baseline transition names: acceptance criteria present, not a draft, at least one task, no open question in the narrative, nothing uncommitted in the worktree, every criterion ticked, every child closed, and a reason where one is needed
- [x] It defines a timer, from one point in the workflow to another or as the sum of the time in the states it names, and gives lead, cycle, queue, and review time as the baseline's timers, leaving the change to `metrics.md` to the timers' story and its own ADR
- [x] It says how a workflow changes with items on the board: renaming or removing a state names where its items go and how their history reads in the timers
- [x] It says that a project with no workflow file runs the baseline flai carries, which is today's board exactly, and that the file changes only through a story the operator accepts, outside every role's `files`
- [x] ADR-0003 and ADR-0004 carry the new ADR in `superseded_by`, or it refines them, as the ADR decides
- [x] `design/system/workflow.md`, `design/system/work-hierarchy.md` § States, and `design/workflows/README.md` describe the workflow file and link the ADR

## Tasks

- T-3187 Write the ADR, proposed, that fixes the workflow file's format, kinds, effects, rules, timers, and migration
- T-3193 Settle the ADR's open choices with the operator, accept it, and mark ADR-0003 and ADR-0004 superseded or refined by it
- T-3200 workflow.md, work-hierarchy.md § States, and the workflows README describe the workflow file and link the accepted ADR

## Notes

- The epic's Notes list what the ADR decides. The rules are a vocabulary flai implements, as a role file's conditions are (`flai/internal/roles/conditions.go`), not an expression language.
- No other story of E-0034 starts before this one is accepted.

### Planning

- Touches, all files, checked again on 2026-11-13 against `flai touches suggest`:
  - `design/adrs/README.md`: co-change. `flai adr new` and `flai adr accept` write the index again.
  - `design/adrs/0003-work-item-hierarchy.md`, `design/adrs/0004-workflow-states-and-transitions.md`: design (criterion 8). They change only by `superseded_by`.
  - `design/system/workflow.md`, `design/system/work-hierarchy.md`, `design/workflows/README.md`: design (criterion 9).
  - No folder touch on `design/adrs`. The ADR's file takes the number `flai adr new` gives it when the story runs.
  - Left out from `flai touches suggest`: `flai-cli.md`, the user and operator guides, `flaiover-dashboard.md`, `metrics.md`, and the template. They describe what is built, and the stories that build it change them. Criterion 5 leaves `metrics.md` to S-0750.
- No `after`: the story reads no other story's output. Every other E-0034 story builds on this ADR, most through S-0737.
- Tasks: T-3187 (the ADR, proposed), then T-3193 (the operator's answers, the acceptance, ADR-0003 and ADR-0004), then T-3200 (the design). One task per layer: each reads the one before. All three still fit the outcome, so none was changed.
- Forecast 45m, kept over the 28m `flai forecast` gives for size 15 (110 s per unit, 131 done feature stories). The ADR settles as much as ADR-0207 did (32 KB, eight decisions), and T-3193 waits on the operator in the middle. The wait is wall time, not agent time, but rewriting the ADR after the answers is agent time. The delivery is the one flai's replan writes from the pull order.
- Cost of delay: 80.17 USD a week, as `flai cod` gives it: this story's share of E-0034's 1500.00 USD a week, the operator's `time_lost_per_cycle: 10h` (TH-1115) at 150 USD an hour, by its forecast of 45m over the 14h2m of the epic's 19 open stories. Kept: the epic's outcome needs every story, so a share by size is the split that fits.
- 2026-11-14T17:58:41Z: moved to in-progress: Rebase needed
