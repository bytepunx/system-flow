---
id: S-0116
type: story
nature: remediation
title: when a story moves to ready, if there is available capacity, it should be assigned to an agent
status: review
parent: E-0003
owner: alex
created: 2026-09-24T08:37:20Z
updated: 2026-09-24T09:20:01Z
transitions:
  - to: ready
    at: 2026-09-24T08:37:26Z
    by: alex
  - to: in-progress
    at: 2026-09-24T08:44:21Z
    by: agent-S-0116
  - to: review
    at: 2026-09-24T09:20:01Z
    by: agent-S-0116
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd, flai/internal/serve, flai/internal/hostapi, flai/internal/harness, flai/internal/config, flai/internal/manifest, flai/internal/itemedit, design/adrs, design/system, design/issues, docs/operators, docs/users]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0116 when a story moves to ready, if there is available capacity, it should be assigned to an agent

## Goal

When stories are moved from the backlog into ready, if there is available capacity in the in process WIP, the story should be assigned to an agent.

## Acceptance criteria
- [x] a story moved from backlog to ready should be assigned to an agent if in process WIP isn't full
- [x] a story that had agent information changed in ready state should be assigned to an agent if in process WIP isn't full
- [x] someone attending the project no longer holds a ready story back while the in-progress limit has room, recorded in an ADR that supersedes ADR-0042's hold (TH-0010)
- [x] a story in ready or in-progress whose agent dropped or failed gets a new agent from `flai serve agent restart <story>` or the story page's **Restart agent** button, without re-entering ready (TH-0010)

## Tasks
- T-0411 A ready story whose agent is changed after its agent ended gets another
- T-0412 The story page says changing the agent of a ready story starts another
- T-0413 A story moved from backlog to ready with room in progress is started, per the designer's answer on TH-0010
- T-0414 Record the added case in the living design and the operator and user docs
- T-0415 An ADR: flai serve starts a ready story's agent whenever the in-progress limit has room, whoever is attending
- T-0416 The launcher starts ready stories with room whoever is attending, and attended_minutes is retired
- T-0417 flai serve agent restart <story> starts a new agent for a story whose agent dropped or failed
- T-0418 The story page has a Restart agent button for a story whose agent dropped or failed
- T-0419 Record the dropped hold, the retired setting, and restart in the design and the docs

## Notes

- Criterion 2 is met by T-0411 and pinned by `TestAChangedAgentStartsAReadyStoryAgain` in `flai/internal/serve/agents_test.go`. It covers a failed start, a changed harness, model, or config value, the in-progress limit, and a change made while the agent runs.
- Criterion 1 is pinned by `TestAStoryMovedFromBacklogToReadyGetsItsAgentWhenThereIsRoom`.
- Criterion 3 was decided on TH-0010 and recorded as ADR-0043. It is pinned by `TestSomeoneAttendingHoldsNothingBack`: a fresh MCP cursor and a narrative written just now, and the story is still started. `attended_minutes` is retired.
- Criterion 4 is pinned by `TestAStoryWhoseAgentDroppedOrFailedIsRestarted` (a failed agent, a dropped process, and each refusal), `TestServeAgentRestartSaysWhyItRefuses`, `TestARestartedAgentIsToldHowTheLastOneEnded`, the hostapi write tables, the route test in `flaiover/src/routes/api/items/[id]/agent/agent.test.ts`, and the button tests in `StoryAgent.svelte.test.ts`.
- Not tried live: the serving flai on this host is the installed 1.15.5, which has neither S-0114 nor this story. After acceptance, publish flai, run `flai host upgrade`, and try Restart agent on a failed story.
