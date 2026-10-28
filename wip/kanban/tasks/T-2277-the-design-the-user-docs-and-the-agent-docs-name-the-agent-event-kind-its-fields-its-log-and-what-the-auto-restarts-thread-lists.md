---
id: T-2277
type: task
nature: improvement
title: The design, the user docs, and the agent docs name the agent event kind, its fields, its log, and what the auto_restarts thread lists
status: backlog
parent: S-0532
owner: alex
created: 2026-10-28T21:09:49Z
updated: 2026-10-28T21:09:49Z
transitions: []
stream: S-0532
tags: [flai, docs, events]
touches: [design/system/flai-cli.md, design/system/strategic-agents.md, design/system/agent-narrative.md, docs/users/flai.md, docs/operators/settings.md]
after: [T-2273, T-2274, T-2275, T-2276]
---
# T-2277 The design, the user docs, and the agent docs name the agent event kind, its fields, its log, and what the auto_restarts thread lists

## Work

Describe what the code tasks built, citing ADR-0209 §6:

- `design/system/flai-cli.md`: where `flai serve`'s agent runs are described, and its internal structure beside `internal/inbox`, the `agent` event, its four values of `to`, its fields, and its log `.flai-cache/agent-events.jsonl`; and that the ADR-0108 thread lists the counted runs' events.
- `design/system/strategic-agents.md`: in each role's "run and how it ended" section, that the run's start and end are `agent` events the orchestrator, the analyzer, and the operator see in `inbox` and `wait_for_events`.
- `design/system/agent-narrative.md` § What an agent is told through MCP, the agent docs: the `agent` kind beside `edited`, `overlapped`, and `message`, with what an agent does with one (nothing is owed; it tells of an agent failing or restarting).
- `docs/users/flai.md`: the MCP tool table's `inbox` and `wait_for_events` rows, and the "since this agent last looked" paragraph, name the kind and its fields; the `agent.auto_restarts` row says what the thread now lists.
- `docs/operators/settings.md`: the `agent.auto_restarts` entry, the same.

Bump each document's `updated`. This task waits for T-2273, T-2274, T-2275, and T-2276, so that it describes what was built, and T-2276 waits for T-2272.

## Done when

- [ ] Each of the five documents names the `agent` kind, `to`'s four values, the item or role, and `why`, where it describes events or the auto-restart thread
- [ ] `flai check --strict` and the markdown lint are clean on the changed documents

## Notes

Planned by planner-S-0532. No convention changes: an `agent` event asks nothing of a story's agent, so `work-management.md` and its template copy stay as they are, as the plan's thread says. `docs/users/flai-reference.md` is generated from the commands' help, which this story does not change.
