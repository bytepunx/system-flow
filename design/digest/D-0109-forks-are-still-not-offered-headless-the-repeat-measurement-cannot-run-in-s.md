---
id: D-0109
title: "Forks are still not offered headless: the repeat measurement cannot run in S-0230: use rec"
kind: answer
sources: [TH-0069]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, template, conventions, code, all]
paths: [flai/internal/usage, flai/internal/serve, flai/internal/harness, flai/cmd, design/conventions, template, design/adrs, design/issues, design/system, docs]
date: 2026-10-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0109 Forks are still not offered headless: the repeat measurement cannot run in S-0230: use rec

The operator, alex, answered TH-0069 at 2026-10-02T17:13:52Z.

What was asked:

```text
The fourth criterion asks that the measurement be repeated once forked sub-agents are offered to a headless session. They still are not: in this `claude -p` session (Claude Code 2.1.286), started by flai serve at 17:1xZ today, the Agent tool refused `subagent_type: fork` with "Agent type 'fork' not found. Available agents: claude, Explore, explorer, general-purpose, Plan, statusline-setup, verifier", the same as S-0176 found, and `claude --help` has no flag that offers them.
```

The answer:

```text
use rec
```

Source: TH-0069.
