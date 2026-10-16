---
id: T-1686
type: task
nature: feature
title: An ADR records where flai writes a story's OpenCode configuration so the worktree stays clean, and how it merges the project's own opencode.json
status: backlog
parent: S-0406
owner: alex
created: 2026-10-16T05:14:16Z
updated: 2026-10-16T05:14:16Z
transitions: []
stream: S-0406
tags: [cli]
touches: [design/adrs]
---
# T-1686 An ADR records where flai writes a story's OpenCode configuration so the worktree stays clean, and how it merges the project's own opencode.json

## Work

- Write the ADR with `flai adr new`, refining ADR-0130. Decide where the configuration with the MCP entry, the permission rules, and later the agents and the provider lives, so that a story's worktree has no uncommitted change from it: a file under `.flai-cache` that `OPENCODE_CONFIG` names, or a git-ignored path in the worktree.
- Say how the project's own `opencode.json` is merged in: which keys flai sets, and that it keeps every other key.
- Waits for nothing. The configuration writer in the next layer follows it.

## Done when

The ADR is accepted in `design/adrs`, listed in its index, and names the path, the merge, and the alternatives.

## Notes
