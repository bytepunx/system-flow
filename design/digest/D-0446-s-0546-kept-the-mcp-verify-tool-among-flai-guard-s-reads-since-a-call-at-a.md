---
id: D-0446
title: "S-0546 kept the MCP verify tool among flai guard's reads, since a call at a verified commit answers from the record and a story in review has nothing to commit; I-0240 records that a sub-agent can still start a committing close-out"
kind: decision
sources: [S-0546, I-0240, S-0556]
topics: [mcp, orchestration, flai, testing]
paths: [flai/internal/guard, flai/internal/mcpserver/verify.go, design/conventions/delegation.md]
date: 2026-11-01
status: absorbed
absorbed_by: "design/system/flai-cli.md#commands"
asked_again: []
cited: [S-0556, T-2459, TH-0827, TH-0849]
---

# D-0446 S-0546 kept the MCP verify tool among flai guard's reads, since a call at a verified commit answers from the record and a story in review has nothing to commit; I-0240 records that a sub-agent can still start a committing close-out

S-0546 made the MCP tool `verify` ask flai serve for the close-out job (ADR-0229), which runs `scripts/close-out.sh` and commits in the story's worktree. It left the tool in `flai guard`'s list of reads, which every sub-agent and strategic agent may call. Its narrative gives the reason (wip/archive/agents/S-0546.md, Decisions):

> The MCP `verify` tool stays among the guard's reads: a sub-agent or the orchestrator asking at a verified commit gets the record, and a story in review has nothing to commit.

The story's agent also filed I-0240 against that choice: "guard.MCPReads still lets every sub-agent and strategic agent call it, so a verifier or task sub-agent could start a close-out that commits; the delegation convention now forbids it unless the prompt asks, but the guard does not enforce it."

S-0556 reversed the choice and closed I-0240 on 2026-11-07, as the operator confirmed on TH-0849 (D-0495): `verify` left `guard.MCPReads`, and the guard refuses it and `flai verify --job` to every sub-agent and every strategic agent but the orchestrator, which holds the tool by name. The two calls S-0546 kept the tool for stay: the orchestrator's verify of a story in review, and a reader's look at the record, now through `flai verify --last`. `design/system/flai-cli.md` § Commands, the `flai guard` row, states the rule as it stands.

Sources: S-0546, I-0240, S-0556.
