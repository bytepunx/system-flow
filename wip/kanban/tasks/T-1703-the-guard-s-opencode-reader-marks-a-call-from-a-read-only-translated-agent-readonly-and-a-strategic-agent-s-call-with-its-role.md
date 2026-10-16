---
id: T-1703
type: task
nature: feature
title: The guard's opencode reader marks a call from a read-only translated agent ReadOnly and a strategic agent's call with its role
status: backlog
parent: S-0409
owner: alex
created: 2026-10-16T05:17:34Z
updated: 2026-10-16T05:17:34Z
transitions: []
stream: S-0409
tags: [cli]
touches: [flai/internal/guard/opencode.go, flai/internal/guard/opencode_test.go]
after: [T-1701]
---
# T-1703 The guard's opencode reader marks a call from a read-only translated agent ReadOnly and a strategic agent's call with its role

## Work

- In `flai/internal/guard/opencode.go`, set `ReadOnly` for a sub-agent call whose agent is one T-1701 translated with no edit, and `Role` from `FLAI_ROLE` for a strategic agent's own calls, as the Claude Code reader does.
- Test the explorer's shell write refused, the verifier's read let through, and a planner's call outside planning refused.
- Waits for T-1701: it reads the names and tools that task gives the agents.

## Done when

The reader's marks match the Claude Code reader's for the same callers, and the guard's tests pass.

## Notes
