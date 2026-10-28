---
id: ADR-0210
title: "flai guard holds a role a role file defines to its own name and session as ADR-0143 holds the planner, its own name being the hook's FLAI_AGENT or the name flai serve gives its run"
status: accepted
date: 2026-10-28
supersedes: []
superseded_by: []
refines: [ADR-0143]
topics: [cli]
---

# ADR-0210 flai guard holds a role a role file defines to its own name and session as ADR-0143 holds the planner, its own name being the hook's FLAI_AGENT or the name flai serve gives its run

## Context

[ADR-0143](0143-flai-guard-refuses-a-story-s-agent-the-planner-and-the-analyzer-a-call-that.md) refuses a served agent's own `Bash` call that names an author or a session other than its own with `--by`, `--placed-by`, `FLAI_AGENT`, or `FLAI_SESSION`, and names the agents it holds: a story's agent, the planner, and the analyzer (the curator joined them later). [ADR-0207](0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) §5 says every role a role file defines writes only as itself, as ADR-0143 holds the planner, and §2 gives such a role's run the `FLAI_AGENT` `<name>`, or `<name>-<item>` for a run on an item. S-0416 has the guard read a defined role's file; without this rule a call the file allows, such as `flai thread reply`, could name the operator with `--by` and answer a permission thread as the owner.

## Decision

`flai guard` holds a session whose `FLAI_ROLE` is a role a role file defines, any but the four built-in roles and the story's agent's and its sub-agents' verbs, to its own name and session exactly as ADR-0143 holds the planner, once the role file has allowed the call.

- **Its own name** is the hook's `FLAI_AGENT`, or the name `flai serve` gives the role's run: the role file's `name`, or `<name>-<item>` (the item from `FLAI_ITEM`, as a canonical ID) for a run on an item. When no role file was read, only the hook's `FLAI_AGENT` is its own.
- **Its own session** is the hook's `FLAI_SESSION`.
- Flags, environment, shells, and the refusal's advice are ADR-0143's, unchanged. A refusal of a defined role names it as `the role <name>`, as the role file's own refusals do.

## Consequences

- A defined role's `--by`, `--placed-by`, `FLAI_AGENT`, or `FLAI_SESSION` naming another is refused, so a thread entry, a move, or a message it writes through the command line carries its own name.
- A call the role file refuses is refused for the file's reason first; the own-name rule decides only what the file lets through.
- The shared paths' and `flai forge push`'s refusals of a defined role's session name it `the role <name>` rather than by its verb.
- The orchestrator stays under its own `--by` rule (S-0220) until a story holds it by its role file.
