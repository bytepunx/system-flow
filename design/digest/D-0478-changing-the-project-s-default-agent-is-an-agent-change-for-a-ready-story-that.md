---
id: D-0478
title: Changing the project's default agent is an agent change for a ready story that names no harness, so flai serve's launcher starts it once more, since a run records the resolved agent
kind: decision
sources: [S-0573, T-2455]
topics: [flai, cli, go, code]
paths: [flai/internal/serve/agents.go, flai/internal/serve/start.go, flai/internal/workitem/agent.go]
date: 2026-11-03
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0478 Changing the project's default agent is an agent change for a ready story that names no harness, so flai serve's launcher starts it once more, since a run records the resolved agent

S-0573 made `flai serve` start a story whose agent names no harness with the project's default agent as it is at the start ([ADR-0251](../adrs/0251-a-story-whose-agent-names-no-harness-is-started-with-the-project-s-default.md)). T-2455 records that resolved agent, from `workitem.Repo.AgentOf`, as the run's `StoryAgent`, not the story's own empty block.

The launcher treats a changed agent as a new try ([ADR-0043](../adrs/0043-flai-serve-starts-a-ready-story-s-agent-whenever-the-in-progress-limit-has-room.md)). So setting or changing the default changes what such a story resolves to, and a ready story with no harness of its own, whose last run went on the host's command, is started once more. S-0573 held this to be what ADR-0251 intends; the design does not yet say it.

- S-0573: "T-2455 records the resolved agent as a run's `StoryAgent`. Setting a default therefore counts as an agent change for an agent-less ready story that last ran on the host's command, and the launcher starts it once more. That is what ADR-0251 intends."

Source: the Decisions of wip/archive/agents/S-0573.md; T-2455.
