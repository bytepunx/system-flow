---
id: S-0603
type: story
nature: improvement
title: agent.status keeps agents.json parsed and reads the conversations once a call, so it stops queueing on agentsMu and takes under 200 ms, not 0.8 to 2.4 s
status: in-progress
parent: E-0033
owner: alex
created: 2026-11-07T22:40:19Z
updated: 2026-11-08T03:31:16Z
transitions:
  - to: ready
    at: 2026-11-08T00:11:17Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-08T01:12:15Z
    by: agent-S-0603
blocked:
  - from: 2026-11-08T02:23:41Z
    until: 2026-11-08T02:31:46Z
    reason: awaiting the operator on TH-0911
  - from: 2026-11-08T02:38:45Z
    until: 2026-11-08T03:23:11Z
    reason: awaiting the operator on TH-0911
  - from: 2026-11-08T03:24:01Z
    until: 2026-11-08T03:30:55Z
    reason: awaiting the operator on TH-0911
  - from: 2026-11-08T03:31:16Z
    reason: awaiting the operator on TH-0911
tags: [cli]
touches: [flai/internal/serve/agents.go, flai/cmd/serve_actions.go, flai/internal/serve/agents_test.go, flai/internal/messages/messages.go, flai/internal/workitem/hold.go, flai/internal/workitem/share.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/writes_test.go, flai/cmd/serve_actions_test.go, design/system/server-performance.md, design/issues/I-0313-agent-status-with-13-projects-at-once-is-not-under-200-ms.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4761
  priced_by: harness
  turns:
    - day: 2026-11-08
      ceremony: 2
      hand_edits: 3
      work: 97
  models:
    - model: claude-opus-5-5
      input: 396
      output: 101915
      cache_read: 56780133
      cache_write: 1930944
      cost: 24.1995
  strategic:
    - kind: orchestrator
      seconds: 139
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 22
          output: 161
          cache_read: 875190
          cache_write: 69806
          cost: 0.2421
cost_of_delay:
  inputs:
    penalty_per_week: 100
    by: orchestrator
    at: 2026-11-07T23:51:56Z
  value: 100
  by: planner-S-0603
  at: 2026-11-07T23:52:33Z
forecast:
  duration: 1h
  delivery: 2026-11-08T02:03:00Z
  basis: "Its own forecast of 1h; 2nd in the pull order with an in-progress limit of 5, behind S-0598, S-0599, S-0600, S-0605, S-0606 and S-0601."
  by: flai
  at: 2026-11-08T00:25:18Z
finalized:
  by: alex
  at: 2026-11-07T23:49:36Z
---
# S-0603 agent.status keeps agents.json parsed and reads the conversations once a call, so it stops queueing on agentsMu and takes under 200 ms, not 0.8 to 2.4 s

## Goal

Cause 12 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `agent.status` is about 30% of flai serve's wall time, 84,417 slow calls in 4 days at p50 796 ms and p90 2,473 ms. Each call reads the host's whole `agents.json`, 663 KB for 14 projects, under the host-wide `agentsMu` (`flai/internal/serve/agents.go` `AgentStates`), so the 13 projects' calls queue: 8 at once take 435 ms each, 32 at once 1,185 ms. It then reads every conversation once for each live run, through `awaitingOther`, and again for `held`. The dashboard asks for it every 30 s for all 13 projects at once, and the board, the activity page, and a story's agent panel every 15 s. The host functions it calls take no context, so it records no phases.

## Acceptance criteria

- [x] `AgentStates` keeps `agents.json` parsed, read again only when its size or modification time changed, and holds `agentsMu` only to read or replace what is kept.
- [x] The activity a call builds reads the conversations once, whatever the number of live runs, and gives the same answer as before, which a test with three live runs shows.
- [ ] `agent.status` records its phases, and on this repository, 13 projects asked at once, each answers under 200 ms on the operator's host; server-performance.md records the before and after.

## Tasks
- T-2621 AgentStates keeps agents.json parsed, read again only when its size or modification time changed, and holds agentsMu only to read or replace it
- T-2623 The activity agent.status builds reads the conversations once, for every live run and for held, with the same answer as before
- T-2629 agent.status passes its context to the host's Agent and records the phases of the state it builds
- T-2638 Measure agent.status with 13 projects asked at once on the operator's host, before and after, and record it in server-performance.md

## Notes

- With the story for cause 13, agent.state's CPU falls from about 450 ms to under 60 ms.
- One agent state file per project instead of the host-wide one would change how flai keeps state, and is out of this epic.

### Planning

Touches, file by file, no folder:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/serve/agents.go` | declared | `AgentStates`, `agentsMu`, `updateAgent`, `Activity`, `runActivity`, `held`, `awaitingOther` |
| `flai/cmd/serve_actions.go` | declared | `host().Agent`, which builds the `agent.status` state; its phases |
| `flai/internal/serve/agents_test.go` | co-change (34%) | the kept parse and the three-live-run test |
| `flai/cmd/serve_actions_test.go` | co-change (41%) | the phases test |
| `flai/internal/hostapi/writes.go` | co-change (23%), layout | `Host.Agent` takes a context |
| `flai/internal/hostapi/writes_test.go` | co-change (20%), layout | its stub of `Host.Agent` |
| `flai/internal/hostapi/hostapi.go` | layout | the `agent.status` method passes its context |
| `flai/internal/messages/messages.go` | layout | `AwaitingOther` over a list read once |
| `flai/internal/workitem/hold.go` | layout | `Repo.Holds` from shares and asks already read |
| `flai/internal/workitem/share.go` | layout | `Repo.Conversations`, which `Holds` reads today |
| `design/system/server-performance.md` | design | criterion 3: the before and after |

Left out of `flai touches suggest`: `docs/users/flai-reference.md`, `docs/operators/settings.md`, `design/system/flai-cli.md`, and the config, harness, restart, and start files. The story adds no command, flag, or setting, and changes no agent start.

No `after`: the story reads no other story's output. S-0602 (cause 13) also changes how conversations are read, in `flai/internal/messages/messages.go`; the trial merge judges that overlap.

Forecast: `flai forecast` gave 22m, 94 s per unit of size over 124 done improvement stories on claude-opus-5-5, times size 14. Raised to 1h. The median is mostly smaller changes. This one changes a host-wide lock with concurrency tests across three packages, and measures 13 concurrent calls on the loaded host. Its close-out runs the integration and `flaiover-flai` tiers, since it changes `flai/internal/hostapi/`. The delivery moves the playout's start, 00:23, on by the hour.

Cost of delay: 100 USD a week, as `flai cod` works it out from the one input, `penalty_per_week: 100`, which the orchestrator set at 2026-11-07T23:51:56Z from TH-0888's table. Kept as given: no revenue rides on the story, and nothing else is lost each cycle.
