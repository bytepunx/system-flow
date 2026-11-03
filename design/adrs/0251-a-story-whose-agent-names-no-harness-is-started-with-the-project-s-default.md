---
id: ADR-0251
title: "A story whose agent names no harness is started with the project's default agent as it is at the start, under what the story's agent does say"
status: accepted
date: 2026-11-03
supersedes: []
superseded_by: []
refines: [ADR-0037]
---

# ADR-0251 A story whose agent names no harness is started with the project's default agent as it is at the start, under what the story's agent does say

## Context

[ADR-0037](0037-a-story-carries-its-agent-copied-from-the-project-s-default-when-it-is-made.md) copies the project's default agent, `agent` in `system-flow.yaml`, into a story when the story is made, and has `flai serve` start a story's agent from the story alone. A story made while the project had no default carries no `agent`, and setting the default later changes nothing for it. `flai serve` then refuses to start it: `S-nnnn names no harness, and no command is set on the host`.

S-0573 found this in a new project, `marketeer`. Its backlog of 32 stories was seeded before its default agent was set. The operator then set the default in `system-flow.yaml` and the harness on the dashboard's settings, and the dashboard's Start agent for S-0001 was still refused at 2026-11-03T17:17:34Z. The only way past it was to edit each story's agent by hand. The operator asked that configuring the project from the dashboard be enough.

The planner, the orchestrator, the analyzer, the curator, and the roles a project defines already resolve their agent at the start, as the default with their own block over it (`Manifest.Agent.With(...)`).

## Decision

**When a story's agent names no harness, `flai serve` starts it with the project's default agent as `system-flow.yaml` holds it at the start, with what the story's agent does say laid over it** (`Agent.With`, field by field, config key by key, and role by role). `workitem.Repo.AgentOf` resolves it, and every path that starts, restarts, resumes, or commits a story's agent, and every view of the agent it will run, reads the agent through it.

- A story whose agent names a harness keeps that agent whole. A later change to the default does not reach it, as ADR-0037 says: a story already worked does not change model under its agent.
- A story with no agent and a project with no default still names no harness. The host's command (`flai serve agent set -- <program>`) starts it, as before, and with neither the start is refused, as before.
- Nothing is written to the story. Its front matter still says only what was given for it, and `flai story new` still copies the default into a new story.
- A story whose agent says `harness: none` names a harness, so it is still run elsewhere ([ADR-0242](0242-a-story-whose-agent-names-harness-none-is-worked-by-an-agent-flai-does-not.md)).

## Consequences

- Setting the default, from the dashboard or with `flai agent set`, is enough for every open story made before it. No story needs editing.
- What such a story will run depends on the manifest at the start, not on its file alone. The dashboard's story page and `flai serve`'s journal name the harness and model of the run.
- A project that has both a default and a host command now starts an agent-less story with the default's harness, not the command. The default is the project's stated choice; the command stays the fallback when the project names none.

## Alternatives considered

- **Copy the default into every open story that has no agent when the default is set.** The front matter would stay the whole truth. But it writes and commits many items in the main checkout on one settings change, misses a story made later by hand or by an import, and leaves a project whose default was set by a hand edit of `system-flow.yaml`, as `marketeer`'s was, as broken as before.
- **Resolve the default under every story's agent at the start.** A change to the default would then reach stories that carry their own agent, which ADR-0037 rules out.
- **Keep the refusal and say how to fix each story.** It is the step the operator asked to remove.
