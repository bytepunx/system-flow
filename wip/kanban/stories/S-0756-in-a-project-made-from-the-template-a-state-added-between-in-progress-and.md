---
id: S-0756
type: story
nature: feature
title: In a project made from the template, a state added between in-progress and review with a rule and a timer of its own, in the workflow file alone, is followed by the board, the moves, the check, the pull, the acceptance, and the charts
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:02:34Z
updated: 2026-11-12T16:02:34Z
transitions: []
tags: [cli, dashboard, template]
topics: [cli, dashboard, template]
touches: [flai/cmd/workflow_e2e_test.go, flai/internal/workflow/testdata/e2e-board.md, flaiover/src/lib/server/workflow.test.ts, scripts/workflow-e2e.sh, system-flow.yaml, design/system/workflow.md, docs/operators/index.md]
after: [S-0738, S-0739, S-0741, S-0742, S-0743, S-0744, S-0746, S-0747, S-0748, S-0749, S-0752, S-0753, S-0754, S-0755]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
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

## Notes

- This is the outcome's last paragraph in E-0034, made a test. It reads every other story's output, which is why its `after` names them.
