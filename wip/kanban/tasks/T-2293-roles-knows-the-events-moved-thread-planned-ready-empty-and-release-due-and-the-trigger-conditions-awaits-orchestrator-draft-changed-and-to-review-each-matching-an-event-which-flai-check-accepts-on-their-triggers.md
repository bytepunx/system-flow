---
id: T-2293
type: task
nature: feature
title: roles knows the events moved, thread, planned, ready-empty, and release-due, and the trigger conditions awaits-orchestrator, draft-changed, and to-review, each matching an event, which flai check accepts on their triggers
status: done
parent: S-0535
owner: alex
created: 2026-10-28T21:37:56Z
updated: 2026-10-28T22:44:59Z
transitions:
  - to: ready
    at: 2026-10-28T22:31:53Z
    by: system-flow
  - to: in-progress
    at: 2026-10-28T22:31:54Z
    by: system-flow
  - to: done
    at: 2026-10-28T22:44:59Z
    by: agent-S-0535
stream: S-0535
tags: [flai, events, roles, check]
touches: [flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/check/roles_test.go]
after: [T-2291]
usage:
  source: log
  seconds: 774
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 47
      output: 210
      cache_read: 1794263
      cache_write: 171458
      cost: 0.6683
---
# T-2293 roles knows the events moved, thread, planned, ready-empty, and release-due, and the trigger conditions awaits-orchestrator, draft-changed, and to-review, each matching an event, which flai check accepts on their triggers

## Work

`roles.Events` (`flai/internal/roles/roles.go`) lists only ADR-0207's seven events, and `triggers()` reports any other as `roles.file`. The trigger conditions in `conditions.go` are `stale-figures`, `replan-agent`, and `each-ready-story`. `QualifiesTrigger` matches an event name only, and nothing matches a condition against an event that happened.

- Add `moved`, `thread`, `planned`, `ready-empty`, and `release-due` to `roles.Events`, keeping `review` and `thread-opened` unless T-2291's ADR says otherwise.
- Add the trigger conditions `awaits-orchestrator` (qualifies `{event: thread}`), `draft-changed` (qualifies `{event: edited}`), and `to-review` (qualifies `{event: moved}`). Let `to-backlog` qualify `{event: moved}` beside its calls, as T-2291's ADR decides, so S-0536's trigger list passes.
- Give each trigger condition a match against an event that happened, in a small type of `roles`' own (kind, `to`, whether the edited item is a draft story, whom a thread awaits), so that `roles` imports nothing of `inbox`:
  - `draft-changed`: an `edited` event on a story that is a draft after the edit (ADR-0211 §3).
  - `to-review`: a `moved` event to `review`.
  - `to-backlog`: a `moved` event to `backlog`.
  - `awaits-orchestrator`: a `thread` event that awaits the orchestrator or the operator (ADR-0181).
- Test each condition's match both ways, and each name in `roles_test.go` and `conditions_test.go`.
- Add rows to `check/roles_test.go`: each new condition is accepted on its event, and refused on another event.

It waits for T-2291, whose ADR settles the names. It shares no path with T-2292, so the two run together.

## Done when

- A role file whose trigger names each new event, and each new condition on its event, passes `flai check`, with a test.
- The same condition on the wrong event is `roles.file`, with a test.
- Each condition's match is tested on a matching and a non-matching event.
- `flai test` on the changed paths passes.

## Notes

S-0536 reads the matches when flai serve queues the orchestrator's reasons; this task only implements and tests them.
