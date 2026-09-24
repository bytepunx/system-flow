---
id: ADR-0043
title: "flai serve starts a ready story's agent whenever the in-progress limit has room, and the operator can restart one that dropped or failed"
status: accepted
date: 2026-09-24
supersedes: [ADR-0042]
superseded_by: []
refines: [ADR-0038, ADR-0041]
---

# ADR-0043 flai serve starts a ready story's agent whenever the in-progress limit has room, and the operator can restart one that dropped or failed

## Context

[ADR-0038](0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md) had `flai serve` start nothing while someone was attending the project. The idea was that an agent that is connected, or at work on a narrative, sees the ready story in its inbox and pulls it. [ADR-0042](0042-someone-attending-holds-a-ready-story-back-for-the-attended-window-then-flai.md) capped that hold at the attended window, `attended_minutes`, 6 by default. On 2026-09-24 the operator filed S-0116: a story moved to ready while the in-progress limit has room should get its agent. In TH-0010 they chose to drop the hold for stories that have room rather than keep the capped one.

A story also got one agent per entry into ready (ADR-0041). A story whose agent failed or dropped therefore had no way to get another short of re-entering ready. Ready accepts only in-progress and cancelled as next states, so for a story in progress that meant cancelling it and making a new one. The operator asked in the same thread for a restart command and button.

## Decision

**`flai serve` starts a ready story's agent whenever the in-progress limit has room, whoever is attending. The operator can restart the agent of a story in ready or in progress whose agent dropped or failed, from the command line or the story's page.** The operator decided this on TH-0010 (2026-09-24). This supersedes ADR-0042's hold and the attended rule of ADR-0038 and ADR-0041. It keeps ADR-0042's naming of the agent to the MCP server (`flai mcp --agent`).

- **No attendance hold.** The launcher no longer looks at MCP cursors or narratives. The in-progress limit is the only thing that holds back a ready story that can be started. `attended_minutes` controls nothing and is retired: `flai serve agent set --attended-minutes` is accepted, deprecated, and does nothing, and the dashboard's settings no longer offer it.
- **A changed agent is a new try.** Each run records the story's agent as it was started. A ready story whose agent says something else now, with no agent running, is started again within the limit. A run recorded before this carries no agent and counts as unchanged.
- **Restart.** `flai serve agent restart <story>`, and the **Restart agent** button that runs it through a hostapi write gated by the `agent` action and journalled, start a new agent in a new session for the story. The agent's name is unchanged. The run is recorded in `serve/agents.json` like the launcher's own, so the serving flai shows it, judges it when it ends, and resumes it on an answer. Restart refuses, saying why, when:
  - the `agent` action is off for the project;
  - the story is not in ready or in-progress;
  - an agent is running for it;
  - its agent is waiting on an answer, which restarts it by itself;
  - no agent was started for it by flai serve;
  - nothing can start it;
  - for a story in ready, the in-progress limit is full.

## Consequences

- A story that becomes ready while the limit has room gets its agent within a second of the change, as the operator asked.
- An agent holding `wait_for_work` and the launcher can both go for the same story. The second move to in-progress is refused. A flai-started agent that loses the race is told by its prompt to pull the next ready story or hold `wait_for_work`, and an interactive agent does the same.
- An operator working interactively no longer holds work back by being connected. An operator who wants no agents started turns the `agent` action off.
- A failed story no longer needs cancelling and recreating.

## Alternatives considered

- **Keep ADR-0042's capped hold.** Recommended on TH-0010 and declined by the operator: a story with room would still wait up to the window whenever anyone was connected.
- **Hold only while an attending agent is idle.** It needs the MCP server to publish which agents hold `wait_for_work`, a new contract, to save a race that the move to in-progress already settles.
- **Restart by moving the story back to ready.** The workflow does not allow in-progress to ready, and a move to restart an agent would put a false transition in the story's history.
