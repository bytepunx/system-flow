---
id: T-1002
type: task
nature: remediation
title: The users' guide says a write permission_prompt refuses reaches the agent with flai's reason
status: backlog
parent: S-0283
owner: alex
created: 2026-10-06T06:26:55Z
updated: 2026-10-06T06:26:55Z
transitions: []
stream: S-0283
tags: [docs]
touches: [docs/users/flai.md]
after: [T-1001]
---
# T-1002 The users' guide says a write permission_prompt refuses reaches the agent with flai's reason

## Work

It waits for T-1001, because it describes what T-1001's fix delivers, in T-1001's words for a refusal.

In `docs/users/flai.md`, under `Writes under .claude/`, say that a write `permission_prompt` refuses at once, such as one to `/tmp`, to the main checkout, or for a story not in progress, reaches the agent with flai's reason, so the agent writes in its story's worktree instead. Keep the `permission_prompt` row of the MCP tools table true; change it only if T-1001 changed what the tool does.

## Done when

- the section names what a refused write tells the agent, matching T-1001's messages
- the markdown lint passes on `docs/users/flai.md`

## Notes
