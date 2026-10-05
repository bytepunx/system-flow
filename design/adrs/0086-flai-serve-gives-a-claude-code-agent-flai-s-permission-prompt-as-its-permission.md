---
id: ADR-0086
title: "flai serve gives a claude-code agent flai's permission_prompt as its permission handler, which asks the story's owner before a write under .claude/, unless the operator turns on auto-approve"
status: accepted
date: 2026-10-05
supersedes: []
superseded_by: []
refines: [ADR-0038]
topics: [cli]
---

# ADR-0086 flai serve gives a claude-code agent flai's permission_prompt as its permission handler, which asks the story's owner before a write under .claude/, unless the operator turns on auto-approve

## Context

flai serve starts a claude-code story agent headless (`claude -p`), with the operator's host arguments, by default `--permission-mode acceptEdits --allowedTools Bash,mcp__flai` ([ADR-0038](0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md)). Claude Code treats a write under any `.claude/` folder, the project's settings, hooks, and agent definitions and the template's copies of them, as sensitive. Only a person, or the permission handler named by `--permission-prompt-tool`, may approve one. Headless there is no person, and flai named no handler, so every such write was refused, to the story's agent and its sub-agents alike, whatever the operator granted on a thread.

I-0069 counts six instances, from S-0208, S-0209, S-0210, and S-0266. Each ended with the operator pasting a whole file the agent sent on a thread, and once with a partial paste that would have dropped the sub-agent guard. A story that changes the guard hook's matcher or an agent definition could not ship it itself.

On TH-0128 the operator chose a handler that asks them, "but extend it so that the operator has the choice to enable auto-approved writes for agents".

## Decision

The claude-code adapter passes `--permission-prompt-tool mcp__flai__permission_prompt` to every session flai serve starts, ahead of the operator's host arguments. flai's MCP tool `permission_prompt` then decides what Claude Code would have asked a person.

1. **It approves one thing: an Edit, Write, MultiEdit, or NotebookEdit of a file in a `.claude/` folder inside an in-progress story's worktree.** The path is absolute, lies below `.flai-cache/worktrees/S-nnnn/` with no symbolic link taking it out, and has a `.claude` folder below the worktree. Anything else it is asked, any other tool, a path in the main checkout or outside every worktree, a story not in progress, is denied at once, as it was before there was a handler.
2. **It asks the story's owner on a thread.** It opens a thread on the story, "Allow <tool> <path>?", that shows the whole content written or the old and new text of each edit. It then waits. A reply by the story's owner whose first word is allow, yes, approve, approved, or ok lets the write through; any other reply by the owner refuses it with their words as the reason. Entries by anyone else, other agents included, are not answers. The tool resolves the thread with what was decided, and refuses when the session ends unanswered.
3. **`auto-approve` is a host action, off by default, that allows such a write at once, with no thread.** `flai serve enable auto-approve` turns it on for the project, `--all-projects` for every project, and `flai serve disable auto-approve` turns it off. flai mcp reads it at every request, so the change applies to agents already running. Like `auto-publish` it is the operator's shell tool ([ADR-0067](0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md)): no dashboard sees it or changes it, so the dashboard token cannot give agents the right to rewrite their own hooks and permissions. Each auto-approved write is logged.

This refines ADR-0038: what the agent may do is still the operator's, and the handler adds no permission the operator did not give, through their answer or through `auto-approve`.

## Consequences

- A story that changes a `.claude/` file, such as the guard hook's matcher or an agent definition, ships it from its own branch. The operator reads the change on a thread and answers `allow` rather than pasting a file.
- An agent waiting on the handler is blocked in a tool call, as one waiting on a thread is, until the owner answers or the session ends.
- Prompts other than these writes are still refused, so a host argument that narrows the agent's permissions keeps its effect.
- With `auto-approve` on, an agent can change its own settings and hooks unattended. That is the operator's choice for a project, made in a shell.
- The handler is Claude Code's. Other harnesses get none, and the command harness runs what the operator configured.
- Bash remains outside the handler: a shell write under `.claude/` is not stopped by it, just as before, and the conventions keep agents from writing these files through the shell.
