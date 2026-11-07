---
id: S-0603
type: story
nature: improvement
title: agent.status keeps agents.json parsed and reads the conversations once a call, so it stops queueing on agentsMu and takes under 200 ms, not 0.8 to 2.4 s
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:40:19Z
updated: 2026-11-07T22:40:19Z
transitions: []
tags: [cli]
touches: [flai/internal/serve/agents.go, flai/cmd/serve_actions.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0603 agent.status keeps agents.json parsed and reads the conversations once a call, so it stops queueing on agentsMu and takes under 200 ms, not 0.8 to 2.4 s

## Goal

Cause 12 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `agent.status` is about 30% of flai serve's wall time, 84,417 slow calls in 4 days at p50 796 ms and p90 2,473 ms. Each call reads the host's whole `agents.json`, 663 KB for 14 projects, under the host-wide `agentsMu` (`flai/internal/serve/agents.go` `AgentStates`), so the 13 projects' calls queue: 8 at once take 435 ms each, 32 at once 1,185 ms. It then reads every conversation once for each live run, through `awaitingOther`, and again for `held`. The dashboard asks for it every 30 s for all 13 projects at once, and the board, the activity page, and a story's agent panel every 15 s. The host functions it calls take no context, so it records no phases.

## Acceptance criteria

- [ ] `AgentStates` keeps `agents.json` parsed, read again only when its size or modification time changed, and holds `agentsMu` only to read or replace what is kept.
- [ ] The activity a call builds reads the conversations once, whatever the number of live runs, and gives the same answer as before, which a test with three live runs shows.
- [ ] `agent.status` records its phases, and on this repository, 13 projects asked at once, each answers under 200 ms on the operator's host; server-performance.md records the before and after.

## Tasks

## Notes

- With the story for cause 13, agent.state's CPU falls from about 450 ms to under 60 ms.
- One agent state file per project instead of the host-wide one would change how flai keeps state, and is out of this epic.
