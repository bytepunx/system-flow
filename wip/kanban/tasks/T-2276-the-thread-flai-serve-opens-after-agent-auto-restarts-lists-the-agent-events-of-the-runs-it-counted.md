---
id: T-2276
type: task
nature: improvement
title: The thread flai serve opens after agent.auto_restarts lists the agent events of the runs it counted
status: backlog
parent: S-0532
owner: alex
created: 2026-10-28T21:09:33Z
updated: 2026-10-29T02:00:45Z
transitions: []
stream: S-0532
tags: [flai, serve, events]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
after: [T-2273]
---
# T-2276 The thread flai serve opens after agent.auto_restarts lists the agent events of the runs it counted

## Work

`stopRestarting()` in `flai/internal/serve/agents.go` opens the ADR-0108 thread when a story's agent ends past `agent.auto_restarts`. Today its text says how many times the agent ended and how the last run ended. Add a list of the story's `agent` events since the run the count started from, oldest first, one line each with its time in UTC, `to`, and `why`, read with the query T-2271 adds to the event log. A started event and its end make two lines, so the operator reads each restart and how it ended. Keep the thread's other lines as they are.

This task waits for T-2272, which writes the story's agent's events and changes `agents.go` too, so the two cannot run together.

## Done when

- [ ] A test ends a story's agent past the limit and reads the thread's text, which lists each `started` and `failed` event of the counted runs with its time and `why`, in order
- [ ] The thread's first lines and its advice on `flai serve agent restart` and Retry are unchanged
- [ ] `flai test` passes on the changed paths

## Notes

Planned by planner-S-0532.
