---
id: T-0876
type: task
nature: remediation
title: A shell-only auto-approve host action lets permission_prompt allow .claude/ writes without asking
status: done
parent: S-0257
owner: alex
created: 2026-10-05T04:15:10Z
updated: 2026-10-05T04:27:24Z
transitions:
  - to: ready
    at: 2026-10-05T04:15:53Z
    by: agent-S-0257
  - to: in-progress
    at: 2026-10-05T04:24:49Z
    by: agent-S-0257
  - to: done
    at: 2026-10-05T04:27:24Z
    by: agent-S-0257
stream: S-0257
tags: []
touches: [flai/internal/hostapi, flai/cmd/mcp.go, flai/cmd/mcp_test.go]
after: [T-0874]
usage:
  source: log
  seconds: 155
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 42
      output: 284
      cache_read: 1347259
      cache_write: 71489
      cost: 0.5858
---
# T-0876 A shell-only auto-approve host action lets permission_prompt allow .claude/ writes without asking

## Work

1. Add the host action `auto-approve` to `hostapi.Actions`, and make it shell-only, like `auto-publish` (ADR-0067): no dashboard sees it or toggles it.
   - Its description says what it has flai do: `permission_prompt` allows, without asking you, a flai-serve agent's Edit or Write of a file in a `.claude/` folder in its own story's worktree.
   - Off, it asks you on a thread.
2. In `flai/cmd/mcp.go`, wire the T-0874 `Options` func so that the MCP server reads `cfg.ActionEnabled("auto-approve", root)` from the operator's configuration at every call.

## Done when

- `flai serve enable auto-approve` and `flai serve disable auto-approve` work, with `--all` too.
- The action is not in what a dashboard sees.
- Tests cover the wiring and the shell-only rule, wherever the existing auto-publish tests do.

## Notes
