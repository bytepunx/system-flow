---
id: T-1705
type: task
nature: feature
title: The user, operator, and design docs say how the roles and the strategic agents run on OpenCode
status: backlog
parent: S-0409
owner: alex
created: 2026-10-16T05:17:44Z
updated: 2026-10-16T05:17:44Z
transitions: []
stream: S-0409
tags: [cli]
touches: [docs/users/flai.md, docs/operators/index.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/system/agent-adapters.md]
after: [T-1704]
---
# T-1705 The user, operator, and design docs say how the roles and the strategic agents run on OpenCode

## Work

- In `docs/users/flai.md` § Sub-agents, say that `.claude/agents/` stays the source and how a role's model reaches OpenCode; update the table in `docs/operators/index.md` § Harnesses.
- In `design/system/flai-cli.md` and `design/system/strategic-agents.md`, say how a strategic agent starts on `opencode`. Record the tool-name table in `agent-adapters.md` § OpenCode.
- Waits for T-1704: the docs describe the tested behaviour.

## Done when

The docs say what the story built, and the markdown lint passes.

## Notes
