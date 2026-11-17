---
title: Board
updated: 2026-11-17
status: active
states:
  - {name: backlog, kind: queue, initial: true, order: true}
  - {name: ready, kind: queue, order: true, limit: true}
  - {name: in-progress, kind: active, limit: true}
  - {name: review, kind: waiting, limit: true, runs: {kluster: [deploy, checks]}}
  - {name: done, kind: accepted}
  - {name: cancelled, kind: cancelled}
transitions:
  - {from: backlog, to: ready, effects: [commit], rules: [criteria-present, not-draft, after-done]}
  - {from: ready, to: in-progress, effects: [pull], rules: [not-blocked, not-held, task-after]}
  - from: in-progress
    to: review
    rules: [has-task, no-open-question, clean-worktree, {rule: criteria-ticked, severity: warn}, kluster-verified]
    verify: {suites: all, severity: warn}
  - {from: review, to: done, types: [epic, story], by: [acceptance], effects: [accept], rules: [children-closed, criteria-ticked], verify: all}
  - {from: review, to: done, types: [task], rules: [children-closed]}
  - {from: in-progress, to: done, types: [task], rules: [children-closed]}
  - {from: review, to: in-progress, effects: [send-back], rules: [reason-given, not-blocked, task-after]}
  - {from: in-progress, to: ready, rules: [criteria-present, not-draft, after-done]}
  - {from: ready, to: backlog}
  - {from: [backlog, ready, in-progress], to: cancelled, effects: [cancel], rules: [reason-given]}
  - {from: review, to: cancelled, by: [cascade], effects: [cancel], rules: [reason-given]}
  - {from: cancelled, to: backlog, effects: [reopen], rules: [parent-open]}
verify_at: {batch: all, ci: all}
timers:
  - {name: lead, title: Lead time, from: created, to: completed}
  - {name: cycle, title: Cycle time, from: {first: in-progress}, to: completed}
  - {name: queue, title: Queue time, sum: [ready], until: {first: in-progress}}
  - {name: review, title: Review time, sum: [review]}
---

# Board

This is the baseline workflow flai carries in its binary: today's board, with its six states, the transitions between them, the rules each move passes, and the timers flai measures over every item's history, as ADR-0301 states it. A project with no `design/workflows/board.md` of its own runs it.

`verify` on a transition names the test suites a move by it needs its last `flai verify` to have passed, `all` for every suite, and `warn` lets the move go with a warning; `runs` on a state names what flai runs when a story enters it, suites and the cluster's deploy and checks; `verify_at` names the suites the merge queue's batch, the publish, and CI run. The values here are today's behaviour, as ADR-0312 §8 states it, and nothing reads them to decide a run yet.

flai reads only the front matter. This body is prose for the operator: why the board is shaped as it is, and what changing it means.
