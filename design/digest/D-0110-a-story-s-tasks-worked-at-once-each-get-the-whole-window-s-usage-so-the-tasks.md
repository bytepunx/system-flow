---
id: D-0110
title: A story's tasks worked at once each get the whole window's usage, so the tasks sum to more than the story
kind: cause
sources: [I-0054]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, template, conventions, code, all]
paths: [flai/internal/usage, flai/internal/serve, flai/internal/harness, flai/cmd, design/conventions, template, design/adrs, design/issues, design/system, docs]
date: 2026-10-02
status: absorbed
absorbed_by: "design/adrs/0071-a-task-s-usage-is-the-calls-of-the-sub-agents-started-for-it-and-an-even-share.md#decision"
asked_again: []
cited: []
---

# D-0110 A story's tasks worked at once each get the whole window's usage, so the tasks sum to more than the story

I-0054 was closed on 2026-10-02.

The cause:

```text
A story's tasks worked at once each get the whole window's usage, so the tasks sum to more than the story
```

The fix:

```text
S-0230: a task is measured by its sub-agents' calls (parent_tool_use_id, matched to the task ID its Agent call names) and an even share of the story's agent's calls while it was in progress (ADR-0071)
```

Source: I-0054.
