---
id: ADR-0044
title: A retry of a failed agent on a ready story with the in-progress limit full is queued until there is room
status: accepted
date: 2026-09-26
supersedes: []
superseded_by: []
refines: [ADR-0043]
---

# ADR-0044 A retry of a failed agent on a ready story with the in-progress limit full is queued until there is room

## Context

[ADR-0043](0043-flai-serve-starts-a-ready-story-s-agent-whenever-the-in-progress-limit-has-room.md) gave the operator a restart for a story whose agent dropped or failed: `flai serve agent restart` and the story page's **Restart agent** button. For a story in ready, the restart was refused while the in-progress limit was full. The operator could only try again later.

S-0118 turns the button into **Retry**, at the top right of the agent pane, hidden once pressed. It asks that Retry "queue the agent as it would normally do". On TH-0011 (2026-09-26) the operator chose to queue the agent rather than refuse it, and asked for a yellow dot while it waits.

## Decision

**A restart of a failed or dropped agent, for a story in ready while the in-progress limit is full, is queued instead of refused. flai serve starts the queued agent at its first look that finds room, as it starts a story that enters ready.** This refines ADR-0043's last refusal. Every other rule of ADR-0043 stands.

- **Queued.** The story's last run in `serve/agents.json` gets `queued`, the time the operator asked. No process starts. The command prints `queued another agent for S-…`, `--json` returns `queued`, and the journal records it.
- **Started when there is room.** The launcher does not hold back a queued run's story as one that has already had its agent since it entered ready. It starts that story in pull order, within the limit, and tells the new agent how the last one ended, as a restart does. The new run replaces the queued one.
- **Waiting, yellow.** While the story is in ready, `agent.status` reports a queued agent as `waiting` with why `queued: flai serve starts another agent when the in-progress limit has room`. The dashboard shows that as a yellow dot and does not offer Retry. Once the story leaves ready, the queued mark means nothing and the agent reads as failed again. For a story in progress, a retry starts the agent at once, whether or not one is queued.
- **Once.** A second retry while one is queued is refused, and the refusal says since when it has been queued.

## Consequences

- The operator presses Retry once, and the story gets its agent as soon as the limit allows. Nobody has to come back and try again.
- A queued story's dot is yellow, as for an agent waiting on an answer. The pane's line says which kind of waiting it is.
- `serve/agents.json` gains an optional field. An older flai ignores it, and its launcher treats a queued story as one that has already had its agent.

## Alternatives considered

- **Keep refusing while the limit is full.** This was offered on TH-0011 and declined. It leaves the operator to watch the board for room.
- **Start past the limit, as Start agent does.** A retry is not a first start the operator is forcing. Going past the limit would let failed stories crowd out the work in progress.
