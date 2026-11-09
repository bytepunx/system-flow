---
id: D-0503
title: "Claude Code BEEFGOOF through anthropic fails flai's permission_prompt check: the entry was removed so we can test this again"
kind: answer
sources: [TH-0819, I-0296]
topics: [agents, flai]
paths: [system-flow.yaml]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0935]
---

# D-0503 Claude Code BEEFGOOF through anthropic fails flai's permission_prompt check: the entry was removed so we can test this again

When flai serve's permission_prompt check of a Claude Code version through a provider fails (ADR-0106), flai does not check that pair again until its entry is removed from `/home/alex/.flai/serve/claude-checks.json`. The operator removes the entry to have it checked again.

flai serve opened TH-0819 on 2026-11-03: the check of Claude Code BEEFGOOF (`/home/alex/.local/bin/rmhclod`) through anthropic failed. The headless `claude -p` saw only the word "This" of its prompt, and the wrapper script itself failed: `rmhclod: 7: Syntax error: Unterminated quoted string`. Until the check passes, an agent's write to a path Claude Code protects (`.claude/`, `.mcp.json`, and the rest ADR-0106 lists) is refused; agents still run.

The operator, alex, answered at 2026-11-07T08:59:18Z: "the entry was removed so we can test this again".

Source: TH-0819.
