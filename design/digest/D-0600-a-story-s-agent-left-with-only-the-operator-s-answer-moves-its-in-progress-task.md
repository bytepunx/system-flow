---
id: D-0600
title: A story's agent left with only the operator's answer moves its in-progress task back to ready after opening the thread, so wait_for_events answers end and flai serve starts the session again on the answer
kind: decision
sources: [S-0659, TH-1103]
topics: [agents, workflow, all]
paths: [design/system/workflow.md, design/system/agent-narrative.md]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0600 A story's agent left with only the operator's answer moves its in-progress task back to ready after opening the thread, so wait_for_events answers end and flai serve starts the session again on the answer

S-0659's agent opened TH-1103 from T-2862, its last task, which needed the operator to confirm ADR-0292's choices before the ADR could be accepted. Nothing was left but the answer, so it moved T-2862 back to `ready`.

From S-0659's Decisions (`wip/archive/agents/S-0659.md`):

```text
T-2862 moved back to ready after opening TH-1103, so that the story holds no task in progress while only the operator's answer is left and flai serve can end and restart the session on it.
```

The reason is flai serve's rule in [workflow.md](../system/workflow.md#transitions-and-who-makes-them) § "An answered agent" and [agent-narrative.md](../system/agent-narrative.md) § "What an agent is told through MCP": `wait_for_events` answers `end: true` only when a thread on the story awaits an answer to the agent's own question and no task of the story is in progress. An agent with a task in progress keeps waiting, and pays for each wake. So a story's agent that has asked and has nothing else to do puts the waiting task back to `ready`, writes its Current state and Next steps with `stream_state`, and ends. On the answer it moves the task to `in-progress` again and finishes it.

Sources: S-0659 (`wip/archive/agents/S-0659.md`, Decisions); TH-1103.
