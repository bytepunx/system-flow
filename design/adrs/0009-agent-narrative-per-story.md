---
id: ADR-0009
title: One agent narrative per story in wip/agents
status: accepted
date: 2026-09-15
supersedes: []
superseded_by: []
---

# ADR-0009 One agent narrative per story in wip/agents

## Context

When an agent's environment crashes, the human should not have to reconstruct context. The recovery material must be small enough to read in a minute, current enough to trust, and tied to the unit of work that gets picked up.

## Decision

Each active story has one narrative file `wip/agents/<story-id>.md` with fixed sections: Context, Current state, Next steps, Decisions, Open questions, and an append-only Log. Tasks report into their story's narrative. An `index.md` lists active streams. The narrative is archived with its story. Agents are obliged by the baseline `CLAUDE.md` to keep Current state and Next steps true after every task transition.

## Consequences

- Recovery is: read index, read Current state and Next steps, reconcile with `git status`, continue.
- Narratives are curated summaries, not transcripts. Raw tool output and secrets are banned.
- One more file to keep honest per story, which `flai stream log` makes cheap.

## Alternatives considered

- One narrative per session: sessions do not map to work; a story may span several sessions and a session several stories.
- Narrative inside the story file: mixes human-facing work definition with agent state and bloats the kanban item.
