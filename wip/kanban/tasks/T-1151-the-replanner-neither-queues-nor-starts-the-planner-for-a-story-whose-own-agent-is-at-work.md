---
id: T-1151
type: task
nature: remediation
title: The replanner neither queues nor starts the planner for a story whose own agent is at work
status: backlog
parent: S-0305
owner: alex
created: 2026-10-07T02:19:38Z
updated: 2026-10-07T02:19:38Z
transitions: []
stream: S-0305
tags: [serve, planner]
touches: [flai/internal/serve/replan.go, flai/internal/serve/replan_test.go]
---
# T-1151 The replanner neither queues nor starts the planner for a story whose own agent is at work

## Work

I-0098's instance: S-0300 was still `ready` when flai serve had already started its agent (it moved to `in-progress` at 23:09:18Z). The agent edited S-0300's touches at 23:03Z. `editTrigger` in `flai/internal/serve/replan.go` checks only that the story is in backlog or ready, so the edit queued the planner. `drain` starts a queued entry after `planCheck`, which does not refuse a story whose agent runs. A second run came at 23:10Z on the agent's criteria edit.

Change two places in `flai/internal/serve/replan.go`:

- `editTrigger` gives no trigger for an edit by the story's own agent, `agent-<ID>` for that story, as it already gives none for a planner's edit.
- `drain` drops a queued entry, logged as `queued planner dropped` with the reason, when its story is `in-progress` or `review`, or when the story's newest run in `AgentState.Stories` is `running()`. This also covers an entry queued before the agent started, and the schedule's and `planning.replan: agent`'s entries.

Leave `planCheck` and `serve.Plan` as they are: the operator may still ask for a story in progress to be planned.

Add tests to `flai/internal/serve/replan_test.go` that reproduce the instance:

- in `TestAnEditOfAPlannedStoryIsATrigger`, a case for an edit by `agent-S-0001` on a ready story that gives no trigger
- with `newReplanLab`, a queued entry for a ready story whose agent runs, and one for a story in progress, both dropped and logged, while another ready story's entry still starts its planner.

This task waits for no other task. It is layer 1.

## Done when

- [ ] An edit by a story's own agent queues no planner, and a test shows it.
- [ ] A queued planner for a story in progress or in review, or for a story whose agent runs, is dropped and logged, and a test shows it.
- [ ] A run the operator asks for with `flai plan` is unchanged.
- [ ] `flai test flai/internal/serve` and `flai check --strict` pass.

## Notes

Drafted by planner-S-0305 from I-0098's instance, `wip/agents/planner.md` (the runs at 23:08:09Z and 23:10:58Z), and S-0300's transitions.
