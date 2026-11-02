---
id: TH-0795
title: Allow Write template/root/.claude/settings.json?
anchor:
  path: wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md
  item: S-0518
status: resolved
participants: [agent-S-0518, orchestrator, alex]
created: 2026-11-02T00:11:57Z
updated: 2026-11-02T00:40:08Z
---

# TH-0795 Allow Write template/root/.claude/settings.json?

On wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md.

## Entries

### 2026-11-02T00:11:57Z agent-S-0518
agent-S-0518 asks to Write `template/root/.claude/settings.json` in S-0518's worktree. Claude Code protects the path and does not write it without a person's approval.

Reply `allow`, as alex, the story's owner, to let it write. Anything else refuses it, and your words go back to the agent as the reason; a reply by anyone else is not an answer.

The whole content it would write:

```text
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash|mcp__flai__.*",
        "hooks": [
          {
            "type": "command",
            "command": "out=$(flai guard 2>&1); s=$?; [ \"$s\" -eq 2 ] && { echo \"$out\" >&2; exit 2; }; case \"$s\" in 126|127) printf 'flai guard did not run (exit %s): %s\\nno guarded call passes until flai is installed and on PATH\\n' \"$s\" \"$out\" >&2; exit 2;; esac; exit 0"
          }
        ]
      },
      {
        "matcher": "Edit|MultiEdit|Write|NotebookEdit",
        "hooks": [
          {
            "type": "command",
            "command": "[ -n \"$FLAI_ROLE\" ] || [ -n \"$FLAI_STORY\" ] || exit 0; out=$(flai guard 2>&1); s=$?; [ \"$s\" -eq 2 ] && { echo \"$out\" >&2; exit 2; }; case \"$s\" in 126|127) printf 'flai guard did not run (exit %s): %s\\nno guarded call passes until flai is installed and on PATH\\n' \"$s\" \"$out\" >&2; exit 2;; esac; exit 0"
          }
        ]
      }
    ],
    "SubagentStart": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "[ -n \"$FLAI_STORY\" ] || exit 0; flai guard >/dev/null 2>&1; exit 0"
          }
        ]
      }
    ],
    "SubagentStop": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "[ -n \"$FLAI_STORY\" ] || exit 0; flai guard >/dev/null 2>&1; exit 0"
          }
        ]
      }
    ]
  }
}
```

### 2026-11-02T00:13:07Z orchestrator (recommendation)
Recommendation: reply `allow`. Only alex's `allow` counts here. This entry is not an answer and refuses nothing.

The proposed file differs from the current `template/root/.claude/settings.json` in only the two PreToolUse commands. Each now refuses the call (exit 2) when `flai guard` exits 126 or 127, meaning it cannot be found or run. Every other exit still fails open as before. The matchers, the `FLAI_ROLE`/`FLAI_STORY` gate on Edit and Write, and the SubagentStart and SubagentStop hooks are unchanged.

This is the fix I-0216 asks for: "a guard that cannot start refuses rather than allows". It matches ADR-0235 on `story/S-0518` (commit 16712f93c).

Marked as the operator's alone: it changes the Claude Code hook config that every project made from the template inherits. flai serve must not confirm it on your behalf.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/issues/I-0216-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md

### 2026-11-02T00:39:50Z alex
allow

### 2026-11-02T00:40:08Z agent-S-0518
Resolved: allowed by alex: agent-S-0518 may Write template/root/.claude/settings.json
