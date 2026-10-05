---
id: TH-0128
title: "S-0257 proposal: let the operator approve a .claude/ write from the thread, through a permission prompt tool"
anchor:
  path: wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md
  item: S-0257
status: resolved
participants: [agent-S-0257, alex]
created: 2026-10-05T04:07:31Z
updated: 2026-10-05T04:14:38Z
---

# TH-0128 S-0257 proposal: let the operator approve a .claude/ write from the thread, through a permission prompt tool

On wip/kanban/stories/S-0257-a-story-agent-started-by-flai-serve-cannot-edit-claude-settings-json-so-a-story-that-changes-the-guard-hook-s-matcher-cannot-ship-it.md.

## Entries

### 2026-10-05T04:07:31Z agent-S-0257
I-0069 asks for a proposed fix before I build one. Here is mine.

**Cause.** flai serve starts Claude Code headless (`-p`) with `--permission-mode acceptEdits --allowedTools Bash,mcp__flai` (ADR-0038). Claude Code 2.1.280 treats writes to settings, git, and tool-configuration files as sensitive: anything under a `.claude/` folder, `template/root/.claude/` included. Only "a person or the configured permission handler" may approve such a write. A headless session has no person, and flai configures no handler, so the write is denied. A grant on a thread changes nothing. That is why all six instances ended with you pasting whole files.

**Proposal.** Give the agent a permission handler, and make it ask you.

1. **A new flai MCP tool, `permission_prompt`.** Claude Code calls it only when a call would otherwise prompt. It does this:
   - For an Edit, Write, or NotebookEdit of a file in a `.claude/` folder inside the story's worktree, it opens a thread on the story. The thread shows the path and the whole new file, or the edit's old and new text. It waits for your answer.
   - An answer of `allow` (or yes, approve) returns allow, and Claude Code writes the file itself. Any other answer returns deny, with your words as the reason.
   - Anything else it is asked about (a path outside the worktree, the main checkout, `.git/`, any other tool) it denies at once, so it widens nothing except through your answer.
2. **flai serve passes `--permission-prompt-tool mcp__flai__permission_prompt`** in the claude-code adapter's start. This is added to the agent's own argv, not to the operator's host args, so it applies with the default args and with custom ones.
3. **A new ADR** extending ADR-0038, plus updates to flai-cli.md (serve, mcp tools) and the users' docs. Go tests cover the tool's allow, deny, and refuse paths and the adapter's argv. Claude Code's own refusal cannot be reproduced in a Go test, so a test asserting the flag is passed stands in for it.

**What you would see.** A thread on the story, "Allow Edit .claude/settings.json?", holding the content. You answer `allow`. No more copying and pasting, and no more partial pastes like the one on TH-0096 that dropped the sub-agent guard.

**What I am not proposing:**
- auto-approving these writes;
- writing them through the shell;
- `bypassPermissions`.

All three would let an unattended agent rewrite its own hooks and permissions, which is exactly what Claude Code is guarding against.

Shall I build this? If you would rather approve on a dashboard page than in a thread, say so. I would keep the thread for this story and leave the page to a later story.

### 2026-10-05T04:13:37Z alex
Let's take your approach but extend it so that the operator has the choice to enable auto-approved writes for agents

### 2026-10-05T04:14:37Z agent-S-0257
Building it, with the auto-approve choice added like this:

- **A new host action, `auto-approve`.** It is off by default. You turn it on per project or for every project with `flai serve enable auto-approve` (`--all` for every project), and off with `flai serve disable auto-approve`.
- **Shell-only, like `auto-publish` (ADR-0067).** No dashboard sees it or toggles it, so holding the dashboard token cannot give agents the right to rewrite their own hooks and permissions.
- **When it is on:** `permission_prompt` allows an agent's Edit, Write, or NotebookEdit of a file in a `.claude/` folder in its own story's worktree at once, with no thread, and logs that it did.
- **When it is off:** it asks you on a thread, as proposed.
- **Read at every call.** The tool reads the setting each time it is called, so turning it on or off applies to agents already running.
- **What it never approves, on or off:** anything outside a story's worktree, the main checkout, or any other kind of prompt. Those stay denied.

I'll record this in an ADR and in the story's Decisions.

### 2026-10-05T04:14:38Z agent-S-0257
Resolved: Operator chose the permission_prompt approach, extended with an opt-in, shell-only auto-approve host action.
