---
id: T-1702
type: task
nature: feature
title: The opencode adapter puts the translated agents in the story's configuration, starts a strategic agent with --agent, and states role models and strategic definitions where Guard is true
status: backlog
parent: S-0409
owner: alex
created: 2026-10-16T05:17:30Z
updated: 2026-10-16T05:23:20Z
transitions: []
stream: S-0409
tags: [cli]
touches: [flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/opencodeconfig.go, flai/internal/harness/opencodeconfig_test.go, flai/internal/harness/harness_test.go]
after: [T-1701]
---
# T-1702 The opencode adapter puts the translated agents in the story's configuration, starts a strategic agent with --agent, and states role models and strategic definitions where Guard is true

## Work

- In `opencodeconfig.go`, add the translated agents where S-0406's ADR puts the story's configuration, keeping the worktree clean.
- In `opencode.go`, pass `--agent <name>` when the start is a strategic role's, as the `claude-code` adapter does for `FLAI_ROLE`.
- Where `Guard` is true, set the role-model and strategic-definition capabilities; where it is false, leave them off. Update `TestEachAdapterStatesItsCapabilities`.
- Test a story's start with roles, a planner's start, and the clean worktree.
- Waits for T-1701: it writes what that task translates.

## Done when

A start carries the agents and the strategic `--agent`, the capabilities follow `Guard`, and the harness tests pass.

## Notes
