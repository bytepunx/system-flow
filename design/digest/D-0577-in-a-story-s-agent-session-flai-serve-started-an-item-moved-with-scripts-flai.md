---
id: D-0577
title: In a story's agent session flai serve started, an item moved with scripts/flai.sh records the session's FLAI_AGENT, system-flow, not agent-S-nnnn, so the agent moves items through the MCP tools, which record its own name
kind: decision
sources: [S-0665]
topics: [harness, cli, flai, all]
paths: [scripts/flai.sh]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0577 In a story's agent session flai serve started, an item moved with scripts/flai.sh records the session's FLAI_AGENT, system-flow, not agent-S-nnnn, so the agent moves items through the MCP tools, which record its own name

S-0665's agent found that its moves carried the wrong author. From S-0665's Decisions (`wip/archive/agents/S-0665.md`):

```text
Moves made with `scripts/flai.sh` record `system-flow`, the `FLAI_AGENT` flai serve gave this session, not `agent-S-0665`; later moves go through the MCP tools, which record `agent-S-0665`.
```

So in this repository a story's agent makes its transitions, and any other write that records an author, through the flai MCP tools (`item_move`, `task_done`, `stream_state`), not through `scripts/flai.sh` in its shell, whose author is the session's `FLAI_AGENT`. It does not set `FLAI_AGENT` itself to work round it: `flai guard` refuses a call that sets it to anything but the session's own (ADR-0143, D-0414). No issue records the `system-flow` name yet.

Source: S-0665 (Decisions).
