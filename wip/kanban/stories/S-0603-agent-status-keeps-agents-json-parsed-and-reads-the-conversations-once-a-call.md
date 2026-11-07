---
id: S-0603
type: story
nature: improvement
title: agent.status keeps agents.json parsed and reads the conversations once a call, so it stops queueing on agentsMu and takes under 200 ms, not 0.8 to 2.4 s
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:40:19Z
updated: 2026-11-07T23:49:36Z
transitions: []
tags: [cli]
touches: [flai/internal/serve/agents.go, flai/cmd/serve_actions.go, flai/internal/serve/agents_test.go, flai/internal/messages/messages.go, flai/internal/workitem/hold.go, flai/internal/workitem/share.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/writes_test.go, flai/cmd/serve_actions_test.go, design/system/server-performance.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
finalized:
  by: alex
  at: 2026-11-07T23:49:36Z
---
# S-0603 agent.status keeps agents.json parsed and reads the conversations once a call, so it stops queueing on agentsMu and takes under 200 ms, not 0.8 to 2.4 s

## Goal

Cause 12 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `agent.status` is about 30% of flai serve's wall time, 84,417 slow calls in 4 days at p50 796 ms and p90 2,473 ms. Each call reads the host's whole `agents.json`, 663 KB for 14 projects, under the host-wide `agentsMu` (`flai/internal/serve/agents.go` `AgentStates`), so the 13 projects' calls queue: 8 at once take 435 ms each, 32 at once 1,185 ms. It then reads every conversation once for each live run, through `awaitingOther`, and again for `held`. The dashboard asks for it every 30 s for all 13 projects at once, and the board, the activity page, and a story's agent panel every 15 s. The host functions it calls take no context, so it records no phases.

## Acceptance criteria

- [ ] `AgentStates` keeps `agents.json` parsed, read again only when its size or modification time changed, and holds `agentsMu` only to read or replace what is kept.
- [ ] The activity a call builds reads the conversations once, whatever the number of live runs, and gives the same answer as before, which a test with three live runs shows.
- [ ] `agent.status` records its phases, and on this repository, 13 projects asked at once, each answers under 200 ms on the operator's host; server-performance.md records the before and after.

## Tasks
- T-2621 AgentStates keeps agents.json parsed, read again only when its size or modification time changed, and holds agentsMu only to read or replace it
- T-2623 The activity agent.status builds reads the conversations once, for every live run and for held, with the same answer as before
- T-2629 agent.status passes its context to the host's Agent and records the phases of the state it builds
- T-2638 Measure agent.status with 13 projects asked at once on the operator's host, before and after, and record it in server-performance.md

## Notes

- With the story for cause 13, agent.state's CPU falls from about 450 ms to under 60 ms.
- One agent state file per project instead of the host-wide one would change how flai keeps state, and is out of this epic.
