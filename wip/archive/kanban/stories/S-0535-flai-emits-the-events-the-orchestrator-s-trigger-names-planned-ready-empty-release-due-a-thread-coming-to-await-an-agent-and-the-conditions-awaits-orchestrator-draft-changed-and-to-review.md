---
id: S-0535
type: story
nature: feature
title: "flai emits the events the orchestrator's trigger names: planned, ready-empty, release-due, a thread coming to await an agent, and the conditions awaits-orchestrator, draft-changed, and to-review"
status: done
parent: E-0030
owner: alex
created: 2026-10-28T21:02:24Z
updated: 2026-10-29T01:07:13Z
transitions:
  - to: ready
    at: 2026-10-28T22:21:19Z
    by: alex
  - to: in-progress
    at: 2026-10-28T22:28:24Z
    by: agent-S-0535
  - to: review
    at: 2026-10-29T00:37:08Z
    by: system-flow
  - to: done
    at: 2026-10-29T01:07:13Z
    by: orchestrator
tags: [flai, serve, events]
topics: [orchestration, cli, mcp, template]
touches: [flai/internal/serve, flai/internal/mcpserver, design/system/flai-cli.md, design/workflows/roles/README.md, design/adrs, design/adrs/README.md, flai/internal/inbox/events.go, flai/internal/inbox/events_test.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/inbox/threadevents.go, flai/internal/inbox/threadevents_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/check/roles_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/readyempty.go, flai/internal/serve/readyempty_test.go, flai/internal/serve/agents.go, flai/cmd/accept.go, flai/cmd/accept_release_due_test.go, flai/internal/release/evaluate.go, flai/internal/release/evaluate_test.go, template/root/design/workflows/roles/README.md, template/CHANGELOG.md, design/system/agent-narrative.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/accept_test.go, design/issues/I-0227-flai-test-s-go-test-tier-runs-the-mcpserver-package-with-race-into-go-test-s-10-minute-default-timeout-on-a-loaded-host-and-names-only-the-package-so-a-story-that-changes-mcpserver-cannot-pass-its-tests.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 7739
  priced_by: harness
  turns:
    - day: 2026-10-28
      ceremony: 2
      test_runs: 3
      hand_edits: 2
      work: 64
    - day: 2026-10-29
      ceremony: 1
      test_runs: 1
      hand_edits: 1
      work: 21
  models:
    - model: claude-opus-5-5
      input: 692
      output: 234337
      cache_read: 43846036
      cache_write: 1941344
      cost: 24.3762
  strategic:
    - kind: orchestrator
      seconds: 3676
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 480
          output: 8305
          cache_read: 77036466
          cache_write: 102425
          cost: 16.9749
        - model: claude-sonnet-5-5
          input: 18
          output: 94
          cache_read: 385099
          cache_write: 70002
          cost: 0.3149
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: orchestrator
    at: 2026-10-28T21:36:58Z
  value: 37.5
  by: planner-S-0535
  at: 2026-10-28T21:39:56Z
forecast:
  duration: 4h
  delivery: 2026-10-29T03:04:00Z
  basis: "Its own forecast of 4h; 2nd in the pull order with an in-progress limit of 5, behind S-0416 and S-0532."
  by: flai
  at: 2026-10-28T22:22:16Z
---
# S-0535 flai emits the events the orchestrator's trigger names: planned, ready-empty, release-due, a thread coming to await an agent, and the conditions awaits-orchestrator, draft-changed, and to-review

## Goal

ADR-0209 §1 lists the reasons the orchestrator is started for in ADR-0207's trigger vocabulary. Some are events flai emits today with a condition to name (`moved` to review, `edited` on a draft, `reordered`, `accepted`); others are not events at all: a planner's run ending with its items written, the ready column draining below the in-progress limit, the release policy's figures changing with an acceptance, a thread coming to await the orchestrator or the operator (ADR-0181). This story defines each as an event or a named condition flai implements (ADR-0207 §3), emits it from where it happens, and reports it through `wait_for_events` and `inbox` like the others, so that the following story can name them in `orchestrator.md`'s `trigger` and flai serve can queue them.

## Acceptance criteria
- [x] `planned` is an event, with the item planned and the IDs written, emitted when a planner's run ends with its items written, and reported by `wait_for_events` and `inbox`
- [x] `ready-empty` is an event emitted when the ready column drains below the in-progress limit's room, once per draining
- [x] `release-due` is an event emitted when an acceptance changes what `flai release --evaluate` would answer
- [x] `awaits-orchestrator`, `draft-changed`, and `to-review` are conditions `flai check` accepts on `thread`, `edited`, and `moved` triggers, each implemented and tested
- [x] A `thread` event exists for a thread coming to await an agent, with the thread and whom it awaits, under ADR-0181's rule
- [x] The design (`flai-cli.md`'s event table, `design/workflows/roles/README.md`'s trigger table) and the MCP descriptions name each event and condition

## Tasks
- T-2291 An ADR records the trigger vocabulary S-0535 adds: the events moved, thread, planned, ready-empty, and release-due, and the conditions awaits-orchestrator, draft-changed, and to-review
- T-2292 flai appends the events it emits to a log under .flai-cache, and inbox.CatchUp reads them back once, described, under the 50-event cap, waking wait_for_events
- T-2293 roles knows the events moved, thread, planned, ready-empty, and release-due, and the trigger conditions awaits-orchestrator, draft-changed, and to-review, each matching an event, which flai check accepts on their triggers
- T-2294 A planner's run that ends with items written appends a planned event naming the item planned and the IDs written
- T-2295 flai serve appends one ready-empty event when the ready column drains below the in-progress limit's room, and not again until ready reaches the room
- T-2296 flai accept appends a release-due event when the acceptance changes what flai release --evaluate answers
- T-2297 inbox.CatchUp reports a thread event when a thread comes to await the reader, or for the orchestrator the operator, naming the thread, its item, and whom it awaits under ADR-0181
- T-2298 The wait_for_events and inbox descriptions and the server's instructions name planned, ready-empty, release-due, and thread, and a test reads each back through both tools
- T-2299 flai-cli.md gains an event table, the roles README and its template copy list the new trigger events and conditions, and the agent and user docs name the new event kinds

## Notes

A named condition is one flai implements, not one a role file defines (ADR-0207 §3). The names are ADR-0209's proposals; the story may rename them with the ADR's text if a better name fits, before the ADR is accepted.

### Planning

Planned by planner-S-0535 on 2026-10-28. The tasks, their layers, and the assumptions are on the plan's thread.

The design: `wait_for_events` and `inbox` read `inbox.CatchUp`, which derives events from the work item files (`workitem.Changes`), from two logs under `.flai-cache` (`edits.jsonl`, `overlaps.jsonl`), from messages, and from the board order. `planned`, `ready-empty`, and `release-due` happen in flai serve and in `flai accept`, so they get a log of their own (T-2292). The `thread` event is derived from the thread files, as messages are (T-2297). ADR-0207's `roles.Events` lacks `moved`, `thread`, and the three new events, and ADR-0209 is now accepted, so T-2291 records the vocabulary in an ADR of its own before the code.

Touches, and where each came from:

- `flai/internal/serve`, `flai/internal/mcpserver`: declared, and kept. The tasks name every file below them the plan expects, so the story's claim narrows each folder to those files (ADR-0096).
- `flai/internal/workitem`: declared, and kept. No task names a file under it: the new kinds live in `inbox`, and `workitem.Changes` already reports `moved` and `edited`. The story's agent may drop it with `flai touches` when it reviews the tasks.
- `design/adrs`: declared, and kept, because T-2291 adds an ADR whose number `flai adr new` allocates only then, a file no task can name yet. `design/adrs/README.md` is its index.
- `design/system/flai-cli.md`, `design/workflows/roles/README.md`: declared; criterion 6.
- `flai/internal/inbox/inbox.go`, `inbox_test.go`: layout; `CatchUp`, the kinds, `Describe`, and the 50 cap.
- `flai/internal/inbox/events.go`, `events_test.go`, `threadevents.go`, `threadevents_test.go`: layout, new files for the log and the thread event.
- `flai/internal/mcpserver/server.go`, `server_test.go`, `folder.go`, `folder_test.go`: layout; `noticeLogs()`, the tool descriptions, the instructions, and the read-back tests.
- `flai/internal/roles/roles.go`, `roles_test.go`, `conditions.go`, `conditions_test.go`, `flai/internal/check/roles_test.go`: layout; `roles.Events`, the trigger conditions, and the `roles.file` check (criterion 4).
- `flai/internal/serve/plan.go`, `plan_test.go`: layout; `planEnded` (criterion 1).
- `flai/internal/serve/agents.go`, `readyempty.go`, `readyempty_test.go`: layout; the launcher's `look()` sees the ready column and the room (criterion 2).
- `flai/cmd/accept.go`, `accept_release_due_test.go`, `flai/internal/release/evaluate.go`, `evaluate_test.go`: layout; every acceptance passes `acceptItem`, and `release.EvaluateRepo` gives the answer (criterion 3).
- `template/root/design/workflows/roles/README.md`, `template/CHANGELOG.md`: design; the template's copy of the roles README changes with this one.
- `design/system/agent-narrative.md`, `docs/users/flai.md`: co-change and design; the agent's and the user's description of the MCP events.
- Left out of the suggestion: `flai-reference.md`, `docs/operators/*`, `flaiover-dashboard.md`, `workflow.md`, `work-management.md`, `project-manifest.md`, the host API, `serve_actions.go`, `root.go`, and `check.go`. They co-change with past serve stories, but nothing in them changes here: no setting, command, or check rule is added. `strategic-agents.md` is S-0536's, which changes the orchestrator's trigger.

Forecast: `flai forecast` gave 32m (median 159 s per unit of size over 4 done feature stories on claude-opus-5-5, size 12: 6 criteria, 6 touches). Raised to 4h: nine tasks in five layers across the inbox, roles and check, three writers in serve and accept, the MCP tools, an ADR, and five documents, beside S-0532's 3h for seven tasks of a similar spread. Delivery is flai's 2026-10-29T13:57Z moved by the added 3h28m, 39th in the pull order with an in-progress limit of 5.

Cost of delay: `flai cod` gave 37.50 USD a week from the input the orchestrator set on TH-0714 (15m of time lost per 168h cycle at 150 USD an hour, 1.00 cycles a week). The value stands as computed; the input is not changed.

### Accepted by the orchestrator

- Verified: cfd998e467399780aa040f64c7ab027eabd71d67
- Batch: 20261029-005305
- Enqueued verified: d3b740cc839d38ce149806d2d8e9aaa26f706635
- At: 2026-10-29T01:07:13Z

Verdict: meets all six criteria at d3b740cc; flai verify passed, with integration and smoke left to the queue; no convention findings and no file outside the touches.
- 1: flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/inbox/events.go, flai/internal/inbox/events_test.go, flai/internal/inbox/inbox.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go
- 2: flai/internal/serve/readyempty.go, flai/internal/serve/readyempty_test.go, flai/internal/serve/agents.go
- 3: flai/cmd/accept.go, flai/cmd/accept_release_due_test.go, flai/cmd/accept_test.go, flai/internal/release/evaluate.go, flai/internal/release/evaluate_test.go
- 4: flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/check/roles_test.go
- 5: flai/internal/inbox/threadevents.go, flai/internal/inbox/threadevents_test.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/mcpserver/cursor.go, flai/internal/mcpserver/server_test.go
- 6: design/system/flai-cli.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, template/CHANGELOG.md, design/system/agent-narrative.md, docs/users/flai.md, docs/users/flai-reference.md, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go
