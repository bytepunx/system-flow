---
title: Board
updated: 2026-11-17
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: testing, title: Testing, kind: active, limit: true}
  - {name: review, kind: waiting, limit: true}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit], rules: [criteria-present, not-draft, after-done]}
  - {from: ready, to: in-progress, effects: [pull], rules: [not-blocked, not-held, task-after]}
  - {from: in-progress, to: testing, rules: [has-task]}
  - from: testing
    to: review
    rules: [no-open-question, clean-worktree, {rule: criteria-ticked, severity: warn}, kluster-verified]
  - {from: review, to: done, types: [epic, story], by: [acceptance], effects: [accept], rules: [children-closed, criteria-ticked]}
  - {from: review, to: done, types: [task], rules: [children-closed]}
  - {from: [in-progress, testing], to: done, types: [task], rules: [children-closed]}
  - {from: review, to: in-progress, effects: [send-back], rules: [reason-given, not-blocked, task-after]}
  - {from: testing, to: in-progress}
  - {from: in-progress, to: ready, rules: [criteria-present, not-draft, after-done]}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, in-progress, testing], to: cancelled, effects: [cancel], rules: [reason-given]}
  - {from: review, to: cancelled, by: [cascade], effects: [cancel], rules: [reason-given]}
  - {from: cancelled, to: backlog, effects: [reopen], rules: [parent-open]}
timers:
  - {name: lead, title: Lead time, from: created, to: completed}
  - {name: cycle, title: Cycle time, from: {first: in-progress}, to: completed}
  - {name: queue, title: Queue time, sum: [ready], until: {first: in-progress}}
  - {name: review, title: Review time, sum: [review]}
  - {name: testing-to-accept, title: Testing to acceptance, types: [story], from: {first: testing}, to: completed}
---

# Board

The template's baseline board with one state added between `in-progress` and `review`: `testing`, of the active kind, with a WIP limit. A story enters it from `in-progress` only once it has a task (`has-task`), and goes on to `review` under the rules the baseline gives the move from `in-progress` to `review`, which this board no longer has. The timer `testing-to-accept` measures a story from its first move into `testing` to its acceptance.

S-0756's end-to-end test, `flai/cmd/workflow_e2e_test.go`, renders the template into a new project, writes this file over its `design/workflows/board.md`, and checks that every surface of flai follows the added state.
