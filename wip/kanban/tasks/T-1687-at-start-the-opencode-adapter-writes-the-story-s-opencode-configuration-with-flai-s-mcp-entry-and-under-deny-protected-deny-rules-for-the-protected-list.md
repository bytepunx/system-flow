---
id: T-1687
type: task
nature: feature
title: At start the opencode adapter writes the story's OpenCode configuration with flai's MCP entry and, under deny-protected, deny rules for the protected list
status: backlog
parent: S-0406
owner: alex
created: 2026-10-16T05:14:27Z
updated: 2026-10-16T05:14:27Z
transitions: []
stream: S-0406
tags: [cli]
touches: [flai/internal/harness/opencodeconfig.go, flai/internal/harness/opencodeconfig_test.go, flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go]
after: [T-1683, T-1685, T-1686]
---
# T-1687 At start the opencode adapter writes the story's OpenCode configuration with flai's MCP entry and, under deny-protected, deny rules for the protected list

## Work

- In `flai/internal/harness/opencodeconfig.go`, build the configuration where T-1686's ADR puts it: the project's own `opencode.json`, if any, with `mcp.flai` set to a local entry, `command: ["flai", "mcp", "--agent", "<name>"]`, and every other key kept.
- When the host sets `--deny-protected` for `opencode`, add a `deny` rule to `permission` for each path `protected.For("opencode")` holds, for both edits and shell writes.
- Have `Start` in `opencode.go` write it and point OpenCode at it, so that the worktree has no uncommitted change from it.
- Test the merge with and without a project file, the deny rules, and that `git status` in a test worktree is clean after a start.
- Waits for T-1683 (the protected list), T-1685 (the adapter), and T-1686 (the ADR on where the file lives).

## Done when

A start writes the configuration where the ADR says, keeps the project's keys, adds the deny rules under `--deny-protected`, leaves the worktree clean, and the tests pass.

## Notes
