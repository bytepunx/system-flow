---
id: TH-0792
title: Allow Edit template/root/.claude/settings.json?
anchor:
  path: wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md
  item: S-0518
status: resolved
participants: [agent-S-0518, orchestrator]
created: 2026-11-01T23:57:41Z
updated: 2026-11-02T00:20:13Z
---

# TH-0792 Allow Edit template/root/.claude/settings.json?

On wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md.

## Entries

### 2026-11-01T23:57:41Z agent-S-0518
agent-S-0518 asks to Edit `template/root/.claude/settings.json` in S-0518's worktree. Claude Code protects the path and does not write it without a person's approval.

Reply `allow`, as alex, the story's owner, to let it write. Anything else refuses it, and your words go back to the agent as the reason; a reply by anyone else is not an answer.

It would replace this text:

```text
            "command": "out=$(flai guard 2>&1); [ $? -eq 2 ] || exit 0; echo \"$out\" >&2; exit 2"
```

with this text:

```text
            "command": "out=$(flai guard 2>&1); s=$?; [ \"$s\" -eq 2 ] && { echo \"$out\" >&2; exit 2; }; case \"$s\" in 126|127) printf 'flai guard did not run (exit %s): %s\\nno guarded call passes until flai is installed and on PATH\\n' \"$s\" \"$out\" >&2; exit 2;; esac; exit 0"
```

### 2026-11-02T00:00:44Z orchestrator (recommendation)
Recommendation: reply `allow`, as alex.

- The edit does what S-0518 is for (I-0216). The hook now refuses a call when `flai guard` could not start (exit 126 or 127) instead of letting it through. A refusal from the guard (exit 2) is blocked as before. Every other exit still passes.
- The cost: in a project made from the template, every guarded call is refused until flai is installed and on PATH. That is the fail-closed behaviour the issue asks for.
- Under ADR-0106 the write stays on `story/S-0518` and changes nothing yet. Only you can accept the story, because it changes a protected path, so you will review this diff again then.
- Only your reply counts as the answer here. Mine does not, so this is marked as yours alone.

One thing to check at acceptance: exit 1, the guard starting and then crashing, still lets the call through. Decide then whether S-0518 should refuse that too.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md § Decision

### 2026-11-02T00:20:13Z agent-S-0518
Resolved: Superseded by TH-0795: that write changes both PreToolUse hooks of template/root/.claude/settings.json in one go, so this narrower one is withdrawn. Answer TH-0795 instead.
