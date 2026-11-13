---
id: S-0756
type: story
nature: feature
title: In a project made from the template, a state added between in-progress and review with a rule and a timer of its own, in the workflow file alone, is followed by the board, the moves, the check, the pull, the acceptance, and the charts
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:02:34Z
updated: 2026-11-13T14:36:34Z
transitions: []
tags: [cli, dashboard, template]
topics: [cli, dashboard, template]
touches: [flai/cmd/workflow_e2e_test.go, flai/internal/workflow/testdata/e2e-board.md, flaiover/src/lib/server/workflow.test.ts, design/system/workflow.md, docs/operators/index.md]
after: [S-0738, S-0739, S-0741, S-0742, S-0743, S-0744, S-0746, S-0747, S-0748, S-0749, S-0752, S-0753, S-0754, S-0755, S-0758]
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
      seconds: 10
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 14
          cache_read: 121381
          cache_write: 20507
          cost: 0.0405
draft: true
forecast:
  duration: 45m
  delivery: 2026-11-13T23:21:00Z
  basis: "Its own forecast of 45m; 25th in the pull order with an in-progress limit of 3, behind S-0662, S-0663, S-0771, S-0760, S-0759, S-0751, S-0757, S-0772, S-0241, S-0730, S-0736, S-0737, S-0738, S-0739, S-0741, S-0742, S-0743, S-0744, S-0746, S-0747, S-0748, S-0749, S-0750, S-0752, S-0753, S-0754 and S-0755."
  by: flai
  at: 2026-11-13T14:36:34Z
---
# S-0756 In a project made from the template, a state added between in-progress and review with a rule and a timer of its own, in the workflow file alone, is followed by the board, the moves, the check, the pull, the acceptance, and the charts

## Goal

The epic's outcome is shown end to end. In a project rendered from the template, the operator adds a state between `in-progress` and `review`, with a rule on the move into it and a timer of its own, by editing the workflow file alone, and every part of system-flow follows it with no change to flai or flaiover. What does not follow is a story of its own.

## Acceptance criteria

- [ ] A test renders the template into a new project and adds to its workflow file a state between `in-progress` and `review`, of the active-work or wait kind as the ADR allows, a rule on the move into it, and a timer from it to the acceptance
- [ ] `flai check --strict` passes on the project, and fails, naming the state, when the new state is left unreachable
- [ ] `flai board` and `board.get` show the new column in its place with its kind and WIP limit
- [ ] A story moves into the new state only when its rule passes, and on to `review` and `done` from it; `flai move`, `item_move`, and the dashboard's `item.move` each make and refuse the same moves
- [ ] flai serve keeps the story's agent running while the story is in the new state, and the acceptance merges and archives it as today
- [ ] `flai stats --json` and `stats.get` report the new timer for the story, and flaiover's Charts page draws it, with time in state stacking the new state
- [ ] The test runs in a tier `system-flow.yaml` declares, so `flai test` and the close-out select it for a story that changes `flai/internal/workflow`
- [ ] Anything that does not follow the new state is recorded as an issue or a story under E-0034 before this story goes to review

## Tasks

- T-3261 An integration test renders the template, adds a state with a rule and a timer to its workflow file, and drives the check, the board, the moves, serve, the acceptance, and stats through it
- T-3263 A flaiover test of the built flai draws the end-to-end project's board, item.move, and charts with the added state and its timer
- T-3264 Record each surface that did not follow the added state as an issue or a story under E-0034, and say in workflow.md that the outcome holds

## Notes

- This is the outcome's last paragraph in E-0034, made a test. It reads every other story's output, which is why its `after` names them.

### Planning

- Touches, all files:
  - `flai/cmd/workflow_e2e_test.go`, `flai/internal/workflow/testdata/e2e-board.md`: layout. A Go test that skips under `-short` runs in the integration tier, which `system-flow.yaml` already selects for a change under `flai/`, so criterion 7 needs no manifest change. The first draft's `scripts/workflow-e2e.sh` and `system-flow.yaml` were dropped for that reason.
  - `flaiover/src/lib/server/workflow.test.ts`: layout. A test importing `testing.ts` runs in the `flaiover-flai` tier against the built flai.
  - `design/system/workflow.md`, `docs/operators/index.md`: design, where the outcome is stated.
- `after`: every story whose surface the test drives, each for what this story reads from it:
  - S-0738: the rule on the move into the new state.
  - S-0739: the board view's columns and an epic that follows by kind.
  - S-0741: the `workflow.` findings, unreachable state included.
  - S-0742: the agent kept running in an active-work state.
  - S-0743: the deploy, forge, and merge queue left alone in the new state.
  - S-0744: `flai move` and `flai accept` by who may make each transition.
  - S-0746: `item_move` and `story_start`.
  - S-0747: the dashboard's `item.move`.
  - S-0748: a role started on a move into the new state.
  - S-0749: the guard's refusals and the forecast.
  - S-0752: the timer in `flai stats` and `stats.get`.
  - S-0753: the board's new column in flaiover.
  - S-0754: the timer's chart and time in state.
  - S-0755: the template's file, read in place of the baseline.
  - S-0758: the item page's moves and the cards.
- `flai check` may advise `wip.after-hub` on S-0737 or S-0755, which several stories name; it is advisory, and each `after` reads that story's output.
- Tasks: T-3261 (the Go test) first; T-3263 (the flaiover test) after it, using its fixture; T-3264 (recording what did not follow) after both.
- Forecast 45m, up from the 24m `flai forecast`'s rate gives for size 13. The test drives every surface of flai and flaiover, and each failure it finds is diagnosed here before it is recorded.
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
