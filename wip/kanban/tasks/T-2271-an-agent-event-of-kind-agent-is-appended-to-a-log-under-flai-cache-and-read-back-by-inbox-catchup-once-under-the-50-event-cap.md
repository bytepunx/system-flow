---
id: T-2271
type: task
nature: improvement
title: An agent event of kind agent is appended to a log under .flai-cache and read back by inbox.CatchUp, once, under the 50-event cap
status: backlog
parent: S-0532
owner: alex
created: 2026-10-28T21:08:40Z
updated: 2026-10-28T21:08:40Z
transitions: []
stream: S-0532
tags: [flai, events, inbox]
touches: [flai/internal/inbox/agentevents.go, flai/internal/inbox/agentevents_test.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/serve/agent_events.go, flai/internal/serve/agent_events_test.go]
---
# T-2271 An agent event of kind agent is appended to a log under .flai-cache and read back by inbox.CatchUp, once, under the 50-event cap

## Work

Give `flai/internal/inbox` an event of kind `agent`, beside `edited`, `overlapped`, and `message` (`inbox.go`, the kind constants and `Event`). It carries the item it concerns (a story's ID, or the epic or story a planner run planned) or the role (`planner`, `orchestrator`, `analyzer`, `curator`, or a project's own role), `to` (`started`, `ended`, `failed`, or `stopped`), `why` (the journal's words: the exit, the operator's stop, the auto-restart count), `by` (the agent's name), and `at`.

Keep the events in an append-only JSONL log in the project's cache folder, `.flai-cache/agent-events.jsonl`, as `overlaps.jsonl` and `edits.jsonl` are kept (`itemedit/notice.go`). The MCP server runs in the agent's process and cannot read `~/.flai/serve`, so the log is the bridge. Put the writer and the reader in a new `agentevents.go`, with a query that lists one item's events since a time; the ADR-0108 thread task uses it.

`CatchUp` merges the log's events after the cursor with the others, sorted by time, under the 50 cap, each reported once. An agent is not told of its own runs, as `workitem.Changes` drops a move whose `by` is the agent. `Describe` gives each a summary such as `S-0532 agent failed: ended (exit 1) with S-0532 in in-progress`.

Add `flai/internal/serve/agent_events.go`, the one helper flai serve calls to turn an `AgentRun` and its moment into an event and append it. The story's agent task and the roles' task call it, and the usage-limit task changes only it, so those three can run together. This task waits for nothing: every other task of the story needs it.

## Done when

- [ ] An event appended to the log is reported once by `CatchUp`, after the cursor, sorted with the other kinds, within the 50 cap, and counted in `omitted` when it falls out
- [ ] An agent is not reported its own run's events
- [ ] `Describe` summarises each `to` with its item or role and its `why`
- [ ] The serve helper writes one event per call from an `AgentRun`, with the item or the role set as the run's kind says
- [ ] `flai test` passes on the changed paths

## Notes

Planned by planner-S-0532. The log's name and the own-run filter are the planner's assumptions, on the plan's thread.
