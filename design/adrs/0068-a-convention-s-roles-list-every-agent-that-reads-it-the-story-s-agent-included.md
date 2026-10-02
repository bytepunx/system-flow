---
id: ADR-0068
title: "A convention's roles list every agent that reads it, the story's agent included, and a convention without roles is read by all"
status: accepted
date: 2026-10-02
supersedes: []
superseded_by: []
refines: [ADR-0059]
---

# ADR-0068 A convention's roles list every agent that reads it, the story's agent included, and a convention without roles is read by all

## Context

ADR-0059 gave conventions a front matter key, `roles`, listing the sub-agents that read a file as well as the story's agent: `explore`, `verify`, or both. A convention without it is read by the story's agent only, and the story's agent reads every convention whatever its roles. So a convention the sub-agents need must name them, and nothing can say that a convention is for the story's agent alone or for the sub-agents alone. The designer decided on 2026-10-02 that a convention without `roles` is read by every agent, and that the story's agent no longer reads everything.

## Decision

A convention's `roles` lists every agent that reads it, the story's agent included: `story` for the agent working a story, `explore` for the explorer, and `verify` for the verifier. A convention without `roles` is read by all of them. `flai prime --story S-nnnn` prints the conventions whose roles are empty or list `story`, and `flai prime --story S-nnnn --role explore|verify` those whose roles are empty or list the role; the story's topics apply to both as before. `flai prime --cat` still prints every convention. `flai check` accepts `story`, `explore`, and `verify`, and warns about any other value. This refines ADR-0059's priming; the rest of ADR-0059 stands.

## Consequences

- The baseline's `roles` are rewritten in the same change as the code: a file the story's agent reads lists `story` or no roles, so the story's agent reads what it reads today unless the baseline deliberately leaves a file to the sub-agents.
- A convention that only the story's agent needs, such as work management, says `roles: [story]` and stops costing sub-agents' packs; one every agent needs, such as safety, carries no roles.
- A project that lists `explore` or `verify` without `story` on a file its story's agent needs must add `story` when it upgrades; `flai upgrade` replaces the baseline's front matter, so this concerns a project's own convention files only.
- `design/system/conventions.md` and the user guide describe the new meaning.

## Alternatives considered

- Keep ADR-0059's meaning: a convention cannot be scoped to the sub-agents alone, and no roles meaning "story only" surprises a reader.
- A separate key for the story's agent: two keys for one question, who reads this file.
