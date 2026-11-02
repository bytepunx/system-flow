---
id: TH-0793
title: Allow Edit .claude/settings.json?
anchor:
  path: wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md
  item: S-0518
status: open
participants: [agent-S-0518, orchestrator]
created: 2026-11-02T00:01:48Z
updated: 2026-11-02T00:02:57Z
---

# TH-0793 Allow Edit .claude/settings.json?

On wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md.

## Entries

### 2026-11-02T00:01:48Z agent-S-0518
agent-S-0518 asks to Edit `.claude/settings.json` in S-0518's worktree. Claude Code protects the path and does not write it without a person's approval.

Reply `allow`, as alex, the story's owner, to let it write. Anything else refuses it, and your words go back to the agent as the reason; a reply by anyone else is not an answer.

It would replace this text:

```text
            "command": "out=$(\"$CLAUDE_PROJECT_DIR\"/scripts/flai.sh guard 2>&1); [ $? -eq 2 ] || exit 0; echo \"$out\" >&2; exit 2"
```

with this text:

```text
            "command": "out=$(\"$CLAUDE_PROJECT_DIR\"/scripts/flai.sh guard 2>&1); s=$?; [ \"$s\" -eq 2 ] && { echo \"$out\" >&2; exit 2; }; case \"$s\" in 126|127) printf 'flai guard did not run (exit %s): %s\\nno guarded call passes until flai is installed and on PATH\\n' \"$s\" \"$out\" >&2; exit 2;; esac; exit 0"
```

### 2026-11-02T00:02:57Z orchestrator (recommendation)
Recommendation: reply `allow`, as alex. This is the same change as TH-0792, made to this repository's own `.claude/settings.json` instead of the template's.

- It changes only the first PreToolUse hook (`Bash|mcp__flai__.*`). That hook now refuses a call when `scripts/flai.sh guard` could not run (exit 126 or 127), instead of letting the call through. A refusal from the guard (exit 2) is still blocked, and every other exit still passes. This is what S-0518 is for (I-0216).
- Under ADR-0106 the write stays on `story/S-0518`. Only you can accept the story, because it changes a protected path, so you will see this diff again at acceptance.

Check these at acceptance:

- Here the hook runs `scripts/flai.sh`, which builds `bin/flai` from the tree. Exit 127 can also mean `go` or the script is missing. The message "until flai is installed and on PATH" is then slightly wrong for this repository.
- The second hook, the one that runs on `FLAI_ROLE` or `FLAI_STORY`, still lets a call through on 126 or 127. Decide whether S-0518 should close that hook as well.
- Exit 1 (the guard starts, then crashes) still lets the call through, as on TH-0792.

Only your reply counts as the answer here, so this is marked as yours alone.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md § Decision
