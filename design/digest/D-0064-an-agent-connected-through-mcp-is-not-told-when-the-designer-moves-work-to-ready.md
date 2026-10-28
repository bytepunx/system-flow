---
id: D-0064
title: An agent connected through MCP is not told when the designer moves work to ready
kind: cause
sources: [I-0020]
topics: [flai, cli, go, template, conventions, code, all]
paths: [flai/internal/mcpserver, design/conventions, template]
date: 2026-10-01
status: absorbed
absorbed_by: "design/system/agent-narrative.md#what-an-agent-is-told-through-mcp"
asked_again: []
cited: []
---

# D-0064 An agent connected through MCP is not told when the designer moves work to ready

I-0020 was closed on 2026-10-01.

The cause:

```text
An agent connected through MCP is not told when the designer moves work to ready
```

The fix:

```text
Fixed by S-0058 (7d3f502): inbox lists ready stories and changes since the agent last looked, with the board tool and wait_for_work.
```

Source: I-0020.
