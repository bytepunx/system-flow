---
id: S-0175
type: story
nature: improvement
title: Agents the claude-code adapter starts hand noisy work and verification to sub-agents that cannot act for the story
status: in-progress
owner: alex
created: 2026-10-01T07:17:20Z
updated: 2026-10-01T08:13:01Z
transitions:
  - to: ready
    at: 2026-10-01T07:38:46Z
    by: alex
  - to: in-progress
    at: 2026-10-01T07:48:30Z
    by: agent-S-0175
tags: [flai, template]
topics: [conventions]
touches: [flai/internal/harness, flai/internal/context, flai/cmd/prime.go, flai/cmd/prime_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, flai/cmd/root.go, flai/internal/guard, flai/internal/mcpserver, flai/internal/conventions, template, design/conventions, design/system/agent-context.md, design/system/flai-cli.md, design/system/conventions.md, design/adrs, ".claude", docs/users, docs/operators/index.md, docs/operators/settings.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1260
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 246
      output: 1222
      cache_read: 21410959
      cache_write: 317798
      cost: 9.7457
---
# S-0175 Agents the claude-code adapter starts hand noisy work and verification to sub-agents that cannot act for the story

## Goal

An agent flai serve starts with the `claude-code` harness works its whole story in one context. Every search, test run, lint log, and file it reads stays there and is read again on every later turn: S-0118's run cost $5.18 with 10.5M cache-read tokens for one Explore sub-agent it spawned on its own. The `Task` tool is already in these headless sessions (the init event lists it, though `--allowedTools` is `Bash,mcp__flai`), but nothing in `harness.Prompt`, the conventions, or the template says when to use it, and a sub-agent shares the parent's flai MCP connection, so it acts as `agent-S-nnnn`: it can move the story, edit items, and consume the parent's inbox.

This story makes sub-agents a deliberate, safe part of the claude-code adapter:

1. **Offload noisy work.** The story's agent hands code search, test and lint runs, and long log reading to sub-agents that return a summary, keeping its own context to decisions and edits.
2. **Verify before review.** Before moving the story to review, the agent has a fresh sub-agent check the diff against the acceptance criteria and the conventions, and acts on what it finds.
3. **Sub-agents cannot act for the story.** The template ships agent definitions (an explorer and a verifier) whose tools leave out what changes a work item or consumes the inbox, and a smaller prime for those roles.

Sub-agents start from a fresh context (only the prompt the parent writes), see none of the parent's primed pack, and return only their final message. Item 3 needs care because tasks are where questions for the designer most often arise.

## Acceptance criteria
- [ ] `harness.Prompt` for `claude-code`, and a sub-agent section in the conventions (this repository and `template/`), say when to delegate (search, test, lint, long logs, pre-review verification), what to put in a sub-agent's prompt (the worktree, the story and task IDs, the question), and that a sub-agent returns a summary, not raw output
- [ ] The template ships `.claude/agents/` definitions for an explorer (read-only) and a verifier (reads, runs tests and lint, no edits), whose tools exclude `item_move`, `item_edit`, `item_new`, `inbox`, `wait_for_work`, `wait_for_events`, and `thread_resolve`; the claude-code adapter's sessions load them
- [ ] `prime` takes a role (for example `flai prime --story S-nnnn --role explore|verify`, and the MCP tool's equivalent) that returns a pack sized for a sub-agent: the conventions its role needs, the story's goal and criteria, and briefs, within a smaller budget
- [ ] The story's agent stays the only one that talks to the designer for the story. How a sub-agent raises a question is decided and built, from: (a) it returns the question in its final message and the parent opens the thread; (b) it may call `thread_open` on its task, attributed to the parent's agent, and the parent's `wait_for_events` sees it; (c) the harness's own messaging between agents (Claude Code's `SendMessage` to a running background agent, or agent teams), if it works headless. Record the choice and why in `design/system/agent-context.md`; an ADR if it changes who may open threads
- [ ] A sub-agent's flai calls are distinguishable from its parent's in `flai serve`'s logs or the MCP server's attribution, or the design says why they need not be
- [ ] Measured on at least two stories with the cost tracking from ADR-0051 (sub-agent usage included): cost, cache reads, and turns against comparable runs without delegation, recorded in `design/system/agent-context.md`
- [ ] The design (`design/system/flai-cli.md`, `agent-context.md`) and the user guide describe what agents delegate and what sub-agents may do

## Tasks
- T-0619 Decide how sub-agents are defined, primed, attributed, and ask questions
- T-0620 flai prime --role and the MCP prime tool's role return a sub-agent's pack
- T-0621 The template ships explorer and verifier sub-agents and a delegation convention
- T-0622 harness.Prompt for claude-code says when and how to delegate
- T-0623 The design and user guide describe what agents delegate and what sub-agents may do
- T-0624 Measure delegation on two stories against runs without it
- T-0627 flai guard refuses a sub-agent's writes as a Claude Code hook

## Notes

- Forked sub-agents inherit the parent's conversation and prompt cache, so they start primed; check whether the `claude` that `flai serve` runs headless exposes them, and whether they inherit the parent's tool restrictions. They are the basis for the task fan-out experiment that follows this story, not for this one.
- ADR-0049 rejected a scout sub-agent that reads the whole pack on the parent's behalf. This story does the reverse: sub-agents get less than the parent, not more.
- The `command` harness is out of scope: flai cannot know what its program does with sub-agents.
- Decided in the story (ADR-0059): questions come back in a sub-agent's final message and the story's agent opens the thread (option (a)); a sub-agent's calls are told apart in the agent's log by `parent_tool_use_id`, and the MCP server is unchanged.
- Added in the story (ADR-0060, T-0627): the verifier needs `Bash`, and a Bash pattern in a definition's `disallowedTools` removes Bash whole, so `flai guard`, a project `PreToolUse` hook, refuses any sub-agent's flai and git writes by the `agent_type` Claude Code passes it.
- Measurement (TH-0042): this run and S-0118 are measured here; S-0188 measures two stories run with the released prompt.
