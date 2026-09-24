---
id: ADR-0041
title: A story in ready gets its agent whether it entered ready before or after flai serve started
status: accepted
date: 2026-09-24
supersedes: [ADR-0038]
superseded_by: []
---

# ADR-0041 A story in ready gets its agent whether it entered ready before or after flai serve started

## Context

[ADR-0038](0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md) says a story is started once each time it enters ready, and that "a restart of flai serve still starts nothing that was already ready". The launcher did this by retiring every story ready at its first look until it left ready. The rule dates from S-0079, when the operator started `flai serve` by hand and a restart was rare.

Since S-0106, `flai host` restarts `flai serve` on every upgrade, every host restart, and every crash. On 2026-09-24, S-0109 entered ready at 03:52Z; serve was restarted at 04:44Z and 05:37Z; neither serve started an agent for it or said why (S-0112). Each restart retired every ready story without a word.

Since S-0104, `serve/agents.json` keeps each story's newest run across a restart. The launcher already skips a story whose run started after it entered ready, and settles runs that outlived an earlier serve.

## Decision

**A story in ready gets its agent whether it entered ready before or after flai serve started. What stops a restart from starting a story twice is the story's run in `serve/agents.json`, not the time serve started.** This supersedes ADR-0038 in part: only the clause on restarts. The operator asked for it in S-0112 (2026-09-24).

- The first look after `flai serve` starts is like any other. It settles runs that outlived the earlier serve, then starts each ready story that has no agent running and none started since it entered ready. The in-progress limit and the attended rule still apply.
- Each ready story the launcher does not start, and that has no agent running, is named with its reason in `agent.status`'s `waiting`. The reasons are: its agent was already started since it entered ready (with how that run ended), it names no harness and no command is set, the `agent` action is off, someone is attending, or the in-progress limit is full. The launcher logs `agent not started` for the story each time its reason changes.

## Consequences

- An upgrade, host restart, or crash no longer strands ready stories. A story that entered ready while no serve ran is started at the next serve's first look.
- A restart can now start agents that nobody asked for since the restart. Each of those stories is in ready and has had no agent since it entered ready, so it was asked for when it was moved. Once the `agent` action is enabled, the operator should expect agents to start when a host comes up.
- A story whose run record was lost (`serve/agents.json` deleted) is started again after a restart. This is the only way a restart can start a story twice.
- A story never waits in ready without a stated reason.

## Alternatives considered

- **Keep the rule, and have `flai host` tell serve that a restart is its own.** This keeps the gap for a crash or a host restart, and it duplicates what the run record already knows.
- **Start after a restart only stories that entered ready while no serve ran.** This needs serve's stop time, which a crash does not record. The run record already covers the case.
