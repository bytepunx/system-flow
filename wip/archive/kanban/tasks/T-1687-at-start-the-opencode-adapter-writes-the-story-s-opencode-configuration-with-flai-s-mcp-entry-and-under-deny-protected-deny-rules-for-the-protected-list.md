---
id: T-1687
type: task
nature: feature
title: At start the opencode adapter writes the story's OpenCode configuration with flai's MCP entry and, under deny-protected, deny rules for the protected list
status: done
parent: S-0406
owner: alex
created: 2026-10-16T05:14:27Z
updated: 2026-11-07T10:55:17Z
transitions:
  - to: ready
    at: 2026-11-07T10:43:48Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:43:51Z
    by: system-flow
  - to: done
    at: 2026-11-07T10:55:17Z
    by: agent-S-0406
stream: S-0406
tags: [cli]
touches: [flai/internal/harness/opencodeconfig.go, flai/internal/harness/opencodeconfig_test.go, flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/harness.go]
after: [T-1683, T-1685, T-1686]
usage:
  source: log
  seconds: 686
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 60
      output: 20355
      cache_read: 4034512
      cache_write: 116320
      cost: 1.8975
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
