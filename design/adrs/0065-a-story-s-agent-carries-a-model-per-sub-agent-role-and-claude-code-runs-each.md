---
id: ADR-0065
title: "A story's agent carries a model per sub-agent role, and claude-code runs each role's sub-agent on it"
status: accepted
date: 2026-10-01
supersedes: []
superseded_by: []
refines: [ADR-0037, ADR-0059]
topics: [cli, conventions]
---

# ADR-0065 A story's agent carries a model per sub-agent role, and claude-code runs each role's sub-agent on it

## Context

S-0188 measured that the explorer and the verifier ran on the story's agent's model, and that what they cost was a tenth of each run. ADR-0059 left a cheaper model to each project's own definitions, in `.claude/agents/`. The designer asked for agent and model configuration by role: the story's own, `explore`, and `verify`, with room for more (S-0189). A story's `agent` (ADR-0037) named only the story's harness, model, and config, so a story could not run its verifier on another model without editing a file every story shares. The designer chose this shape in TH-0057.

## Decision

A story's and the project's `agent` carry an optional `roles` map from a sub-agent role to its own harness, model, and config, and `flai serve` starts `claude-code` with each role's model over the project's definition of that sub-agent.

1. **Shape.** `agent.roles.<role>` has `harness`, `model`, and `config`, checked as the agent's own are. A role is a lower-case name; `explore` and `verify` are the ones flai's harnesses know. A role sets something or is left out. The top-level fields stay the story's own agent, so an agent with no roles reads and writes as before.
2. **Default and merge.** A story made while the project's default has roles gets them, merged role by role and key by key under what it is given, as ADR-0037 merges the rest. `Same` compares roles, so a changed role restarts a ready story's agent as a changed model does.
3. **Setting them.** `flai agent set`, `flai story new`, and `flai edit` take `--role-harness role=h`, `--role-model role=m`, `--role-config role.key=value` (no value removes the key), and `--unset-role role`. The MCP `item_new` and `item_edit` take `roles` in `agent`. The dashboard shows roles and keeps them on a save; it does not edit them.
4. **claude-code.** A role with a model is passed as `--agents`: the project's `.claude/agents/<definition>.md` (`explore` is `explorer`, `verify` is `verifier`), its front matter but `name`, and its body as `prompt`, with the role's model. A session's `--agents` outranks the project's files, so the definition stays the source of everything but the model. A role with no definition, on another harness, or with config refuses the start, because a sub-agent runs in the story's own session and takes only a model.
5. **command.** The operator's command gets the roles as `FLAI_AGENT_ROLES`, a JSON object.

## Consequences

- A project can run its sub-agents cheaper, or a story can run its verifier on a stronger model, without editing the definitions every story shares. The template's definitions name `haiku` and `sonnet` (S-0189), so most stories need no roles.
- A flai older than this one does not know `agent.roles`. One from before S-0181 refuses the story; one from S-0181 on reads past it and would drop it when it rewrote the story, and its `flai agent set` would drop the default's roles. `front-matter-fields.txt` therefore lists the agent block's keys and a role's, so the release that carries this raises `flai.minimum` (S-0181) and an older flai stops before it reads the project.
- The adapter now reads two files in the project to start an agent. A malformed definition stops a story that names its role, not one that does not.
- Corrections stay the story's agent's whatever a role's model (ADR-0059, ADR-0060).
