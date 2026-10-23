---
id: D-0201
title: The MCP prime pack for a story is larger than Claude Code's tool result limit, so the agent reads it back from a saved file
kind: cause
sources: [I-0068]
topics: [cli, conventions, flai, go, code, all]
paths: [flai/internal/context, flai/internal/mcpserver, flai/cmd/prime.go, flai/cmd/prime_test.go, design/adrs, design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/issues/I-0068-the-mcp-prime-pack-for-a-story-is-larger-than-claude-code-s-tool-result-limit-so-the-agent-reads-it-back-from-a-saved-file.md, design/issues/summary.md]
date: 2026-10-06
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0201 The MCP prime pack for a story is larger than Claude Code's tool result limit, so the agent reads it back from a saved file

I-0068 was closed on 2026-10-06.

The cause:

```text
The MCP prime pack for a story is larger than Claude Code's tool result limit, so the agent reads it back from a saved file
```

The fix:

```text
S-0261, ADR-0104: the MCP prime tool returns the context pack in parts of at most 40,000 bytes of JSON (context.PartLimit, a fifth under the 50,000 characters over which Claude Code saves a tool result to a file), and the agent reads every part; flai prime takes --part N. Reproduced and checked by TestAPackOverTheToolResultLimitIsCutIntoPartsThatFit (flai/internal/context) and TestPrimeComesInPartsOnThisRepository (flai/internal/mcpserver), which prime S-0205, S-0210, S-0211, S-0261, and the planner's pack for S-0261: each was over 160 KB as one result and now comes in 5 or 6 parts of at most 39,994 bytes.
```

Source: I-0068.
