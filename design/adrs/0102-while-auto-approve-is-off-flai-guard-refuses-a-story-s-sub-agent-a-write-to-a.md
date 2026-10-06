---
id: ADR-0102
title: "While auto-approve is off, flai guard refuses a story's sub-agent a write to a file in a .claude/ folder at once, and the story's agent makes that write itself"
status: proposed
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0060, ADR-0086]
---

# ADR-0102 While auto-approve is off, flai guard refuses a story's sub-agent a write to a file in a .claude/ folder at once, and the story's agent makes that write itself

## Context

[ADR-0086](0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md) has flai's `permission_prompt` decide a flai serve agent's Edit or Write of a file in a `.claude/` folder: it opens a thread on the story and holds the call until the owner answers, unless the operator has turned on `auto-approve`. It decides a sub-agent's call the same way, since Claude Code passes every session's prompts to the one handler.

I-0093 records what that cost. On S-0223, T-0957's task sub-agent wrote `template/root/.claude/agents/analyzer.md` and edited `template/root/.claude/settings.json`. `permission_prompt` opened TH-0192 and held the call. The owner was away, so the sub-agent, its layer, and the story's agent waiting on it all stopped, until the board watcher restarted the session about six minutes later; without the watcher it would have run to the thirty-minute timeout. The warning about `.claude/` writes reached the story's agent on a thread only after it had launched the layer. The way out, which agents then followed by hand (TH-0194, TH-0195), was for the story's agent to write the files into `.flai-cache/` and ask the owner on a thread to `cp` them in.

Under [ADR-0060](0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md) `flai guard` guarded no file, and the settings ran it before `Edit`, `Write`, and `NotebookEdit` only in a planner's, an orchestrator's, or an analyzer's session.

## Decision

While the project's `auto-approve` host action is off, `flai guard` refuses a story's sub-agent a write to a file in a `.claude/` folder at once, and the story's agent makes that write itself.

- The guard refuses a sub-agent's `Edit`, `MultiEdit`, `Write`, or `NotebookEdit` (the input carries an `agent_id`) of a file with a folder named `.claude` along its path. The refusal names the sub-agent and the file and tells it to put the file's whole new content in its final message. With `auto-approve` on, the write passes to `permission_prompt`, which allows it at once. A configuration or project the guard cannot read counts `auto-approve` as off.
- The settings, `.claude/settings.json` and the template's, run `flai guard` before `Edit`, `MultiEdit`, `Write`, and `NotebookEdit` in a story's session (`FLAI_STORY` set) as well as a strategic agent's. They still exit at once in the operator's own session.
- The story's agent's start prompt says, before it launches any sub-agent, that such a write is never a sub-agent's. It tells the agent to say so in the prompt of a sub-agent whose task changes such a file, and to make the write itself once the layer's sub-agents are back, through `permission_prompt`. When the owner may be away and nothing else is left, it writes each whole file into the worktree's ignored `.flai-cache/` folder instead, on a path with no folder named `.claude` along it, opens one thread on the story with the exact `cp` commands, and ends rather than wait.
- `permission_prompt` still decides the story's agent's own write as ADR-0086 says.

## Consequences

- A sub-agent that reaches for a `.claude/` file loses one turn rather than its layer's time, and the story's agent learns the content from its final message.
- The only call that waits on the owner for a `.claude/` write is the story's agent's own, made when nothing else of the story is running.
- The guard now runs before every file write in a story's session, about 18 ms a call.
- A project whose `.claude/settings.json` predates this keeps the hold until it takes the template's settings.
- Only Claude Code's sessions are covered, as with ADR-0060 and ADR-0086.

## Alternatives considered

- Have `permission_prompt` deny a sub-agent's write at once: its input, the tool, the tool's input, and the call's ID, does not say whether a sub-agent made the call, so it cannot tell a sub-agent's write from the story's agent's.
- A shorter `permission_prompt` timeout for every write: the story's agent's own write would also lose the operator's answer.
- The start prompt's warning alone: it would rest the guarantee on a model following its prompt, which ADR-0060 rejected for the same reason.
- Allow sub-agents these writes with no thread: that is what `auto-approve` already offers the operator, and it stays their choice.
