---
id: T-3067
type: task
nature: improvement
title: The operators guide, the container runbook, and the user guide say how to give a story's agent outside folders, and that Monitor is allowed
status: backlog
parent: S-0701
owner: alex
created: 2026-11-10T10:35:51Z
updated: 2026-11-10T10:35:51Z
transitions: []
stream: S-0701
tags: [docs]
touches: [docs/operators/index.md, docs/operators/runbooks/agent-container.md, docs/users/flai.md, docs/operators/settings.md]
after: [T-3063, T-3064]
---
# T-3067 The operators guide, the container runbook, and the user guide say how to give a story's agent outside folders, and that Monitor is allowed

## Work

In `docs/operators/index.md`, update the claude-code harness paragraph and its example arguments: the new default `--allowedTools Bash,Monitor,mcp__flai`, and `--outside`, with what it grants and what it refuses. In `docs/operators/runbooks/agent-container.md`, add the outside folders to the container's mounts. In `docs/users/flai.md` § Sub-agents and the protected-path writes section, say that Read, Grep, Glob, Write, and Monitor reach the outside folders the operator granted and nothing else outside the project. In `docs/operators/settings.md`, update the default in the row for `agent.harnesses.<name>.args`.

This task waits for T-3063 and T-3064, so that it documents the mounts and the prompt as built.

## Done when

- Each of the four documents says what the built behaviour does, and links the ADR from T-3048.
- The markdown lint and `flai check --strict` are clean.

## Notes
