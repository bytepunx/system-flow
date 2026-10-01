---
id: T-0619
type: task
nature: improvement
title: Decide how sub-agents are defined, primed, attributed, and ask questions
status: done
parent: S-0175
owner: arobson
created: 2026-10-01T07:53:01Z
updated: 2026-10-01T07:54:42Z
transitions:
  - to: ready
    at: 2026-10-01T07:53:25Z
    by: agent-S-0175
  - to: in-progress
    at: 2026-10-01T07:53:25Z
    by: agent-S-0175
  - to: done
    at: 2026-10-01T07:54:42Z
    by: agent-S-0175
stream: S-0175
tags: []
touches: [design/adrs, design/system/agent-context.md]
usage:
  source: log
  seconds: 77
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 12
      output: 4770
      cache_read: 1137593
      cache_write: 14449
      cost: 0.4225
---
# T-0619 Decide how sub-agents are defined, primed, attributed, and ask questions

## Work

- Write an ADR: the claude-code adapter's agents hand noisy work and pre-review verification to two sub-agent roles, explore and verify, defined in the template with tools that only read; a sub-agent raises a question in its final message and the story's agent opens the thread; a sub-agent's calls are told apart in the agent's log by `parent_tool_use_id`, and the MCP server need not tell them apart because every tool they have is a read that uses no agent identity.
- Record the probe of 2026-10-01 (claude 2.1.286, headless): project `.claude/agents/` load, `tools` restricts MCP tools, an inline `mcpServers` on a sub-agent does not connect under `--strict-mcp-config`.
- Add a `## Sub-agents` section to `design/system/agent-context.md` with the choice and why.

## Done when

- The ADR is in `design/adrs` and indexed, and `agent-context.md` records the question choice among (a), (b), (c) and the attribution finding.

## Notes
