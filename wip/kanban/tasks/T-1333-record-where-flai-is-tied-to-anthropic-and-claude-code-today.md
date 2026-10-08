---
id: T-1333
type: task
nature: research
title: Record where flai is tied to Anthropic and Claude Code today
status: done
parent: S-0339
owner: alex
created: 2026-10-08T07:20:27Z
updated: 2026-10-08T07:48:34Z
transitions:
  - to: ready
    at: 2026-10-08T07:39:55Z
    by: agent-S-0339
  - to: in-progress
    at: 2026-10-08T07:39:56Z
    by: agent-S-0339
  - to: done
    at: 2026-10-08T07:48:34Z
    by: agent-S-0339
stream: S-0339
tags: [research, agents]
touches: [design/system/agent-adapters.md, design/system/README.md]
usage:
  source: log
  seconds: 518
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 492
      output: 31502
      cache_read: 1772090
      cache_write: 193363
      cost: 4.9027
    - model: claude-haiku-4-5-20251001
      input: 853006
      output: 55817
      cache_read: 3654023
      cache_write: 260277
      cost: 2.0628
---
# T-1333 Record where flai is tied to Anthropic and Claude Code today

## Work

Start `design/system/agent-adapters.md` as the finding of S-0339, in the shape of `release-signing.md` and `agent-coordination.md`. Write a `## Today` section that states, from the code and the ADRs, each place flai assumes Claude Code or Anthropic. Group it by concern:

- **Starting and resuming an agent.** The `Adapter` interface and registry in `flai/internal/harness/harness.go`, and the `claude-code` and `command` adapters in `flai/internal/harness/adapters.go`. Cover `claude -p`, session resume, and how harness, model, config, and roles reach the argument list (ADR-0037, ADR-0038, ADR-0065).
- **Permissions.** `permission_prompt` as `--permission-prompt-tool` (`flai/internal/mcpserver/permission.go`, ADR-0086, ADR-0097, ADR-0124), and the check against each new Claude Code (`flai/internal/serve/claudecheck.go`, ADR-0106).
- **The guard.** `flai guard` as a `PreToolUse`, `SubagentStart`, and `SubagentStop` hook (`flai/internal/guard/`, `template/root/.claude/settings.json`, ADR-0060, ADR-0102).
- **Sub-agents and strategic agents.** The definitions in `template/root/.claude/agents/` and the roles they serve (ADR-0059, ADR-0075, ADR-0082).
- **Usage and cost.** What is read from stream-json `result` events and their `modelUsage` and `costUSD` (`flai/internal/usage/log.go`), turn classes (`flai/internal/usage/turns.go`), empty wakes, and the task apportionment (ADR-0051, ADR-0071, ADR-0105, ADR-0116).
- **Model names and files the harness reads.** Model names in the manifest and stories, `CLAUDE.md`, and `.mcp.json`.

For each concern, say what is neutral already, such as MCP, and what only Claude Code provides. Add the document's row to `design/system/README.md`. This task waits for no other task.

## Done when

- `design/system/agent-adapters.md` exists with front matter and a `## Today` section, one subsection per concern above.
- Every statement in `## Today` names the file or ADR it comes from.
- `design/system/README.md` lists the document.

## Notes
