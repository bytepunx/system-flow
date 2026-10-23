---
id: D-0133
title: "The planner's Edit and Write refusal needs the guard hook to see those tools: settings.json, or --settings on the planner's run: I have made the changes to both files as requested"
kind: answer
sources: [TH-0096]
topics: [orchestration, planning, analysis, flai, cli, go, flaiover, dashboard, sveltekit, template, conventions, code, all]
paths: [flai/internal/serve, flai/internal/harness, flai/internal/hostapi, flai/cmd, flai/internal/mcpserver, flai/internal/manifest, flai/internal/guard, flaiover/src, ".claude/agents", template, design/issues, design/system, design/adrs, docs/users, docs/operators, ".claude/settings.json"]
date: 2026-10-04
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0133 The planner's Edit and Write refusal needs the guard hook to see those tools: settings.json, or --settings on the planner's run: I have made the changes to both files as requested

The operator, alex, answered TH-0096 at 2026-10-04T03:21:57Z.

What was asked:

```text
T-0785 taught `flai guard` to refuse the planner's Edit, Write, and NotebookEdit calls (with `FLAI_ROLE=plan`). Claude Code runs the hook only for tools its matcher names, and the matcher in `.claude/settings.json` and `template/root/.claude/settings.json` is `Bash|mcp__flai__.*`. Changing either file is refused by this session's permissions, so I have not changed them.
```

The answer:

```text
I have made the changes to both files as requested
```

Source: TH-0096.
