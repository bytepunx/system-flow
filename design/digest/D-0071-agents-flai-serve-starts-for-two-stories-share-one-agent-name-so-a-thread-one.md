---
id: D-0071
title: Agents flai serve starts for two stories share one agent name, so a thread one opens for the other never reads as awaiting it
kind: cause
sources: [I-0037]
topics: [flai, cli, go, code, all]
paths: [flai/cmd, flai/internal/serve, flai/internal/harness, design/adrs, design/system, docs/users, docs/operators]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0071 Agents flai serve starts for two stories share one agent name, so a thread one opens for the other never reads as awaiting it

I-0037 was closed on 2026-10-01.

The cause:

```text
Agents flai serve starts for two stories share one agent name, so a thread one opens for the other never reads as awaiting it
```

The fix:

```text
Fixed by S-0114 (0eac71b): the claude-code adapter passes the agent's name to flai mcp as --agent, over any FLAI_AGENT in the harness's settings.
```

Source: I-0037.
