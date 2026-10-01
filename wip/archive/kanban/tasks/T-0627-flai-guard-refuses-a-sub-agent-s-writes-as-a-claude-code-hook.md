---
id: T-0627
type: task
nature: improvement
title: flai guard refuses a sub-agent's writes as a Claude Code hook
status: done
parent: S-0175
owner: arobson
created: 2026-10-01T08:03:40Z
updated: 2026-10-01T08:06:56Z
transitions:
  - to: ready
    at: 2026-10-01T08:03:44Z
    by: agent-S-0175
  - to: in-progress
    at: 2026-10-01T08:03:44Z
    by: agent-S-0175
  - to: done
    at: 2026-10-01T08:06:56Z
    by: agent-S-0175
stream: S-0175
tags: []
touches: [flai/internal/guard, flai/cmd/guard.go, template/, ".claude/"]
usage:
  source: log
  seconds: 192
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 30
      output: 11933
      cache_read: 2845740
      cache_write: 36145
      cost: 1.0569
---
# T-0627 flai guard refuses a sub-agent's writes as a Claude Code hook

## Work

- `flai guard` reads a Claude Code PreToolUse hook's input and, when `agent_type` says a sub-agent makes the call, refuses flai's MCP writes and inbox, flai commands that write, and git commands that change the worktree or history, with exit 2 and the reason on standard error.
- The template's `.claude/settings.json` runs it before `Bash` and flai's MCP tools; this repository's runs it through `scripts/flai.sh`.
- ADR-0060 records it; the probe of 2026-10-01 grounds it: a Bash pattern in a sub-agent's `disallowedTools` removes Bash whole, a sub-agent's own `hooks` did not fire headless, and a project hook sees `agent_type` for a sub-agent's calls and none for the session's agent.

## Done when

- `guard` and `cmd` tests cover MCP writes, flai and git writes through the shell, reads, the story's agent, and unreadable input; a headless session shows a sub-agent refused.

## Notes
