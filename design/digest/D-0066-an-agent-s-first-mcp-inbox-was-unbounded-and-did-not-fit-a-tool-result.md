---
id: D-0066
title: An agent's first MCP inbox was unbounded and did not fit a tool result
kind: cause
sources: [I-0023]
topics: [flai, cli, go, code, all]
paths: [flai/internal/mcpserver]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0066 An agent's first MCP inbox was unbounded and did not fit a tool result

I-0023 was closed on 2026-10-01.

The cause:

```text
An agent's first MCP inbox was unbounded and did not fit a tool result
```

The fix:

```text
Fixed by S-0061 (866c980): inbox is capped at 50 changes with changes_omitted, and a first look covers 24 hours of stories and epics.
```

Source: I-0023.
