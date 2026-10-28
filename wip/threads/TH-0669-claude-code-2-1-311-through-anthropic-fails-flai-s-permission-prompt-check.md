---
id: TH-0669
title: Claude Code 2.1.311 through anthropic fails flai's permission_prompt check
anchor:
  path: system-flow.yaml
status: open
participants: [flai serve]
created: 2026-10-26T21:55:54Z
updated: 2026-10-26T21:55:54Z
---

# TH-0669 Claude Code 2.1.311 through anthropic fails flai's permission_prompt check

On system-flow.yaml.

## Entries

### 2026-10-26T21:55:54Z flai serve
flai serve checked Claude Code 2.1.311 (/home/alex/.local/share/claude/versions/2.1.311) through anthropic against permission_prompt (ADR-0106), and it failed. A headless claude -p in a scratch project, with auto-approve on and the permission arguments a story's agent gets, was asked to Write .claude/flai-check.md in an in-progress story's worktree. What went wrong, and what Claude Code said:

```text
.flai-cache/worktrees/S-0001/.claude/flai-check.md was not written; claude ended: exit status 1
session ended (success, 1 turns, 0s): You've hit your weekly limit · resets Oct 28, 2am (America/Denver)
```

Until this is fixed, writes to the paths Claude Code protects (.claude/, .mcp.json, and the rest ADR-0106 lists) will not go through permission_prompt: an agent's write to one is refused. The agents still run, and a story that writes no protected path is not held up. flai does not check 2.1.311 through anthropic again; to check it again, remove the anthropic entry under 2.1.311 from /home/alex/.flai/serve/claude-checks.json.
