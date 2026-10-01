---
id: ADR-0060
title: "A Claude Code PreToolUse hook, flai guard, refuses any sub-agent's call that would change a work item, a thread, or the repository's history"
status: accepted
date: 2026-10-01
supersedes: []
superseded_by: []
refines: [ADR-0059]
topics: [cli, conventions, template]
---

# ADR-0060 A Claude Code PreToolUse hook, flai guard, refuses any sub-agent's call that would change a work item, a thread, or the repository's history

## Context

ADR-0059 keeps sub-agents from acting for the story with their definitions' tool allowlists, which leave out every flai MCP tool that writes. The verifier needs `Bash` to run tests and lint, and with `Bash` it can run `flai move`, `flai thread new`, `git commit`, or `git stash` in the story's worktree, under the parent's `FLAI_AGENT`. The built-in sub-agents (`Explore`, `general-purpose`) have every flai MCP tool.

Probes on 2026-10-01 with Claude Code 2.1.286, headless, as `flai serve` starts it:

- A Bash pattern such as `Bash(flai move:*)` in a sub-agent's `disallowedTools` removes `Bash` from it whole, so a definition cannot keep the shell and fence off writes.
- A `hooks` block in a sub-agent's own definition did not fire.
- A `PreToolUse` hook in the project's `.claude/settings.json` fires for a sub-agent's calls too, and its input then carries `agent_type` (the definition's name) and `agent_id`; for the session's own agent it carries neither. A hook that exits 2 refuses the call and its standard error reaches the caller.

## Decision

A Claude Code `PreToolUse` hook, `flai guard`, refuses any sub-agent's call that would change a work item, a thread, or the repository's history; the story's agent's calls pass untouched.

- The template's `.claude/settings.json` runs `flai guard` before `Bash` and every `mcp__flai__` tool.
- When the input names no `agent_type`, it allows the call. Otherwise it refuses flai's MCP tools other than `board`, `doc_get`, `doc_search`, `item_get`, `prime`, `thread_get`, and `who_touches`; a flai command other than one that reads; and a git command other than one that reads. It finds the commands in a shell line by splitting at `;`, `&`, `|`, newlines, subshells, and command substitutions, and follows `sh -c`, `bash -c`, and `eval`.
- A refusal exits 2 with the reason and what to do instead (put it in the final message) on standard error. Input it cannot read passes: the guard fails open, so a missing or older `flai` costs the guarantee, not the session.
- It guards every sub-agent, the built-in ones included, not only the explorer and verifier.

## Consequences

- A sub-agent cannot move, create, or edit an item, write to a thread, read the inbox, commit, push, or switch branches in the story's worktree, whatever its definition says, as long as the project's settings run the hook and `flai` on `PATH` has `guard`.
- Files are not guarded: a sub-agent with `Bash` can still write a file through the shell. The verifier's instructions forbid it; a file it wrote is visible in `git status` before the story can go to review.
- A story that wants sub-agents to write, such as the task fan-out experiment after S-0175, changes the guard by a decision of its own.
- `flai guard` is Claude Code's hook format. Another harness with sub-agents and hooks gets an adapter of its own.
- This repository runs it through `scripts/flai.sh guard`, so the guard is the one in the tree; projects run the installed `flai`.

## Alternatives considered

- Bash patterns in the verifier's `disallowedTools`: they remove the shell whole.
- A `hooks` block in each definition: it did not fire headless, and it would guard only the sub-agents the template defines.
- A shell script under `scripts/` instead of a flai command: untested by flai's suite, and one more file each project keeps in step.
- Refusing in the flai CLI and MCP server: a sub-agent's calls arrive with the parent's name and environment, so flai cannot tell them apart there.
- No guard, the verifier's instructions only: the guarantee the story asks for would rest on a model following its prompt.
